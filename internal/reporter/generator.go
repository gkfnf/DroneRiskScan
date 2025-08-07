package reporter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/pkg/models"
)

// ReportGenerator 报告生成器接口
type ReportGenerator interface {
	GenerateReport(result *models.ScanResult, format string, outputPath string) error
	GetSupportedFormats() []string
}

// Config 报告生成器配置
type Config struct {
	Formats       []string
	TemplateDir   string
	OutputDir     string
	Debug         bool
}

// DefaultReportGenerator 默认报告生成器实现
type DefaultReportGenerator struct {
	config    *Config
	templates map[string]*template.Template
}

// NewReportGenerator 创建新的报告生成器
func NewReportGenerator(config *Config) ReportGenerator {
	if config == nil {
		config = &Config{
			Formats: []string{"json", "html", "markdown"},
			Debug:   false,
		}
	}

	generator := &DefaultReportGenerator{
		config:    config,
		templates: make(map[string]*template.Template),
	}

	// 初始化模板
	generator.initTemplates()

	return generator
}

// GenerateReport 生成报告
func (rg *DefaultReportGenerator) GenerateReport(result *models.ScanResult, format string, outputPath string) error {
	if result == nil {
		return fmt.Errorf("扫描结果不能为空")
	}

	// 预处理报告数据
	if err := rg.preprocessReport(result); err != nil {
		return fmt.Errorf("预处理报告失败: %w", err)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 根据格式生成报告
	switch strings.ToLower(format) {
	case "json":
		return rg.generateJSONReport(result, outputPath)
	case "html":
		return rg.generateHTMLReport(result, outputPath)
	case "markdown", "md":
		return rg.generateMarkdownReport(result, outputPath)
	default:
		return fmt.Errorf("不支持的报告格式: %s", format)
	}
}

// GetSupportedFormats 获取支持的格式
func (rg *DefaultReportGenerator) GetSupportedFormats() []string {
	return []string{"json", "html", "markdown"}
}

// preprocessReport 预处理报告数据
func (rg *DefaultReportGenerator) preprocessReport(result *models.ScanResult) error {
	// 设置默认值
	if result.ID == "" {
		result.ID = fmt.Sprintf("scan_%d", time.Now().Unix())
	}

	// 排序漏洞（按严重程度和发现时间）
	sort.Slice(result.Vulnerabilities, func(i, j int) bool {
		if result.Vulnerabilities[i].Severity != result.Vulnerabilities[j].Severity {
			return result.Vulnerabilities[i].Severity > result.Vulnerabilities[j].Severity
		}
		return result.Vulnerabilities[i].Timestamp.Before(result.Vulnerabilities[j].Timestamp)
	})

	return nil
}

// generateJSONReport 生成JSON报告
func (rg *DefaultReportGenerator) generateJSONReport(result *models.ScanResult, outputPath string) error {
	// 创建报告数据结构
	reportData := map[string]interface{}{
		"id":              result.ID,
		"title":          "DroneRiskScan Security Assessment Report",
		"description":    "Comprehensive security assessment report",
		"scan_info": map[string]interface{}{
			"start_time": result.StartTime.Format(time.RFC3339),
			"end_time":   result.EndTime.Format(time.RFC3339),
			"duration":   result.Duration.Nanoseconds(),
			"status":     string(result.Status),
		},
		"targets":         result.Targets,
		"vulnerabilities": result.Vulnerabilities,
		"statistics":      result.Statistics,
		"generated_by":    "DroneRiskScan",
		"version":         "1.0",
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	// 生成修复建议
	recommendations := rg.generateRecommendations(result.Vulnerabilities)
	reportData["recommendations"] = recommendations

	// 序列化为JSON
	data, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	// 写入文件
	return os.WriteFile(outputPath, data, 0644)
}

// generateHTMLReport 生成HTML报告
func (rg *DefaultReportGenerator) generateHTMLReport(result *models.ScanResult, outputPath string) error {
	// 使用内置HTML模板
	htmlTemplate := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DroneRiskScan 安全扫描报告</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            min-height: 100vh; 
            color: #333; 
            line-height: 1.6; 
        }
        .container { max-width: 1400px; margin: 0 auto; padding: 20px; }
        .report-header { 
            background: white; 
            border-radius: 20px; 
            padding: 40px; 
            margin-bottom: 30px; 
            box-shadow: 0 20px 60px rgba(0,0,0,0.1); 
        }
        .report-header h1 { color: #333; font-size: 2.5em; margin-bottom: 10px; }
        .report-header .meta { color: #666; }
        .stats-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); 
            gap: 20px; 
            margin-bottom: 30px; 
        }
        .stat-card { 
            background: white; 
            padding: 25px; 
            border-radius: 15px; 
            box-shadow: 0 10px 30px rgba(0,0,0,0.1); 
            text-align: center; 
            transition: all 0.3s; 
        }
        .stat-card:hover { 
            transform: translateY(-10px); 
            box-shadow: 0 20px 40px rgba(0,0,0,0.15); 
        }
        .stat-card .number { 
            font-size: 3em; 
            font-weight: bold; 
            margin-bottom: 10px; 
        }
        .stat-card.critical .number { color: #dc3545; }
        .stat-card.high .number { color: #fd7e14; }
        .stat-card.medium .number { color: #ffc107; }
        .stat-card.low .number { color: #28a745; }
        .stat-card.info .number { color: #17a2b8; }
        .stat-card.total .number { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            -webkit-background-clip: text; 
            -webkit-text-fill-color: transparent; 
        }
        .stat-card .label { 
            color: #888; 
            font-size: 1.1em; 
            text-transform: uppercase; 
            letter-spacing: 1px; 
        }
        .vulns-section { 
            background: white; 
            border-radius: 20px; 
            padding: 40px; 
            box-shadow: 0 20px 60px rgba(0,0,0,0.1); 
            margin-bottom: 20px;
        }
        .vulns-section h2 { color: #333; margin-bottom: 30px; font-size: 2em; }
        .vuln-item { 
            border-left: 5px solid #ddd; 
            padding: 20px; 
            margin-bottom: 20px; 
            background: #f8f9fa; 
            border-radius: 10px; 
            transition: all 0.3s; 
        }
        .vuln-item:hover { 
            box-shadow: 0 5px 15px rgba(0,0,0,0.1); 
            transform: translateX(5px); 
        }
        .vuln-item.severity-4 { 
            border-left-color: #dc3545; 
            background: linear-gradient(to right, #fff5f5, white); 
        }
        .vuln-item.severity-3 { 
            border-left-color: #fd7e14; 
            background: linear-gradient(to right, #fff9f5, white); 
        }
        .vuln-item.severity-2 { 
            border-left-color: #ffc107; 
            background: linear-gradient(to right, #fffef5, white); 
        }
        .vuln-item.severity-1 { 
            border-left-color: #28a745; 
            background: linear-gradient(to right, #f5fff5, white); 
        }
        .vuln-item.severity-0 { 
            border-left-color: #17a2b8; 
            background: linear-gradient(to right, #f5fffe, white); 
        }
        .vuln-header { 
            display: flex; 
            justify-content: space-between; 
            align-items: center; 
            margin-bottom: 15px; 
        }
        .vuln-title { font-size: 1.3em; font-weight: bold; color: #333; }
        .severity-badge { 
            padding: 5px 15px; 
            border-radius: 20px; 
            font-size: 0.9em; 
            font-weight: bold; 
            text-transform: uppercase; 
        }
        .severity-4 { background: #dc3545; color: white; }
        .severity-3 { background: #fd7e14; color: white; }
        .severity-2 { background: #ffc107; color: black; }
        .severity-1 { background: #28a745; color: white; }
        .severity-0 { background: #17a2b8; color: white; }
        .vuln-details { color: #666; line-height: 1.8; }
        .vuln-details strong { color: #333; }
        .no-vulns { 
            text-align: center; 
            padding: 60px; 
            color: #28a745; 
            font-size: 1.5em; 
        }
        .footer { 
            text-align: center; 
            color: white; 
            margin-top: 40px; 
            padding: 20px; 
        }
        .footer a { color: white; text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <div class="report-header">
            <h1>🚁 DroneRiskScan Security Report</h1>
            <div class="meta">
                <p>📅 Generated: {{.GeneratedTime}}</p>
                <p>🎯 Targets: {{.TargetCount}}</p>
                <p>📊 Total Vulnerabilities: {{.VulnCount}}</p>
                <p>⏱️ Scan Duration: {{.Duration}}</p>
            </div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card total">
                <div class="number">{{.VulnCount}}</div>
                <div class="label">Total</div>
            </div>
            <div class="stat-card critical">
                <div class="number">{{.CriticalCount}}</div>
                <div class="label">Critical</div>
            </div>
            <div class="stat-card high">
                <div class="number">{{.HighCount}}</div>
                <div class="label">High</div>
            </div>
            <div class="stat-card medium">
                <div class="number">{{.MediumCount}}</div>
                <div class="label">Medium</div>
            </div>
            <div class="stat-card low">
                <div class="number">{{.LowCount}}</div>
                <div class="label">Low</div>
            </div>
            <div class="stat-card info">
                <div class="number">{{.InfoCount}}</div>
                <div class="label">Info</div>
            </div>
        </div>
        
        <div class="vulns-section">
            <h2>🔍 Vulnerability Details</h2>
            {{if .HasVulnerabilities}}
                {{range .Vulnerabilities}}
                <div class="vuln-item severity-{{.Severity.Value}}">
                    <div class="vuln-header">
                        <div class="vuln-title">{{.Title}}</div>
                        <span class="severity-badge severity-{{.Severity.Value}}">{{.Severity.String}}</span>
                    </div>
                    <div class="vuln-details">
                        <p><strong>URL:</strong> {{.URL}}</p>
                        {{if .Parameter}}<p><strong>Parameter:</strong> {{.Parameter}} ({{.Position}})</p>{{end}}
                        {{if .Payload}}<p><strong>Payload:</strong> <code>{{.Payload}}</code></p>{{end}}
                        <p><strong>Description:</strong> {{.Description}}</p>
                        {{if .Evidence}}<p><strong>Evidence:</strong> {{.Evidence}}</p>{{end}}
                        {{if .CWE}}<p><strong>CWE:</strong> {{.CWE}} | <strong>CVSS:</strong> {{.CVSS}}</p>{{end}}
                        <p><strong>Confidence:</strong> {{printf "%.0f%%" (mul .Confidence 100)}}</p>
                    </div>
                </div>
                {{end}}
            {{else}}
                <div class="no-vulns">✅ No vulnerabilities found!</div>
            {{end}}
        </div>
        
        <div class="footer">
            <p>Generated by DroneRiskScan v1.0 | Professional Drone Security Scanner</p>
        </div>
    </div>
</body>
</html>`

	// 准备模板数据
	data := rg.prepareHTMLData(result)

	// 解析模板
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("解析HTML模板失败: %w", err)
	}

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建HTML文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	return tmpl.Execute(file, data)
}

// generateMarkdownReport 生成Markdown报告
func (rg *DefaultReportGenerator) generateMarkdownReport(result *models.ScanResult, outputPath string) error {
	var md strings.Builder

	// 标题和基本信息
	md.WriteString("# DroneRiskScan Security Assessment Report\n\n")
	md.WriteString("**Generated by:** DroneRiskScan 1.0  \n")
	md.WriteString("**Generated on:** " + time.Now().Format("2006-01-02 15:04:05") + "  \n")
	md.WriteString("**Report ID:** " + result.ID + "\n\n")

	// 执行摘要
	md.WriteString("## Executive Summary\n\n")
	md.WriteString(fmt.Sprintf("This security assessment report presents the findings from scanning %d targets over %s. ",
		len(result.Targets), result.Duration))
	
	vulnCount := len(result.Vulnerabilities)
	if vulnCount > 0 {
		md.WriteString(fmt.Sprintf("The scan identified **%d vulnerabilities** across different severity levels.\n\n", vulnCount))
	} else {
		md.WriteString("No vulnerabilities were identified.\n\n")
	}

	// 关键发现
	md.WriteString("### Key Findings\n\n")
	stats := result.Statistics
	md.WriteString(fmt.Sprintf("- **Critical:** %d vulnerabilities\n", stats.VulnsBySeverity[models.SeverityCritical]))
	md.WriteString(fmt.Sprintf("- **High:** %d vulnerabilities  \n", stats.VulnsBySeverity[models.SeverityHigh]))
	md.WriteString(fmt.Sprintf("- **Medium:** %d vulnerabilities\n", stats.VulnsBySeverity[models.SeverityMedium]))
	md.WriteString(fmt.Sprintf("- **Low:** %d vulnerabilities\n", stats.VulnsBySeverity[models.SeverityLow]))
	md.WriteString(fmt.Sprintf("- **Info:** %d vulnerabilities\n\n", stats.VulnsBySeverity[models.SeverityInfo]))

	// 风险概述
	md.WriteString("### Risk Overview\n\n")
	md.WriteString(fmt.Sprintf("- **Targets Scanned:** %d\n", stats.TargetsScanned))
	md.WriteString(fmt.Sprintf("- **Targets with Vulnerabilities:** %d\n", stats.TargetsWithVulns))
	
	overallRisk := "Low"
	if stats.VulnsBySeverity[models.SeverityCritical] > 0 {
		overallRisk = "Critical"
	} else if stats.VulnsBySeverity[models.SeverityHigh] > 0 {
		overallRisk = "High"
	} else if stats.VulnsBySeverity[models.SeverityMedium] > 0 {
		overallRisk = "Medium"
	}
	md.WriteString(fmt.Sprintf("- **Overall Risk Level:** %s\n\n\n", overallRisk))

	// 漏洞详情
	if len(result.Vulnerabilities) > 0 {
		md.WriteString("## Vulnerability Details\n\n")
		
		for i, vuln := range result.Vulnerabilities {
			md.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, vuln.Title))
			md.WriteString(fmt.Sprintf("- **Severity:** %s\n", vuln.Severity.String()))
			md.WriteString(fmt.Sprintf("- **Type:** %s\n", vuln.Type))
			md.WriteString(fmt.Sprintf("- **URL:** %s\n", vuln.URL))
			if vuln.Parameter != "" {
				md.WriteString(fmt.Sprintf("- **Parameter:** %s\n", vuln.Parameter))
			}
			if vuln.Evidence != "" {
				md.WriteString(fmt.Sprintf("- **Evidence:** %s\n", vuln.Evidence))
			}
			if vuln.CWE != "" {
				md.WriteString(fmt.Sprintf("- **CWE:** %s\n", vuln.CWE))
			}
			if vuln.CVSS > 0 {
				md.WriteString(fmt.Sprintf("- **CVSS:** %.1f\n", vuln.CVSS))
			}
			md.WriteString("\n**Description:** " + vuln.Description + "\n\n")
			if vuln.Risk != "" {
				md.WriteString("**Risk:** " + vuln.Risk + "\n\n")
			}
			md.WriteString("---\n\n")
		}
	}

	// 修复建议
	recommendations := rg.generateRecommendations(result.Vulnerabilities)
	if len(recommendations) > 0 {
		md.WriteString("\n## Recommendations\n\n")
		
		for i, rec := range recommendations {
			md.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, rec.Title))
			md.WriteString(fmt.Sprintf("- **Priority:** %s\n", rec.Priority))
			md.WriteString(fmt.Sprintf("- **Category:** %s\n", rec.Category))
			md.WriteString(fmt.Sprintf("- **Effort:** %s\n", rec.Effort))
			md.WriteString(fmt.Sprintf("- **Impact:** %s\n\n", rec.Impact))
			md.WriteString("**Description:** " + rec.Description + "\n\n")
			md.WriteString("**Solution:** " + rec.Solution + "\n\n")
		}
	}

	// 写入文件
	return os.WriteFile(outputPath, []byte(md.String()), 0644)
}

// prepareHTMLData 准备HTML模板数据
func (rg *DefaultReportGenerator) prepareHTMLData(result *models.ScanResult) map[string]interface{} {
	stats := result.Statistics
	
	// 转换漏洞数据，添加字符串表示
	vulnerabilities := make([]map[string]interface{}, len(result.Vulnerabilities))
	for i, vuln := range result.Vulnerabilities {
		vulnerabilities[i] = map[string]interface{}{
			"Title":       vuln.Title,
			"Description": vuln.Description,
			"URL":         vuln.URL,
			"Parameter":   vuln.Parameter,
			"Position":    vuln.Position,
			"Payload":     vuln.Payload,
			"Evidence":    vuln.Evidence,
			"CWE":         vuln.CWE,
			"CVSS":        vuln.CVSS,
			"Confidence":  vuln.Confidence,
			"Severity": map[string]interface{}{
				"Value":  vuln.Severity.Value(),
				"String": vuln.Severity.String(),
			},
		}
	}

	return map[string]interface{}{
		"GeneratedTime":     time.Now().Format("2006-01-02 15:04:05"),
		"TargetCount":       len(result.Targets),
		"VulnCount":         len(result.Vulnerabilities),
		"Duration":          result.Duration.String(),
		"CriticalCount":     stats.VulnsBySeverity[models.SeverityCritical],
		"HighCount":         stats.VulnsBySeverity[models.SeverityHigh],
		"MediumCount":       stats.VulnsBySeverity[models.SeverityMedium],
		"LowCount":          stats.VulnsBySeverity[models.SeverityLow],
		"InfoCount":         stats.VulnsBySeverity[models.SeverityInfo],
		"HasVulnerabilities": len(result.Vulnerabilities) > 0,
		"Vulnerabilities":   vulnerabilities,
	}
}

// Recommendation 修复建议结构
type Recommendation struct {
	ID                string              `json:"id"`
	Title             string              `json:"title"`
	Priority          string              `json:"priority"`
	Category          models.Category     `json:"category"`
	VulnerabilityTypes []models.VulnType  `json:"vulnerability_types"`
	Description       string              `json:"description"`
	Solution          string              `json:"solution"`
	References        []string            `json:"references,omitempty"`
	Effort            string              `json:"effort"`
	Impact            string              `json:"impact"`
}

// generateRecommendations 生成修复建议
func (rg *DefaultReportGenerator) generateRecommendations(vulnerabilities []*models.Vulnerability) []Recommendation {
	// 按类别分组漏洞
	categoryVulns := make(map[models.Category][]models.VulnType)
	for _, vuln := range vulnerabilities {
		if _, exists := categoryVulns[vuln.Category]; !exists {
			categoryVulns[vuln.Category] = []models.VulnType{}
		}
		categoryVulns[vuln.Category] = append(categoryVulns[vuln.Category], vuln.Type)
	}

	var recommendations []Recommendation
	
	for category, vulnTypes := range categoryVulns {
		rec := rg.generateCategoryRecommendation(category, vulnTypes)
		if rec != nil {
			recommendations = append(recommendations, *rec)
		}
	}

	// 按优先级排序
	sort.Slice(recommendations, func(i, j int) bool {
		priorities := map[string]int{
			"Critical": 4, "High": 3, "Medium": 2, "Low": 1,
		}
		return priorities[recommendations[i].Priority] > priorities[recommendations[j].Priority]
	})

	return recommendations
}

// generateCategoryRecommendation 为特定类别生成建议
func (rg *DefaultReportGenerator) generateCategoryRecommendation(category models.Category, vulnTypes []models.VulnType) *Recommendation {
	// 去重漏洞类型
	uniqueTypes := make(map[models.VulnType]bool)
	for _, vt := range vulnTypes {
		uniqueTypes[vt] = true
	}
	
	typeList := make([]models.VulnType, 0, len(uniqueTypes))
	for vt := range uniqueTypes {
		typeList = append(typeList, vt)
	}

	switch category {
	case models.CategoryInjection:
		return &Recommendation{
			ID:                 fmt.Sprintf("rec_%s_1", category),
			Title:              "Fix " + string(category) + " Vulnerabilities",
			Priority:           "Critical",
			Category:           category,
			VulnerabilityTypes: typeList,
			Description:        fmt.Sprintf("Found %d %s vulnerabilities that need attention", len(vulnTypes), category),
			Solution:           "Implement input validation, use parameterized queries, and apply least privilege principles",
			References: []string{
				"https://owasp.org/www-community/attacks/SQL_Injection",
				"https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html",
			},
			Effort: "Medium",
			Impact: "Critical",
		}
	case models.CategoryAuth:
		return &Recommendation{
			ID:                 fmt.Sprintf("rec_%s_1", category),
			Title:              "Strengthen Authentication Controls",
			Priority:           "High",
			Category:           category,
			VulnerabilityTypes: typeList,
			Description:        fmt.Sprintf("Found %d authentication-related vulnerabilities", len(vulnTypes)),
			Solution:           "Implement proper authentication mechanisms, session management, and access controls",
			References: []string{
				"https://owasp.org/www-project-top-ten/2017/A2_2017-Broken_Authentication",
			},
			Effort: "High",
			Impact: "High",
		}
	case models.CategoryXSS:
		return &Recommendation{
			ID:                 fmt.Sprintf("rec_%s_1", category),
			Title:              "Implement XSS Protection",
			Priority:           "High",
			Category:           category,
			VulnerabilityTypes: typeList,
			Description:        fmt.Sprintf("Found %d XSS vulnerabilities", len(vulnTypes)),
			Solution:           "Implement proper output encoding, input validation, and Content Security Policy",
			References: []string{
				"https://owasp.org/www-community/attacks/xss/",
			},
			Effort: "Medium",
			Impact: "High",
		}
	default:
		return &Recommendation{
			ID:                 fmt.Sprintf("rec_%s_1", category),
			Title:              "Address " + string(category) + " Issues",
			Priority:           "Medium",
			Category:           category,
			VulnerabilityTypes: typeList,
			Description:        fmt.Sprintf("Found %d %s related issues", len(vulnTypes), category),
			Solution:           "Review and remediate the identified issues according to security best practices",
			Effort:             "Low",
			Impact:             "Medium",
		}
	}
}

// initTemplates 初始化模板
func (rg *DefaultReportGenerator) initTemplates() {
	// 这里可以加载外部模板文件
	// 目前使用内置模板，所以暂时为空
}