package injection

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/internal/detector"
	"github.com/dronesec/droneriskscan/internal/transport"
	"github.com/dronesec/droneriskscan/pkg/models"
)

// SQLiDetector SQL注入检测器
type SQLiDetector struct {
	*detector.BasePlugin
	requestModifier  *detector.RequestModifier
	responseAnalyzer *detector.ResponseAnalyzer
	paramExtractor   *detector.ParameterExtractor
}

// NewSQLiDetector 创建SQL注入检测器
func NewSQLiDetector(httpClient transport.HTTPClient) *SQLiDetector {
	base := detector.NewBasePlugin(
		"sqli-detector",
		detector.PluginTypeActive,
		models.CategoryInjection,
		models.SeverityHigh,
	)
	
	base.SetDescription("检测SQL注入漏洞，包括布尔盲注、错误注入和时间盲注")
	base.SetAuthor("DroneRiskScan Team")
	base.SetVersion("1.2.0")
	base.SetHTTPClient(httpClient)

	return &SQLiDetector{
		BasePlugin:       base,
		requestModifier:  detector.NewRequestModifier(httpClient),
		responseAnalyzer: detector.NewResponseAnalyzer(),
		paramExtractor:   detector.NewParameterExtractor(),
	}
}

// SetSessionCookies 设置会话Cookie
func (s *SQLiDetector) SetSessionCookies(cookies []*http.Cookie) {
	s.requestModifier.SetSessionCookies(cookies)
}

// Execute 执行SQL注入检测
func (s *SQLiDetector) Execute(ctx context.Context, target *detector.ScanTarget) (*detector.DetectionResult, error) {
	result := &detector.DetectionResult{
		IsVulnerable:    false,
		Vulnerabilities: []*models.Vulnerability{},
		Evidence:        []detector.Evidence{},
		Metadata:        make(map[string]interface{}),
	}

	// 提取所有注入点
	injectPoints := s.paramExtractor.ExtractParameters(target)
	if len(injectPoints) == 0 {
		result.Metadata["message"] = "未发现可注入参数"
		return result, nil
	}
	
	// 调试信息
	fmt.Printf("[DEBUG] SQL注入检测器找到 %d 个注入点\n", len(injectPoints))
	for _, point := range injectPoints {
		fmt.Printf("[DEBUG] 注入点: %s=%s (位置: %s, 类型: %s)\n", point.Name, point.Value, point.Position, point.Type)
	}

	// 对每个注入点进行测试
	for _, point := range injectPoints {
		// 获取基准响应
		baselineResp, baselineBody, err := s.getBaselineResponse(ctx, target, point)
		if err != nil {
			continue // 跳过无法获取基准响应的参数
		}

		// 执行不同类型的SQL注入测试
		vulns := []*models.Vulnerability{}
		
		fmt.Printf("[DEBUG] 正在测试参数: %s\n", point.Name)
		
		// 1. 错误注入检测
		if errorVuln := s.testErrorBasedInjection(ctx, target, point, baselineResp, baselineBody); errorVuln != nil {
			vulns = append(vulns, errorVuln)
			fmt.Printf("[DEBUG] 发现错误注入漏洞: %s\n", point.Name)
		}
		
		// 2. 布尔盲注检测
		if boolVuln := s.testBooleanBlindInjection(ctx, target, point, baselineResp, baselineBody); boolVuln != nil {
			vulns = append(vulns, boolVuln)
			fmt.Printf("[DEBUG] 发现布尔盲注漏洞: %s\n", point.Name)
		}
		
		// 3. 时间盲注检测
		if timeVuln := s.testTimeBasedInjection(ctx, target, point, baselineResp, baselineBody); timeVuln != nil {
			vulns = append(vulns, timeVuln)
			fmt.Printf("[DEBUG] 发现时间盲注漏洞: %s\n", point.Name)
		}

		// 如果发现漏洞，标记结果并添加到结果中
		if len(vulns) > 0 {
			result.IsVulnerable = true
			result.Vulnerabilities = append(result.Vulnerabilities, vulns...)
			
			// 只测试第一个发现漏洞的参数，避免过度测试
			break
		}
	}

	result.Metadata["tested_parameters"] = len(injectPoints)
	result.Metadata["detection_time"] = time.Now().Format(time.RFC3339)

	return result, nil
}

// getBaselineResponse 获取基准响应
func (s *SQLiDetector) getBaselineResponse(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint) ([]byte, []byte, error) {
	// 如果已有基准响应，直接使用
	if target.BaselineResponse != nil && target.BaselineBody != nil {
		return target.BaselineBody, target.BaselineBody, nil
	}

	// 发送原始请求获取基准响应
	resp, err := s.requestModifier.ModifyParameter(ctx, target, point, point.Value)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	helper := transport.NewResponseHelper()
	body, err := helper.ReadBody(resp)
	if err != nil {
		return nil, nil, err
	}

	return body, body, nil
}

// testErrorBasedInjection 测试错误注入
func (s *SQLiDetector) testErrorBasedInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp, baselineBody []byte) *models.Vulnerability {
	errorPayloads := s.getErrorPayloads(point.Type)
	
	for _, payload := range errorPayloads {
		resp, err := s.requestModifier.ModifyParameter(ctx, target, point, point.Value+payload.Value)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		helper := transport.NewResponseHelper()
		body, err := helper.ReadBody(resp)
		if err != nil {
			fmt.Printf("[DEBUG] 读取响应体失败: %v\n", err)
			continue
		}

		// 检查是否包含SQL错误信息
		hasErrors, foundErrors := s.responseAnalyzer.ContainsErrorPatterns(body, s.getSQLErrorPatterns())
		fmt.Printf("[DEBUG] 测试参数: %s, Payload: %s, 响应长度: %d\n", point.Name, payload.Value, len(body))
		fmt.Printf("[DEBUG] 错误检测结果: hasErrors=%t, foundErrors=%v\n", hasErrors, foundErrors)
		
		// 调试: 搜索可能的错误信息位置
		bodyStr := string(body)
		if strings.Contains(strings.ToLower(bodyStr), "error") {
			fmt.Printf("[DEBUG] 发现错误关键字，在响应中搜索上下文:\n")
			lines := strings.Split(bodyStr, "\n")
			for i, line := range lines {
				if strings.Contains(strings.ToLower(line), "error") {
					start := max(0, i-2)
					end := min(len(lines), i+3)
					for j := start; j < end; j++ {
						fmt.Printf("[DEBUG] 第%d行: %s\n", j+1, strings.TrimSpace(lines[j]))
					}
					break
				}
			}
		} else {
			fmt.Printf("[DEBUG] 响应中未找到'error'关键字，响应长度: %d\n", len(body))
			if len(body) > 200 {
				fmt.Printf("[DEBUG] 响应内容开头: %s...\n", string(body[:200]))
			}
			
			// 保存完整响应内容用于分析
			if point.Name == "title" && payload.Value == "'" {
				os.WriteFile("/tmp/sqli_response_debug.html", body, 0644)
				fmt.Printf("[DEBUG] 响应已保存到 /tmp/sqli_response_debug.html\n")
			}
		}
		
		// 检查响应状态码和内容变化
		if resp.StatusCode >= 500 {
			fmt.Printf("[DEBUG] 发现HTTP错误状态码: %d\n", resp.StatusCode)
			hasErrors = true
			foundErrors = append(foundErrors, fmt.Sprintf("HTTP %d Error", resp.StatusCode))
		}
		
		if hasErrors {
			fmt.Printf("[INFO] 发现SQL注入漏洞 - 参数: %s, 错误: %v\n", point.Name, foundErrors)
			return s.buildVulnerability(
				models.VulnSQLi,
				"SQL Error-based Injection",
				fmt.Sprintf("参数 %s 存在基于错误的SQL注入漏洞", point.Name),
				target,
				point,
				point.Value+payload.Value,
				fmt.Sprintf("SQL错误信息: %v", foundErrors),
				0.95,
			)
		}
	}

	return nil
}

// testBooleanBlindInjection 测试布尔盲注
func (s *SQLiDetector) testBooleanBlindInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp, baselineBody []byte) *models.Vulnerability {
	booleanPairs := s.getBooleanPayloadPairs(point.Type)
	
	for _, pair := range booleanPairs {
		// 测试True payload
		trueResp, err := s.requestModifier.ModifyParameter(ctx, target, point, point.Value+pair.TruePayload)
		if err != nil {
			continue
		}
		defer trueResp.Body.Close()

		helper := transport.NewResponseHelper()
		trueBody, err := helper.ReadBody(trueResp)
		if err != nil {
			continue
		}

		// 测试False payload
		falseResp, err := s.requestModifier.ModifyParameter(ctx, target, point, point.Value+pair.FalsePayload)
		if err != nil {
			continue
		}
		defer falseResp.Body.Close()

		falseBody, err := helper.ReadBody(falseResp)
		if err != nil {
			continue
		}

		// 分析响应差异
		if s.analyzeBooleanDifference(baselineBody, trueBody, falseBody) {
			return s.buildVulnerability(
				models.VulnSQLi,
				"SQL Boolean-based Blind Injection",
				fmt.Sprintf("参数 %s 存在SQL布尔盲注漏洞", point.Name),
				target,
				point,
				fmt.Sprintf("True: %s%s, False: %s%s", point.Value, pair.TruePayload, point.Value, pair.FalsePayload),
				fmt.Sprintf("布尔盲注测试成功: %s", pair.Description),
				0.85,
			)
		}
	}

	return nil
}

// testTimeBasedInjection 测试时间盲注
func (s *SQLiDetector) testTimeBasedInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp, baselineBody []byte) *models.Vulnerability {
	timePayloads := s.getTimeBasedPayloads(point.Type)
	
	for _, payload := range timePayloads {
		start := time.Now()
		resp, err := s.requestModifier.ModifyParameter(ctx, target, point, point.Value+payload.Value)
		duration := time.Since(start)
		
		if err != nil {
			continue
		}
		resp.Body.Close()

		// 如果响应时间显著增加，可能存在时间盲注
		if duration > payload.ExpectedDelay && duration > 3*time.Second {
			return s.buildVulnerability(
				models.VulnSQLi,
				"SQL Time-based Blind Injection",
				fmt.Sprintf("参数 %s 存在SQL时间盲注漏洞", point.Name),
				target,
				point,
				point.Value+payload.Value,
				fmt.Sprintf("时间延迟检测成功，响应时间: %v (预期: %v)", duration, payload.ExpectedDelay),
				0.80,
			)
		}
	}

	return nil
}

// PayloadPair 布尔payload对
type PayloadPair struct {
	TruePayload  string
	FalsePayload string
	Description  string
}

// Payload 通用payload结构
type Payload struct {
	Value         string
	Description   string
	ExpectedDelay time.Duration // 用于时间盲注
}

// getErrorPayloads 获取错误注入payload
func (s *SQLiDetector) getErrorPayloads(paramType detector.ParamType) []Payload {
	// 基础语法错误 payload
	payloads := []Payload{
		{Value: "'", Description: "单引号语法错误"},
		{Value: "''", Description: "双单引号测试"},
		{Value: "\"", Description: "双引号语法错误"},
		{Value: "\\'", Description: "转义单引号"},
		{Value: "' OR '1'='1", Description: "OR注入测试"},
		{Value: "' AND '1'='1", Description: "AND泣入测试"},
		{Value: "' UNION SELECT 1--", Description: "UNION注入测试"},
		{Value: "';--", Description: "单行注释测试"},
		{Value: "' /*", Description: "多行注释测试"},
	}

	// 数字类型参数的特定 payload
	if paramType == detector.ParamTypeNumeric {
		payloads = append(payloads, []Payload{
			{Value: " OR 1=1", Description: "数字类OR注入"},
			{Value: " AND 1=1", Description: "数字类AND注入"},
			{Value: " UNION SELECT 1", Description: "数字类UNION注入"},
			{Value: "; DROP TABLE users--", Description: "DROP TABLE测试"},
			{Value: " AND 1=CONVERT(int,(SELECT @@version))", Description: "MSSQL版本检测"},
			{Value: " AND 1=1/0", Description: "除零错误"},
		}...)
	}

	return payloads
}

// getBooleanPayloadPairs 获取布尔盲注payload对
func (s *SQLiDetector) getBooleanPayloadPairs(paramType detector.ParamType) []PayloadPair {
	pairs := []PayloadPair{
		{
			TruePayload:  "' AND '1'='1",
			FalsePayload: "' AND '1'='2",
			Description:  "字符型单引号布尔注入",
		},
		{
			TruePayload:  "\" AND \"1\"=\"1",
			FalsePayload: "\" AND \"1\"=\"2",
			Description:  "字符型双引号布尔注入",
		},
	}

	if paramType == detector.ParamTypeNumeric {
		pairs = append(pairs, []PayloadPair{
			{
				TruePayload:  " AND 1=1",
				FalsePayload: " AND 1=2",
				Description:  "数字型AND布尔注入",
			},
			{
				TruePayload:  " OR 1=1",
				FalsePayload: " OR 1=2",
				Description:  "数字型OR布尔注入",
			},
		}...)
	}

	return pairs
}

// getTimeBasedPayloads 获取时间盲注payload
func (s *SQLiDetector) getTimeBasedPayloads(paramType detector.ParamType) []Payload {
	payloads := []Payload{
		{
			Value:         "' AND SLEEP(5)--",
			Description:   "MySQL时间延迟",
			ExpectedDelay: 5 * time.Second,
		},
		{
			Value:         "' AND (SELECT SLEEP(5))--",
			Description:   "MySQL SELECT时间延迟",
			ExpectedDelay: 5 * time.Second,
		},
		{
			Value:         "'; WAITFOR DELAY '00:00:05'--",
			Description:   "MSSQL时间延迟",
			ExpectedDelay: 5 * time.Second,
		},
	}

	if paramType == detector.ParamTypeNumeric {
		payloads = append(payloads, []Payload{
			{
				Value:         " AND SLEEP(5)",
				Description:   "数字型MySQL时间延迟",
				ExpectedDelay: 5 * time.Second,
			},
			{
				Value:         "; WAITFOR DELAY '00:00:05'",
				Description:   "数字型MSSQL时间延迟",
				ExpectedDelay: 5 * time.Second,
			},
		}...)
	}

	return payloads
}

// getSQLErrorPatterns 获取SQL错误模式
func (s *SQLiDetector) getSQLErrorPatterns() []string {
	return []string{
		// MySQL
		"you have an error in your sql syntax",
		"warning: mysql",
		"mysql_fetch_array()",
		"mysql_fetch_assoc()",
		"mysql_fetch_row()",
		"mysql_num_rows()",
		"mysql error",
		"supplied argument is not a valid mysql",
		"column count doesn't match value count",
		
		// PostgreSQL
		"postgresql query failed",
		"warning: pg_",
		"invalid query result",
		"pg_query() expects",
		
		// MSSQL
		"microsoft ole db provider",
		"odbc sql server driver",
		"microsoft sql native client",
		"sqlstate",
		"sqlexception",
		
		// Oracle
		"ora-01756",
		"ora-00936",
		"ora-00942",
		"oracle error",
		"oracle driver",
		
		// SQLite
		"sqlite_error",
		"sqlite3.operationalerror",
		"no such column",
		"sql error or missing database",
		
		// Generic SQL
		"sql syntax",
		"syntax error",
		"unexpected token",
		"unclosed quotation mark",
		"invalid column name",
		"must declare the scalar variable",
		"operand should contain 1 column(s)",
		"the used select statements have different number of columns",
		"table doesn't exist",
		"unknown column",
		"ambiguous column name",
		"division by zero error encountered",
		"data type mismatch",
		
		// PHP/Web specific
		"warning:",
		"fatal error",
		"call to undefined function",
		"cannot execute statement",
	}
}

// analyzeBooleanDifference 分析布尔差异
func (s *SQLiDetector) analyzeBooleanDifference(baseline, trueResp, falseResp []byte) bool {
	// 计算相似度
	trueSim, _ := s.responseAnalyzer.AnalyzeDifference(baseline, trueResp)
	falseSim, _ := s.responseAnalyzer.AnalyzeDifference(baseline, falseResp)

	// True payload应该与基准相似，False payload应该不同
	if trueSim > 0.95 && falseSim < 0.85 {
		return true
	}

	// 检查响应长度差异
	baselineLen := len(baseline)
	trueLen := len(trueResp)
	falseLen := len(falseResp)

	// True响应长度接近基准，False响应长度差异较大
	if abs(trueLen-baselineLen) < 100 && abs(falseLen-baselineLen) > 500 {
		return true
	}

	// 检查状态码等其他指标...

	return false
}

// buildVulnerability 构建漏洞对象
func (s *SQLiDetector) buildVulnerability(
	vulnType models.VulnType,
	title, description string,
	target *detector.ScanTarget,
	point detector.InjectPoint,
	payload, evidence string,
	confidence float64,
) *models.Vulnerability {
	return models.NewVulnerabilityBuilder().
		WithType(vulnType).
		WithCategory(models.CategoryInjection).
		WithSeverity(models.SeverityHigh).
		WithTitle(title).
		WithDescription(description).
		WithURL(target.URL.String()).
		WithMethod(target.Method).
		WithParameter(point.Name, point.Position).
		WithPayload(payload).
		WithEvidence(evidence).
		WithConfidence(confidence).
		WithPlugin(s.Name()).
		WithCWE("CWE-89").
		WithCVSS(9.0).
		WithSolution("使用预编译语句（参数化查询）、输入验证和最小权限原则").
		WithReferences([]string{
			"https://owasp.org/www-community/attacks/SQL_Injection",
			"https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html",
		}).
		Build()
}

// abs 返回绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// isPrintableText 检查内容是否为可读文本
func isPrintableText(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	
	// 检查前100个字节
	checkLen := 100
	if len(data) < checkLen {
		checkLen = len(data)
	}
	
	printableCount := 0
	for i := 0; i < checkLen; i++ {
		b := data[i]
		// ASCII可打印字符范围是32-126，加上换行符
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	
	// 如果80%以上是可打印字符，认为是文本
	return float64(printableCount)/float64(checkLen) > 0.8
}

// max 返回两个数的最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min 返回两个数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}