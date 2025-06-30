package nuclei

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"drone-security-scanner/internal/models"
	"github.com/google/uuid"
)

type Scanner struct {
	templatesPath string
}

type NucleiResult struct {
	TemplateID   string                 `json:"template-id"`
	TemplatePath string                 `json:"template-path"`
	Info         NucleiInfo             `json:"info"`
	Type         string                 `json:"type"`
	Host         string                 `json:"host"`
	MatchedAt    string                 `json:"matched-at"`
	ExtractedResults []string           `json:"extracted-results"`
	Request      string                 `json:"request"`
	Response     string                 `json:"response"`
	Metadata     map[string]interface{} `json:"metadata"`
	Timestamp    time.Time              `json:"timestamp"`
}

type NucleiInfo struct {
	Name        string            `json:"name"`
	Author      []string          `json:"author"`
	Tags        []string          `json:"tags"`
	Description string            `json:"description"`
	Reference   []string          `json:"reference"`
	Severity    string            `json:"severity"`
	Metadata    map[string]string `json:"metadata"`
}

func NewScanner(templatesPath string) *Scanner {
	return &Scanner{
		templatesPath: templatesPath,
	}
}

func (s *Scanner) Scan(ctx context.Context, targets []string, config *models.ScanConfig) ([]*models.Vulnerability, error) {
	// 构建 nuclei 命令
	args := []string{
		"-json",
		"-silent",
		"-no-color",
	}

	// 添加目标
	for _, target := range targets {
		args = append(args, "-target", target)
	}

	// 添加模板路径
	if s.templatesPath != "" {
		args = append(args, "-t", s.templatesPath)
	}

	// 添加特定模板
	if len(config.Templates) > 0 {
		for _, template := range config.Templates {
			args = append(args, "-t", template)
		}
	}

	// 添加并发和速率限制
	if config.Concurrency > 0 {
		args = append(args, "-c", fmt.Sprintf("%d", config.Concurrency))
	}

	if config.RateLimit > 0 {
		args = append(args, "-rl", fmt.Sprintf("%d", config.RateLimit))
	}

	// 添加超时
	if config.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%ds", config.Timeout))
	}

	// 添加自定义参数
	for key, value := range config.CustomParams {
		args = append(args, fmt.Sprintf("-%s", key), value)
	}

	// 执行命令
	cmd := exec.CommandContext(ctx, "nuclei", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start nuclei: %w", err)
	}

	// 解析结果
	var vulnerabilities []*models.Vulnerability
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var result NucleiResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			// 跳过无法解析的行
			continue
		}

		vuln := s.convertToVulnerability(&result)
		vulnerabilities = append(vulnerabilities, vuln)
	}

	if err := cmd.Wait(); err != nil {
		// nuclei 可能会返回非零退出码，但仍然产生有效结果
		// 只有在没有找到任何漏洞时才报告错误
		if len(vulnerabilities) == 0 {
			return nil, fmt.Errorf("nuclei scan failed: %w", err)
		}
	}

	return vulnerabilities, nil
}

func (s *Scanner) convertToVulnerability(result *NucleiResult) *models.Vulnerability {
	vuln := &models.Vulnerability{
		ID:           uuid.New().String(),
		Title:        result.Info.Name,
		Description:  result.Info.Description,
		Severity:     s.normalizeSeverity(result.Info.Severity),
		Category:     s.categorizeByTags(result.Info.Tags),
		Status:       "open",
		DiscoveredAt: result.Timestamp,
		UpdatedAt:    result.Timestamp,
	}

	// 设置 CVE ID（如果存在）
	if cveID := s.extractCVEID(result); cveID != "" {
		vuln.CVEID = cveID
	}

	// 设置 CVSS 评分（如果存在）
	if cvssScore := s.extractCVSSScore(result); cvssScore > 0 {
		vuln.CVSSScore = cvssScore
	}

	// 设置 CVSS 向量（如果存在）
	if cvssVector := s.extractCVSSVector(result); cvssVector != "" {
		vuln.CVSSVector = cvssVector
	}

	// 设置修复建议
	vuln.Remediation = s.generateRemediation(result)

	// 设置参考链接
	if len(result.Info.Reference) > 0 {
		references, _ := json.Marshal(result.Info.Reference)
		vuln.References = string(references)
	}

	return vuln
}

func (s *Scanner) normalizeSeverity(severity string) string {
	severity = strings.ToLower(severity)
	switch severity {
	case "critical", "high", "medium", "low":
		return severity
	case "info", "informational":
		return "low"
	default:
		return "medium"
	}
}

func (s *Scanner) categorizeByTags(tags []string) string {
	// 根据标签确定漏洞类别
	for _, tag := range tags {
		switch strings.ToLower(tag) {
		case "sqli", "sql-injection":
			return "SQL注入"
		case "xss", "cross-site-scripting":
			return "跨站脚本"
		case "rce", "remote-code-execution":
			return "远程代码执行"
		case "lfi", "local-file-inclusion":
			return "本地文件包含"
		case "rfi", "remote-file-inclusion":
			return "远程文件包含"
		case "ssrf", "server-side-request-forgery":
			return "服务端请求伪造"
		case "auth-bypass", "authentication-bypass":
			return "认证绕过"
		case "privilege-escalation":
			return "权限提升"
		case "information-disclosure":
			return "信息泄露"
		case "dos", "denial-of-service":
			return "拒绝服务"
		}
	}
	return "其他"
}

func (s *Scanner) extractCVEID(result *NucleiResult) string {
	// 从模板 ID 或元数据中提取 CVE ID
	if strings.HasPrefix(result.TemplateID, "CVE-") {
		return result.TemplateID
	}

	// 检查元数据
	if cve, exists := result.Info.Metadata["cve-id"]; exists {
		return cve
	}

	// 检查参考链接中的 CVE
	for _, ref := range result.Info.Reference {
		if strings.Contains(ref, "CVE-") {
			// 提取 CVE ID
			parts := strings.Split(ref, "CVE-")
			if len(parts) > 1 {
				cveID := "CVE-" + strings.Split(parts[1], "/")[0]
				return cveID
			}
		}
	}

	return ""
}

func (s *Scanner) extractCVSSScore(result *NucleiResult) float64 {
	// 从元数据中提取 CVSS 评分
	if cvss, exists := result.Info.Metadata["cvss-score"]; exists {
		var score float64
		if _, err := fmt.Sscanf(cvss, "%f", &score); err == nil {
			return score
		}
	}

	// 根据严重程度估算评分
	switch s.normalizeSeverity(result.Info.Severity) {
	case "critical":
		return 9.0
	case "high":
		return 7.5
	case "medium":
		return 5.0
	case "low":
		return 2.5
	default:
		return 0.0
	}
}

func (s *Scanner) extractCVSSVector(result *NucleiResult) string {
	// 从元数据中提取 CVSS 向量
	if vector, exists := result.Info.Metadata["cvss-vector"]; exists {
		return vector
	}
	return ""
}

func (s *Scanner) generateRemediation(result *NucleiResult) string {
	// 根据漏洞类型生成修复建议
	category := s.categorizeByTags(result.Info.Tags)

	switch category {
	case "SQL注入":
		return "使用参数化查询或预编译语句，验证和过滤用户输入，实施最小权限原则。"
	case "跨站脚本":
		return "对用户输入进行适当的编码和验证，使用内容安全策略(CSP)，避免直接输出用户数据。"
	case "远程代码执行":
		return "更新到最新版本，禁用不必要的功能，实施严格的输入验证，使用沙箱环境。"
	case "认证绕过":
		return "修复认证逻辑漏洞，实施多因素认证，定期审查访问控制机制。"
	case "信息泄露":
		return "移除敏感信息的暴露，配置适当的错误页面，实施访问控制。"
	default:
		return "请参考相关安全最佳实践，及时更新系统和应用程序，实施适当的安全控制措施。"
	}
}

func (s *Scanner) GetAvailableTemplates() ([]string, error) {
	// 获取可用的模板列表
	cmd := exec.Command("nuclei", "-tl")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get template list: %w", err)
	}

	var templates []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "[") {
			templates = append(templates, line)
		}
	}

	return templates, nil
}

func (s *Scanner) UpdateTemplates() error {
	// 更新模板
	cmd := exec.Command("nuclei", "-update-templates")
	return cmd.Run()
}
