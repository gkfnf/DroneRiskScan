package injection

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/internal/detector"
	"github.com/dronesec/droneriskscan/internal/transport"
	"github.com/dronesec/droneriskscan/pkg/models"
)

// EnhancedSQLiDetector 增强的SQL注入检测器 - 基于sqlmap的检测逻辑
type EnhancedSQLiDetector struct {
	*detector.BasePlugin
	requestModifier  *detector.RequestModifier
	responseAnalyzer *detector.ResponseAnalyzer
	paramExtractor   *detector.ParameterExtractor
}

// NewEnhancedSQLiDetector 创建增强的SQL注入检测器
func NewEnhancedSQLiDetector(httpClient transport.HTTPClient) *EnhancedSQLiDetector {
	base := detector.NewBasePlugin(
		"enhanced-sqli-detector",
		detector.PluginTypeActive,
		models.CategoryInjection,
		models.SeverityHigh,
	)
	
	base.SetDescription("增强SQL注入检测器 - 基于sqlmap检测逻辑")
	base.SetAuthor("DroneRiskScan Team (Enhanced)")
	base.SetVersion("2.0.0")
	base.SetHTTPClient(httpClient)

	return &EnhancedSQLiDetector{
		BasePlugin:       base,
		requestModifier:  detector.NewRequestModifier(httpClient),
		responseAnalyzer: detector.NewResponseAnalyzer(),
		paramExtractor:   detector.NewParameterExtractor(),
	}
}

// SetSessionCookies 设置会话Cookie
func (e *EnhancedSQLiDetector) SetSessionCookies(cookies []*http.Cookie) {
	e.requestModifier.SetSessionCookies(cookies)
}

// SQLiTest SQL注入测试用例
type SQLiTest struct {
	Title         string
	PayloadType   string // error, boolean, time, union
	Risk          int    // 1-3
	Level         int    // 1-5
	Clause        string // WHERE, ORDER BY, HAVING等
	Payload       string
	Where         int    // 注入位置类型
	Vector        string // 注入向量
	Request       *SQLiRequest
	Response      *SQLiResponse
	Comparison    string // 比较方式
	DBMS          []string // 支持的数据库类型
	Confidence    float64  // 置信度
}

// SQLiRequest 请求配置
type SQLiRequest struct {
	Prefix    string
	Suffix    string
	Comment   string
	Columns   int
	Char      string
}

// SQLiResponse 响应配置
type SQLiResponse struct {
	Union     *UnionTest
	Boolean   *BooleanTest
	Error     *ErrorTest
	Time      *TimeTest
}

// UnionTest UNION注入测试
type UnionTest struct {
	FieldsCount int
	Char        string
	Template    string
}

// BooleanTest 布尔注入测试
type BooleanTest struct {
	TrueRegex      string
	FalseRegex     string
	TrueCondition  string
	FalseCondition string
	Description    string
}

// ErrorTest 错误注入测试
type ErrorTest struct {
	Regex       []string
	Template    string
	Extractor   string
}

// TimeTest 时间注入测试
type TimeTest struct {
	Delay       int
	Template    string
	Threshold   float64
}

// Execute 执行SQL注入检测
func (e *EnhancedSQLiDetector) Execute(ctx context.Context, target *detector.ScanTarget) (*detector.DetectionResult, error) {
	result := &detector.DetectionResult{
		IsVulnerable:    false,
		Vulnerabilities: []*models.Vulnerability{},
		Evidence:        []detector.Evidence{},
		Metadata:        make(map[string]interface{}),
	}

	// 提取所有注入点
	injectPoints := e.paramExtractor.ExtractParameters(target)
	if len(injectPoints) == 0 {
		result.Metadata["message"] = "未发现可注入参数"
		return result, nil
	}
	
	fmt.Printf("[INFO] Enhanced SQL注入检测器找到 %d 个注入点\n", len(injectPoints))

	// 对每个注入点进行测试
	for _, point := range injectPoints {
		fmt.Printf("[INFO] 正在测试参数: %s (类型: %s, 位置: %s)\n", point.Name, point.Type, point.Position)
		
		// 获取基准响应
		baselineResp, baselineBody, err := e.getBaselineResponse(ctx, target, point)
		if err != nil {
			fmt.Printf("[ERROR] 无法获取基准响应: %v\n", err)
			continue
		}
		
		fmt.Printf("[DEBUG] 基准响应长度: %d bytes, 状态码: %d\n", len(baselineBody), baselineResp.StatusCode)
		
		// 执行不同类型的SQL注入测试
		vulns := []*models.Vulnerability{}
		
		// 1. 错误注入检测 (最直接的方法)
		if errorVuln := e.testErrorBasedInjection(ctx, target, point, baselineResp, baselineBody); errorVuln != nil {
			vulns = append(vulns, errorVuln)
			fmt.Printf("[SUCCESS] 发现错误注入漏洞: %s\n", point.Name)
		}
		
		// 2. 布尔盲注检测 (更准确的比较逻辑)
		if boolVuln := e.testBooleanBlindInjection(ctx, target, point, baselineResp, baselineBody); boolVuln != nil {
			vulns = append(vulns, boolVuln)
			fmt.Printf("[SUCCESS] 发现布尔盲注漏洞: %s\n", point.Name)
		}
		
		// 3. UNION注入检测 (新增)
		if unionVuln := e.testUnionBasedInjection(ctx, target, point, baselineResp, baselineBody); unionVuln != nil {
			vulns = append(vulns, unionVuln)
			fmt.Printf("[SUCCESS] 发现UNION注入漏洞: %s\n", point.Name)
		}
		
		// 4. 时间盲注检测 (最后检测，因为耗时)
		if timeVuln := e.testTimeBasedInjection(ctx, target, point, baselineResp, baselineBody); timeVuln != nil {
			vulns = append(vulns, timeVuln)
			fmt.Printf("[SUCCESS] 发现时间盲注漏洞: %s\n", point.Name)
		}

		// 如果发现漏洞，标记结果并添加到结果中
		if len(vulns) > 0 {
			result.IsVulnerable = true
			result.Vulnerabilities = append(result.Vulnerabilities, vulns...)
			
			// 对于高置信度的漏洞，可以继续测试其他参数
			// 对于低置信度的漏洞，只测试第一个发现的
			break
		}
	}

	result.Metadata["tested_parameters"] = len(injectPoints)
	result.Metadata["detection_time"] = time.Now().Format(time.RFC3339)
	result.Metadata["vulnerabilities_found"] = len(result.Vulnerabilities)

	return result, nil
}

// getBaselineResponse 获取基准响应
func (e *EnhancedSQLiDetector) getBaselineResponse(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint) (*http.Response, []byte, error) {
	// 发送原始请求获取基准响应
	resp, err := e.requestModifier.ModifyParameter(ctx, target, point, point.Value)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	helper := transport.NewResponseHelper()
	body, err := helper.ReadBody(resp)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}

// testErrorBasedInjection 测试错误注入 - 基于sqlmap的方法
func (e *EnhancedSQLiDetector) testErrorBasedInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp *http.Response, baselineBody []byte) *models.Vulnerability {
	errorTests := e.getErrorBasedTests()
	
	for _, test := range errorTests {
		// 构造payload
		payload := e.buildPayload(point, test.Payload)
		
		fmt.Printf("[DEBUG] 测试错误注入payload: %s\n", payload)
		
		// 发送请求
		resp, err := e.requestModifier.ModifyParameter(ctx, target, point, payload)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		helper := transport.NewResponseHelper()
		body, err := helper.ReadBody(resp)
		if err != nil {
			continue
		}

		// 检查HTTP状态码变化
		if resp.StatusCode >= 500 && baselineResp.StatusCode < 500 {
			fmt.Printf("[FOUND] HTTP错误状态码变化: %d -> %d\n", baselineResp.StatusCode, resp.StatusCode)
			return e.buildVulnerability(
				models.VulnSQLi,
				"SQL Error-based Injection (HTTP Status)",
				fmt.Sprintf("参数 %s 存在基于错误的SQL注入漏洞 (HTTP状态码变化)", point.Name),
				target, point, payload,
				fmt.Sprintf("HTTP状态码从 %d 变为 %d", baselineResp.StatusCode, resp.StatusCode),
				0.90,
			)
		}

		// 检查SQL错误信息
		bodyStr := string(body)
		for _, pattern := range e.getSQLErrorPatterns() {
			if matched, _ := regexp.MatchString("(?i)"+pattern, bodyStr); matched {
				fmt.Printf("[FOUND] SQL错误模式匹配: %s\n", pattern)
				return e.buildVulnerability(
					models.VulnSQLi,
					"SQL Error-based Injection (Error Message)",
					fmt.Sprintf("参数 %s 存在基于错误的SQL注入漏洞", point.Name),
					target, point, payload,
					fmt.Sprintf("检测到SQL错误模式: %s", pattern),
					0.95,
				)
			}
		}
		
		// 检查响应长度显著变化 (可能是错误页面)
		baselineLen := len(baselineBody)
		currentLen := len(body)
		lengthDiff := absInt(currentLen - baselineLen)
		
		if lengthDiff > 1000 && lengthDiff > baselineLen/4 { // 长度差异超过25%且超过1000字节
			fmt.Printf("[FOUND] 响应长度显著变化: %d -> %d (差异: %d)\n", baselineLen, currentLen, lengthDiff)
			
			// 进一步检查是否包含错误相关内容
			if e.containsErrorIndicators(bodyStr) {
				return e.buildVulnerability(
					models.VulnSQLi,
					"SQL Error-based Injection (Response Change)",
					fmt.Sprintf("参数 %s 存在基于错误的SQL注入漏洞 (响应变化)", point.Name),
					target, point, payload,
					fmt.Sprintf("响应长度显著变化且包含错误指示符"),
					0.75,
				)
			}
		}
	}

	return nil
}

// testBooleanBlindInjection 测试布尔盲注 - 改进版本
func (e *EnhancedSQLiDetector) testBooleanBlindInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp *http.Response, baselineBody []byte) *models.Vulnerability {
	booleanTests := e.getBooleanBlindTests()
	
	for _, test := range booleanTests {
		// 测试True条件
		truePayload := e.buildPayload(point, test.TrueCondition)
		trueResp, err := e.requestModifier.ModifyParameter(ctx, target, point, truePayload)
		if err != nil {
			continue
		}
		defer trueResp.Body.Close()

		helper := transport.NewResponseHelper()
		trueBody, err := helper.ReadBody(trueResp)
		if err != nil {
			continue
		}

		// 测试False条件  
		falsePayload := e.buildPayload(point, test.FalseCondition)
		falseResp, err := e.requestModifier.ModifyParameter(ctx, target, point, falsePayload)
		if err != nil {
			continue
		}
		defer falseResp.Body.Close()

		falseBody, err := helper.ReadBody(falseResp)
		if err != nil {
			continue
		}
		
		fmt.Printf("[DEBUG] 布尔测试: True长度=%d, False长度=%d, Baseline长度=%d\n", 
			len(trueBody), len(falseBody), len(baselineBody))

		// 分析响应差异 - 使用更精确的方法
		if e.analyzeBooleanDifference(baselineBody, trueBody, falseBody) {
			fmt.Printf("[FOUND] 布尔盲注差异检测成功\n")
			return e.buildVulnerability(
				models.VulnSQLi,
				"SQL Boolean-based Blind Injection",
				fmt.Sprintf("参数 %s 存在SQL布尔盲注漏洞", point.Name),
				target, point,
				fmt.Sprintf("True: %s, False: %s", truePayload, falsePayload),
				fmt.Sprintf("布尔盲注测试成功: %s", test.Description),
				0.85,
			)
		}
	}

	return nil
}

// testUnionBasedInjection 测试UNION注入 - 新增功能
func (e *EnhancedSQLiDetector) testUnionBasedInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp *http.Response, baselineBody []byte) *models.Vulnerability {
	// UNION注入需要先确定列数
	for colCount := 1; colCount <= 10; colCount++ {
		// 构造ORDER BY测试确定列数
		orderByPayload := e.buildPayload(point, fmt.Sprintf("' ORDER BY %d--", colCount))
		
		resp, err := e.requestModifier.ModifyParameter(ctx, target, point, orderByPayload)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		helper := transport.NewResponseHelper()
		body, err := helper.ReadBody(resp)
		if err != nil {
			continue
		}

		// 如果ORDER BY出错，说明列数不够
		if resp.StatusCode >= 500 || e.containsErrorIndicators(string(body)) {
			if colCount > 1 {
				// 尝试UNION SELECT
				unionPayload := e.buildUnionPayload(point, colCount-1)
				unionResp, err := e.requestModifier.ModifyParameter(ctx, target, point, unionPayload)
				if err != nil {
					continue
				}
				defer unionResp.Body.Close()

				unionBody, err := helper.ReadBody(unionResp)
				if err != nil {
					continue
				}

				// 检查UNION标识符
				if strings.Contains(string(unionBody), "UNION_TEST_") {
					return e.buildVulnerability(
						models.VulnSQLi,
						"SQL UNION-based Injection",
						fmt.Sprintf("参数 %s 存在UNION注入漏洞", point.Name),
						target, point, unionPayload,
						fmt.Sprintf("成功执行UNION查询，列数: %d", colCount-1),
						0.95,
					)
				}
			}
			break
		}
	}

	return nil
}

// testTimeBasedInjection 测试时间盲注 - 改进版本
func (e *EnhancedSQLiDetector) testTimeBasedInjection(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, baselineResp *http.Response, baselineBody []byte) *models.Vulnerability {
	timeTests := e.getTimeBasedTests()
	
	// 先测量基准响应时间
	baselineTime := e.measureResponseTime(ctx, target, point, point.Value)
	if baselineTime < 0 {
		return nil
	}
	
	fmt.Printf("[DEBUG] 基准响应时间: %v\n", baselineTime)
	
	for _, test := range timeTests {
		payload := e.buildPayload(point, test.Template)
		
		// 测量延迟响应时间
		delayTime := e.measureResponseTime(ctx, target, point, payload)
		if delayTime < 0 {
			continue
		}
		
		fmt.Printf("[DEBUG] 延迟响应时间: %v (预期延迟: %ds)\n", delayTime, test.Delay)
		
		// 判断是否存在明显的时间延迟
		expectedDelay := time.Duration(test.Delay) * time.Second
		actualDelay := delayTime - baselineTime
		
		if actualDelay >= expectedDelay*8/10 && actualDelay >= 3*time.Second {
			return e.buildVulnerability(
				models.VulnSQLi,
				"SQL Time-based Blind Injection",
				fmt.Sprintf("参数 %s 存在SQL时间盲注漏洞", point.Name),
				target, point, payload,
				fmt.Sprintf("时间延迟检测成功，基准: %v, 延迟: %v, 差异: %v", 
					baselineTime, delayTime, actualDelay),
				0.80,
			)
		}
	}

	return nil
}

// 辅助方法

func (e *EnhancedSQLiDetector) buildPayload(point detector.InjectPoint, template string) string {
	// 根据参数类型构造不同的payload
	switch point.Type {
	case detector.ParamTypeNumeric:
		// 数字型参数
		if strings.HasPrefix(template, "'") {
			return point.Value + strings.TrimPrefix(template, "'")
		}
		return point.Value + template
	default:
		// 字符型参数
		return point.Value + template
	}
}

func (e *EnhancedSQLiDetector) buildUnionPayload(point detector.InjectPoint, colCount int) string {
	// 构造UNION SELECT payload
	columns := make([]string, colCount)
	for i := 0; i < colCount; i++ {
		if i == 0 {
			columns[i] = "'UNION_TEST_" + fmt.Sprintf("%d", i) + "'"
		} else {
			columns[i] = fmt.Sprintf("'COL_%d'", i)
		}
	}
	
	unionPart := "UNION SELECT " + strings.Join(columns, ",")
	
	switch point.Type {
	case detector.ParamTypeNumeric:
		return point.Value + " " + unionPart + "--"
	default:
		return point.Value + "' " + unionPart + "--"
	}
}

func (e *EnhancedSQLiDetector) measureResponseTime(ctx context.Context, target *detector.ScanTarget, point detector.InjectPoint, payload string) time.Duration {
	start := time.Now()
	resp, err := e.requestModifier.ModifyParameter(ctx, target, point, payload)
	duration := time.Since(start)
	
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	
	// 读取响应体以确保完整的响应时间
	helper := transport.NewResponseHelper()
	_, err = helper.ReadBody(resp)
	if err != nil {
		return -1
	}
	
	return duration
}

func (e *EnhancedSQLiDetector) containsErrorIndicators(body string) bool {
	errorIndicators := []string{
		"error", "warning", "exception", "fatal", "mysql", "postgres", "oracle", 
		"mssql", "sqlite", "syntax", "invalid", "unexpected", "failed",
	}
	
	bodyLower := strings.ToLower(body)
	for _, indicator := range errorIndicators {
		if strings.Contains(bodyLower, indicator) {
			return true
		}
	}
	return false
}

func (e *EnhancedSQLiDetector) analyzeBooleanDifference(baseline, trueResp, falseResp []byte) bool {
	baselineLen := len(baseline)
	trueLen := len(trueResp)
	falseLen := len(falseResp)
	
	// 1. 基本长度差异检查
	trueDiff := absInt(trueLen - baselineLen)
	falseDiff := absInt(falseLen - baselineLen)
	
	// True应该与基准相似，False应该有差异
	if trueDiff < 100 && falseDiff > 500 {
		return true
	}
	
	// 2. True和False之间的差异
	tfDiff := absInt(trueLen - falseLen)
	if tfDiff > 200 {
		return true
	}
	
	// 3. 内容相似度分析
	trueSim, _ := e.responseAnalyzer.AnalyzeDifference(baseline, trueResp)
	falseSim, _ := e.responseAnalyzer.AnalyzeDifference(baseline, falseResp)
	
	// True应该与基准更相似
	if trueSim > 0.95 && falseSim < 0.85 {
		return true
	}
	
	// 4. 文本内容差异分析
	if e.analyzeTextDifferences(string(trueResp), string(falseResp)) {
		return true
	}
	
	return false
}

func (e *EnhancedSQLiDetector) analyzeTextDifferences(trueBody, falseBody string) bool {
	// 检查关键词差异
	trueWords := strings.Fields(strings.ToLower(trueBody))
	falseWords := strings.Fields(strings.ToLower(falseBody))
	
	trueWordCount := make(map[string]int)
	falseWordCount := make(map[string]int)
	
	for _, word := range trueWords {
		trueWordCount[word]++
	}
	
	for _, word := range falseWords {
		falseWordCount[word]++
	}
	
	// 计算词汇差异
	diffCount := 0
	for word := range trueWordCount {
		if trueWordCount[word] != falseWordCount[word] {
			diffCount++
		}
	}
	
	for word := range falseWordCount {
		if falseWordCount[word] != trueWordCount[word] {
			diffCount++
		}
	}
	
	// 如果词汇差异超过10%，认为存在差异
	totalWords := len(trueWords) + len(falseWords)
	if totalWords > 0 && float64(diffCount)/float64(totalWords) > 0.1 {
		return true
	}
	
	return false
}

// 测试用例定义

func (e *EnhancedSQLiDetector) getErrorBasedTests() []SQLiTest {
	return []SQLiTest{
		{
			Title:       "MySQL Error-based",
			PayloadType: "error", 
			Payload:     "'",
			Confidence:  0.9,
			DBMS:        []string{"MySQL"},
		},
		{
			Title:       "Generic Error",
			PayloadType: "error",
			Payload:     "' AND (SELECT * FROM (SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)a)--",
			Confidence:  0.95,
			DBMS:        []string{"MySQL"},
		},
		{
			Title:       "MSSQL Error",
			PayloadType: "error",
			Payload:     "' AND 1=CONVERT(int,(SELECT @@version))--",
			Confidence:  0.9,
			DBMS:        []string{"MSSQL"},
		},
	}
}

func (e *EnhancedSQLiDetector) getBooleanBlindTests() []BooleanTest {
	return []BooleanTest{
		{
			TrueCondition:  "' AND '1'='1",
			FalseCondition: "' AND '1'='2", 
			Description:    "字符型单引号布尔注入",
		},
		{
			TrueCondition:  "\" AND \"1\"=\"1",
			FalseCondition: "\" AND \"1\"=\"2",
			Description:    "字符型双引号布尔注入",
		},
		{
			TrueCondition:  " AND 1=1",
			FalseCondition: " AND 1=2",
			Description:    "数字型AND布尔注入",
		},
		{
			TrueCondition:  "' AND 'a'='a",
			FalseCondition: "' AND 'a'='b",
			Description:    "字符型字母比较布尔注入",
		},
	}
}

func (e *EnhancedSQLiDetector) getTimeBasedTests() []TimeTest {
	return []TimeTest{
		{
			Delay:    5,
			Template: "' AND SLEEP(5)--",
		},
		{
			Delay:    5,
			Template: "' AND (SELECT SLEEP(5))--", 
		},
		{
			Delay:    5,
			Template: "'; WAITFOR DELAY '00:00:05'--",
		},
		{
			Delay:    5,
			Template: " AND SLEEP(5)",
		},
	}
}

func (e *EnhancedSQLiDetector) getSQLErrorPatterns() []string {
	return []string{
		// MySQL
		"you have an error in your sql syntax",
		"warning: mysql",
		"mysql_fetch",
		"mysql_num_rows",
		"mysql error",
		"supplied argument is not a valid mysql",
		"column count doesn't match value count",
		"operand should contain 1 column",
		"illegal mix of collations",
		"invalid use of group function",
		
		// PostgreSQL  
		"postgresql query failed",
		"warning: pg_",
		"invalid query result", 
		"pg_query\\(\\) expects",
		"pg_exec\\(\\) expects",
		
		// MSSQL
		"microsoft ole db provider",
		"odbc sql server driver", 
		"microsoft sql native client",
		"sqlstate",
		"sqlexception",
		"system\\.data\\.sqlclient\\.sqlexception",
		"unclosed quotation mark after the character string",
		"incorrect syntax near",
		
		// Oracle
		"ora-01756",
		"ora-00936", 
		"ora-00942",
		"oracle error",
		"oracle driver",
		"quoted string not properly terminated",
		
		// SQLite
		"sqlite_error",
		"sqlite3\\.operationalerror",
		"no such column",
		"sql error or missing database",
		"sqlite3.database error",
		
		// Generic
		"sql syntax",
		"syntax error", 
		"unexpected token",
		"invalid column name",
		"must declare the scalar variable",
		"table doesn't exist",
		"unknown column",
		"ambiguous column name", 
		"division by zero error encountered",
		"data type mismatch",
		"conversion failed",
		"invalid object reference",
	}
}

func (e *EnhancedSQLiDetector) buildVulnerability(
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
		WithPlugin(e.Name()).
		WithCWE("CWE-89").
		WithCVSS(9.0).
		WithSolution("使用预编译语句（参数化查询）、输入验证和最小权限原则").
		WithReferences([]string{
			"https://owasp.org/www-community/attacks/SQL_Injection",
			"https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html",
		}).
		Build()
}

// absInt 返回绝对值
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}