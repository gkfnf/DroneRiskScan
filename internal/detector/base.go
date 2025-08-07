package detector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dronesec/droneriskscan/internal/transport"
	"github.com/dronesec/droneriskscan/pkg/models"
)

// Plugin 漏洞检测插件接口
type Plugin interface {
	// 基础信息
	Name() string
	Description() string
	Author() string
	Version() string
	
	// 分类信息
	Type() PluginType
	Category() models.Category
	Severity() models.Severity
	
	// 执行检测
	Execute(ctx context.Context, target *ScanTarget) (*DetectionResult, error)
	
	// 配置检查
	IsEnabled() bool
	SetEnabled(enabled bool)
}

// PluginType 插件类型
type PluginType string

const (
	PluginTypeActive  PluginType = "active"   // 主动扫描
	PluginTypePassive PluginType = "passive"  // 被动扫描
	PluginTypeHybrid  PluginType = "hybrid"   // 混合模式
)

// ScanTarget 扫描目标
type ScanTarget struct {
	URL        *url.URL
	Method     string
	Headers    map[string]string
	Body       string
	Parameters map[string][]string // 支持多值参数
	Cookies    map[string]string
	Metadata   map[string]interface{}
	
	// 基准响应（用于比对）
	BaselineResponse *http.Response
	BaselineBody     []byte
}

// DetectionResult 检测结果
type DetectionResult struct {
	IsVulnerable    bool
	Vulnerabilities []*models.Vulnerability
	Evidence        []Evidence
	Requests        []*http.Request
	Responses       []*http.Response
	Metadata        map[string]interface{}
}

// Evidence 证据信息
type Evidence struct {
	Type        EvidenceType `json:"type"`
	Description string       `json:"description"`
	Data        interface{}  `json:"data"`
	Confidence  float64      `json:"confidence"`
}

// EvidenceType 证据类型
type EvidenceType string

const (
	EvidenceTypeResponse     EvidenceType = "response"
	EvidenceTypeError        EvidenceType = "error"
	EvidenceTypeDifference   EvidenceType = "difference"
	EvidenceTypePattern      EvidenceType = "pattern"
	EvidenceTypeTiming       EvidenceType = "timing"
	EvidenceTypeStatusCode   EvidenceType = "status_code"
)

// BasePlugin 基础插件实现，提供通用功能
type BasePlugin struct {
	name        string
	description string
	author      string
	version     string
	pluginType  PluginType
	category    models.Category
	severity    models.Severity
	enabled     bool
	
	// HTTP客户端
	httpClient transport.HTTPClient
	
	// 配置选项
	options map[string]interface{}
}

// NewBasePlugin 创建基础插件
func NewBasePlugin(name string, pluginType PluginType, category models.Category, severity models.Severity) *BasePlugin {
	return &BasePlugin{
		name:       name,
		pluginType: pluginType,
		category:   category,
		severity:   severity,
		enabled:    true,
		version:    "1.0.0",
		options:    make(map[string]interface{}),
	}
}

// Name 返回插件名称
func (bp *BasePlugin) Name() string {
	return bp.name
}

// Description 返回插件描述
func (bp *BasePlugin) Description() string {
	return bp.description
}

// Author 返回插件作者
func (bp *BasePlugin) Author() string {
	return bp.author
}

// Version 返回插件版本
func (bp *BasePlugin) Version() string {
	return bp.version
}

// Type 返回插件类型
func (bp *BasePlugin) Type() PluginType {
	return bp.pluginType
}

// Category 返回漏洞类别
func (bp *BasePlugin) Category() models.Category {
	return bp.category
}

// Severity 返回严重程度
func (bp *BasePlugin) Severity() models.Severity {
	return bp.severity
}

// IsEnabled 检查插件是否启用
func (bp *BasePlugin) IsEnabled() bool {
	return bp.enabled
}

// SetEnabled 设置插件启用状态
func (bp *BasePlugin) SetEnabled(enabled bool) {
	bp.enabled = enabled
}

// SetDescription 设置描述
func (bp *BasePlugin) SetDescription(description string) {
	bp.description = description
}

// SetAuthor 设置作者
func (bp *BasePlugin) SetAuthor(author string) {
	bp.author = author
}

// SetVersion 设置版本
func (bp *BasePlugin) SetVersion(version string) {
	bp.version = version
}

// SetHTTPClient 设置HTTP客户端
func (bp *BasePlugin) SetHTTPClient(client transport.HTTPClient) {
	bp.httpClient = client
}

// GetHTTPClient 获取HTTP客户端
func (bp *BasePlugin) GetHTTPClient() transport.HTTPClient {
	return bp.httpClient
}

// SetOption 设置选项
func (bp *BasePlugin) SetOption(key string, value interface{}) {
	bp.options[key] = value
}

// GetOption 获取选项
func (bp *BasePlugin) GetOption(key string) (interface{}, bool) {
	value, exists := bp.options[key]
	return value, exists
}

// Execute 基础执行方法，子类需要重写
func (bp *BasePlugin) Execute(ctx context.Context, target *ScanTarget) (*DetectionResult, error) {
	return &DetectionResult{
		IsVulnerable:    false,
		Vulnerabilities: []*models.Vulnerability{},
		Evidence:        []Evidence{},
		Metadata:        make(map[string]interface{}),
	}, nil
}

// ParameterExtractor 参数提取器
type ParameterExtractor struct{}

// NewParameterExtractor 创建参数提取器
func NewParameterExtractor() *ParameterExtractor {
	return &ParameterExtractor{}
}

// ExtractParameters 从目标中提取参数
func (pe *ParameterExtractor) ExtractParameters(target *ScanTarget) []InjectPoint {
	var points []InjectPoint
	
	// 提取GET参数
	for name, values := range target.Parameters {
		for _, value := range values {
			points = append(points, InjectPoint{
				Name:     name,
				Value:    value,
				Position: models.PositionGET,
				Type:     pe.inferParameterType(value),
			})
		}
	}
	
	// 提取POST参数（如果是表单数据）
	if target.Method == "POST" && strings.Contains(target.Headers["Content-Type"], "application/x-www-form-urlencoded") {
		if formValues, err := url.ParseQuery(target.Body); err == nil {
			for name, values := range formValues {
				for _, value := range values {
					points = append(points, InjectPoint{
						Name:     name,
						Value:    value,
						Position: models.PositionPOST,
						Type:     pe.inferParameterType(value),
					})
				}
			}
		}
	}
	
	// 提取Cookie参数（排除认证相关Cookie）
	authCookies := []string{"PHPSESSID", "JSESSIONID", "ASP.NET_SessionId", "security_level", "_token", "csrf_token"}
	for name, value := range target.Cookies {
		// 跳过认证相关Cookie
		isAuthCookie := false
		for _, authCookie := range authCookies {
			if strings.EqualFold(name, authCookie) {
				isAuthCookie = true
				break
			}
		}
		
		if !isAuthCookie {
			points = append(points, InjectPoint{
				Name:     name,
				Value:    value,
				Position: models.PositionCOOKIE,
				Type:     pe.inferParameterType(value),
			})
		}
	}
	
	// 提取头部参数（某些特定头部）
	vulnerableHeaders := []string{"X-Forwarded-For", "X-Real-IP", "User-Agent", "Referer"}
	for _, header := range vulnerableHeaders {
		if value, exists := target.Headers[header]; exists {
			points = append(points, InjectPoint{
				Name:     header,
				Value:    value,
				Position: models.PositionHEADER,
				Type:     ParamTypeString,
			})
		}
	}
	
	return points
}

// InjectPoint 注入点
type InjectPoint struct {
	Name     string
	Value    string
	Position models.Position
	Type     ParamType
}

// ParamType 参数类型
type ParamType string

const (
	ParamTypeString  ParamType = "string"
	ParamTypeNumeric ParamType = "numeric"
	ParamTypeBoolean ParamType = "boolean"
	ParamTypeEmail   ParamType = "email"
	ParamTypeURL     ParamType = "url"
)

// inferParameterType 推断参数类型
func (pe *ParameterExtractor) inferParameterType(value string) ParamType {
	// 数字类型
	if strings.TrimSpace(value) != "" {
		if matched := regexp.MustCompile(`^\d+$`).MatchString(value); matched {
			return ParamTypeNumeric
		}
	}
	
	// 布尔类型
	lowerValue := strings.ToLower(value)
	if lowerValue == "true" || lowerValue == "false" || lowerValue == "1" || lowerValue == "0" {
		return ParamTypeBoolean
	}
	
	// 邮箱类型
	if strings.Contains(value, "@") && strings.Contains(value, ".") {
		return ParamTypeEmail
	}
	
	// URL类型
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return ParamTypeURL
	}
	
	// 默认字符串类型
	return ParamTypeString
}

// RequestModifier 请求修改器
type RequestModifier struct {
	httpClient transport.HTTPClient
	sessionCookies []*http.Cookie
}

// NewRequestModifier 创建请求修改器
func NewRequestModifier(httpClient transport.HTTPClient) *RequestModifier {
	return &RequestModifier{
		httpClient: httpClient,
		sessionCookies: make([]*http.Cookie, 0),
	}
}

// SetSessionCookies 设置会话Cookie
func (rm *RequestModifier) SetSessionCookies(cookies []*http.Cookie) {
	rm.sessionCookies = cookies
}

// ModifyParameter 修改参数并发送请求
func (rm *RequestModifier) ModifyParameter(ctx context.Context, target *ScanTarget, point InjectPoint, payload string) (*http.Response, error) {
	// 根据参数位置修改请求
	switch point.Position {
	case models.PositionGET:
		return rm.modifyGETParameter(ctx, target, point.Name, payload)
	case models.PositionPOST:
		return rm.modifyPOSTParameter(ctx, target, point.Name, payload)
	case models.PositionHEADER:
		return rm.modifyHeaderParameter(ctx, target, point.Name, payload)
	case models.PositionCOOKIE:
		return rm.modifyCookieParameter(ctx, target, point.Name, payload)
	default:
		return nil, fmt.Errorf("不支持的参数位置: %s", point.Position)
	}
}

// modifyGETParameter 修改GET参数
func (rm *RequestModifier) modifyGETParameter(ctx context.Context, target *ScanTarget, paramName, payload string) (*http.Response, error) {
	// 解析URL
	u := *target.URL
	query := u.Query()
	
	// 修改参数值
	query.Set(paramName, payload)
	u.RawQuery = query.Encode()
	
	// 创建新请求
	finalURL := u.String()
	fmt.Printf("[DEBUG] 发送请求URL: %s\n", finalURL)
	req, err := http.NewRequestWithContext(ctx, target.Method, finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 添加头部
	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}
	
	// 添加目标Cookie（但跳过会被会话Cookie覆盖的）
	sessionCookieNames := make(map[string]bool)
	for _, cookie := range rm.sessionCookies {
		sessionCookieNames[cookie.Name] = true
	}
	
	for k, v := range target.Cookies {
		// 如果会话Cookie中有同名的，跳过目标Cookie
		if !sessionCookieNames[k] {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}
	
	// 添加会话Cookie（这些优先级更高）
	for _, cookie := range rm.sessionCookies {
		req.AddCookie(cookie)
	}
	
	// 调试：打印完整的请求信息
	fmt.Printf("[DEBUG] 完整请求信息:\n")
	fmt.Printf("[DEBUG]   URL: %s\n", req.URL.String())
	fmt.Printf("[DEBUG]   Method: %s\n", req.Method)
	fmt.Printf("[DEBUG]   Headers: %+v\n", req.Header)
	cookieHeader := req.Header.Get("Cookie")
	if cookieHeader != "" {
		fmt.Printf("[DEBUG]   Cookie Header: %s\n", cookieHeader)
	} else {
		fmt.Printf("[DEBUG]   Cookie Header: <EMPTY>\n")
	}
	
	return rm.httpClient.Do(req)
}

// modifyPOSTParameter 修改POST参数
func (rm *RequestModifier) modifyPOSTParameter(ctx context.Context, target *ScanTarget, paramName, payload string) (*http.Response, error) {
	// 解析表单数据
	formValues, err := url.ParseQuery(target.Body)
	if err != nil {
		return nil, fmt.Errorf("解析表单数据失败: %w", err)
	}
	
	// 修改参数值
	formValues.Set(paramName, payload)
	newBody := formValues.Encode()
	
	// 创建新请求
	req, err := http.NewRequestWithContext(ctx, target.Method, target.URL.String(), strings.NewReader(newBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 设置Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// 添加其他头部
	for k, v := range target.Headers {
		if k != "Content-Type" && k != "Content-Length" {
			req.Header.Set(k, v)
		}
	}
	
	// 添加目标Cookie
	for k, v := range target.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	
	// 添加会话Cookie
	for _, cookie := range rm.sessionCookies {
		req.AddCookie(cookie)
	}
	
	return rm.httpClient.Do(req)
}

// modifyHeaderParameter 修改头部参数
func (rm *RequestModifier) modifyHeaderParameter(ctx context.Context, target *ScanTarget, headerName, payload string) (*http.Response, error) {
	var body io.Reader
	if target.Body != "" {
		body = strings.NewReader(target.Body)
	}
	
	req, err := http.NewRequestWithContext(ctx, target.Method, target.URL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 添加头部
	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}
	
	// 修改指定头部
	req.Header.Set(headerName, payload)
	
	// 添加目标Cookie
	for k, v := range target.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	
	// 添加会话Cookie
	for _, cookie := range rm.sessionCookies {
		req.AddCookie(cookie)
	}
	
	return rm.httpClient.Do(req)
}

// modifyCookieParameter 修改Cookie参数
func (rm *RequestModifier) modifyCookieParameter(ctx context.Context, target *ScanTarget, cookieName, payload string) (*http.Response, error) {
	var body io.Reader
	if target.Body != "" {
		body = strings.NewReader(target.Body)
	}
	
	req, err := http.NewRequestWithContext(ctx, target.Method, target.URL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 添加头部
	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}
	
	// 添加目标Cookie
	for k, v := range target.Cookies {
		if k == cookieName {
			req.AddCookie(&http.Cookie{Name: k, Value: payload})
		} else {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}
	
	// 添加会话Cookie
	for _, cookie := range rm.sessionCookies {
		req.AddCookie(cookie)
	}
	
	return rm.httpClient.Do(req)
}

// ResponseAnalyzer 响应分析器
type ResponseAnalyzer struct{}

// NewResponseAnalyzer 创建响应分析器
func NewResponseAnalyzer() *ResponseAnalyzer {
	return &ResponseAnalyzer{}
}

// AnalyzeDifference 分析响应差异
func (ra *ResponseAnalyzer) AnalyzeDifference(baseline, test []byte) (float64, map[string]interface{}) {
	analysis := make(map[string]interface{})
	
	// 长度差异
	lenDiff := abs(len(baseline) - len(test))
	analysis["length_difference"] = lenDiff
	analysis["baseline_length"] = len(baseline)
	analysis["test_length"] = len(test)
	
	// 字符串相似度（简单实现）
	similarity := calculateStringSimilarity(string(baseline), string(test))
	analysis["similarity"] = similarity
	
	// 状态码变化等其他分析...
	
	return similarity, analysis
}

// ContainsErrorPatterns 检查是否包含错误模式
func (ra *ResponseAnalyzer) ContainsErrorPatterns(body []byte, patterns []string) (bool, []string) {
	content := strings.ToLower(string(body))
	var foundPatterns []string
	
	for _, pattern := range patterns {
		if strings.Contains(content, strings.ToLower(pattern)) {
			foundPatterns = append(foundPatterns, pattern)
		}
	}
	
	return len(foundPatterns) > 0, foundPatterns
}

// 辅助函数
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 简单的字符串相似度计算
func calculateStringSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 1.0
	}
	
	// 这里可以使用更复杂的算法如Levenshtein距离
	// 简单实现：计算公共字符比例
	common := 0
	for i := 0; i < min(len(s1), len(s2)); i++ {
		if s1[i] == s2[i] {
			common++
		}
	}
	
	return float64(common) / float64(maxLen)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}