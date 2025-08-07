package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dronesec/droneriskscan/internal/auth"
	"github.com/dronesec/droneriskscan/internal/transport"
)

// Crawler 网页爬虫
type Crawler struct {
	httpClient     transport.HTTPClient
	sessionManager *auth.SessionManager
	config         *CrawlerConfig
	visited        map[string]bool
	visitedMutex   sync.RWMutex
	results        []*CrawlResult
	resultsMutex   sync.RWMutex
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	MaxDepth       int           // 最大爬取深度
	MaxPages       int           // 最大页面数量
	RequestTimeout time.Duration // 请求超时时间
	Delay          time.Duration // 请求间隔
	UserAgent      string        // User Agent
	Verbose        bool          // 详细输出
	AllowedDomains []string      // 允许的域名
	ExcludeExts    []string      // 排除的文件扩展名
	FollowRedirect bool          // 是否跟随重定向
}

// CrawlResult 爬取结果
type CrawlResult struct {
	URL          string            `json:"url"`
	Title        string            `json:"title"`
	StatusCode   int               `json:"status_code"`
	ContentType  string            `json:"content_type"`
	ContentSize  int64             `json:"content_size"`
	ResponseTime time.Duration     `json:"response_time"`
	Forms        []*FormInfo       `json:"forms"`
	Links        []string          `json:"links"`
	Inputs       []*InputInfo      `json:"inputs"`
	Headers      map[string]string `json:"headers"`
	Cookies      []*http.Cookie    `json:"cookies"`
	Depth        int               `json:"depth"`
	Timestamp    time.Time         `json:"timestamp"`
	
	// 功能分析结果
	FunctionType string   `json:"function_type"` // 功能类型（登录、搜索、上传等）
	RiskLevel    string   `json:"risk_level"`    // 风险等级
	TestPlugins  []string `json:"test_plugins"`  // 推荐的测试插件
}

// FormInfo 表单信息
type FormInfo struct {
	Action     string       `json:"action"`
	Method     string       `json:"method"`
	EncType    string       `json:"enctype"`
	Inputs     []*InputInfo `json:"inputs"`
	HasUpload  bool         `json:"has_upload"`
	HasHidden  bool         `json:"has_hidden"`
	IsLogin    bool         `json:"is_login"`
	IsSearch   bool         `json:"is_search"`
}

// InputInfo 输入字段信息
type InputInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Placeholder string `json:"placeholder"`
	Required    bool   `json:"required"`
	MaxLength   int    `json:"maxlength"`
}

// NewCrawler 创建新的爬虫
func NewCrawler(httpClient transport.HTTPClient, sessionManager *auth.SessionManager, config *CrawlerConfig) *Crawler {
	if config == nil {
		config = DefaultCrawlerConfig()
	}
	
	return &Crawler{
		httpClient:     httpClient,
		sessionManager: sessionManager,
		config:         config,
		visited:        make(map[string]bool),
		results:        make([]*CrawlResult, 0),
	}
}

// DefaultCrawlerConfig 返回默认配置
func DefaultCrawlerConfig() *CrawlerConfig {
	return &CrawlerConfig{
		MaxDepth:       3,
		MaxPages:       100,
		RequestTimeout: 30 * time.Second,
		Delay:          100 * time.Millisecond,
		UserAgent:      "DroneRiskScan/1.0 Web Crawler",
		Verbose:        false,
		AllowedDomains: []string{},
		ExcludeExts:    []string{"jpg", "jpeg", "png", "gif", "css", "js", "ico", "svg", "woff", "ttf", "pdf"},
		FollowRedirect: true,
	}
}

// Crawl 开始爬取
func (c *Crawler) Crawl(ctx context.Context, startURL string) ([]*CrawlResult, error) {
	if c.config.Verbose {
		fmt.Printf("[INFO] 开始爬取: %s (最大深度: %d, 最大页面: %d)\n", 
			startURL, c.config.MaxDepth, c.config.MaxPages)
	}
	
	// 解析起始URL
	baseURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("解析起始URL失败: %w", err)
	}
	
	// 设置允许的域名
	if len(c.config.AllowedDomains) == 0 {
		c.config.AllowedDomains = []string{baseURL.Host}
	}
	
	// 开始爬取
	err = c.crawlURL(ctx, startURL, 0)
	if err != nil {
		return nil, err
	}
	
	// 分析所有爬取结果的功能
	c.analyzeAllFunctions()
	
	if c.config.Verbose {
		fmt.Printf("[INFO] 爬取完成，共发现 %d 个页面\n", len(c.results))
	}
	
	return c.results, nil
}

// crawlURL 爬取单个URL
func (c *Crawler) crawlURL(ctx context.Context, targetURL string, depth int) error {
	// 检查深度限制
	if depth > c.config.MaxDepth {
		return nil
	}
	
	// 检查页面数量限制
	c.resultsMutex.RLock()
	pageCount := len(c.results)
	c.resultsMutex.RUnlock()
	
	if pageCount >= c.config.MaxPages {
		return nil
	}
	
	// 检查是否已访问
	c.visitedMutex.Lock()
	if c.visited[targetURL] {
		c.visitedMutex.Unlock()
		return nil
	}
	c.visited[targetURL] = true
	c.visitedMutex.Unlock()
	
	// 检查URL是否允许
	if !c.isAllowedURL(targetURL) {
		return nil
	}
	
	if c.config.Verbose {
		fmt.Printf("[CRAWL] 深度 %d: %s\n", depth, targetURL)
	}
	
	// 添加延迟
	if c.config.Delay > 0 {
		time.Sleep(c.config.Delay)
	}
	
	// 发送请求
	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("User-Agent", c.config.UserAgent)
	
	// 应用会话认证
	if c.sessionManager != nil && c.sessionManager.IsLoggedIn() {
		c.sessionManager.ApplyAuth(req)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.config.Verbose {
			fmt.Printf("[ERROR] 请求失败: %s - %v\n", targetURL, err)
		}
		return nil // 继续爬取其他页面
	}
	defer resp.Body.Close()
	
	responseTime := time.Since(startTime)
	
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	
	// 解析页面内容
	result := &CrawlResult{
		URL:          targetURL,
		StatusCode:   resp.StatusCode,
		ContentType:  resp.Header.Get("Content-Type"),
		ContentSize:  int64(len(body)),
		ResponseTime: responseTime,
		Headers:      make(map[string]string),
		Cookies:      resp.Cookies(),
		Depth:        depth,
		Timestamp:    time.Now(),
	}
	
	// 复制响应头
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}
	
	// 只处理HTML内容
	if strings.Contains(result.ContentType, "text/html") {
		c.parseHTMLContent(result, string(body))
		
		// 提取链接并继续爬取
		for _, link := range result.Links {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				absoluteURL := c.resolveURL(targetURL, link)
				if absoluteURL != "" {
					c.crawlURL(ctx, absoluteURL, depth+1)
				}
			}
		}
	}
	
	// 保存结果
	c.resultsMutex.Lock()
	c.results = append(c.results, result)
	c.resultsMutex.Unlock()
	
	return nil
}

// parseHTMLContent 解析HTML内容
func (c *Crawler) parseHTMLContent(result *CrawlResult, htmlContent string) {
	// 提取页面标题
	titleRe := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	if matches := titleRe.FindStringSubmatch(htmlContent); len(matches) > 1 {
		result.Title = strings.TrimSpace(matches[1])
	}
	
	// 提取链接
	result.Links = c.extractLinks(htmlContent)
	
	// 提取表单
	result.Forms = c.extractForms(htmlContent)
	
	// 提取所有输入字段
	result.Inputs = c.extractInputs(htmlContent)
}

// extractLinks 提取链接
func (c *Crawler) extractLinks(htmlContent string) []string {
	var links []string
	linkRe := regexp.MustCompile(`<a[^>]+href=["']([^"']+)["']`)
	matches := linkRe.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			href := strings.TrimSpace(match[1])
			if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "javascript:") {
				links = append(links, href)
			}
		}
	}
	
	return links
}

// extractForms 提取表单
func (c *Crawler) extractForms(htmlContent string) []*FormInfo {
	var forms []*FormInfo
	formRe := regexp.MustCompile(`<form[^>]*>(.*?)</form>`)
	matches := formRe.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			formHTML := match[0]
			formContent := match[1]
			
			form := &FormInfo{
				Method:  "GET", // 默认值
				EncType: "application/x-www-form-urlencoded", // 默认值
			}
			
			// 提取action属性
			actionRe := regexp.MustCompile(`action=["']([^"']+)["']`)
			if actionMatch := actionRe.FindStringSubmatch(formHTML); len(actionMatch) > 1 {
				form.Action = actionMatch[1]
			}
			
			// 提取method属性
			methodRe := regexp.MustCompile(`method=["']([^"']+)["']`)
			if methodMatch := methodRe.FindStringSubmatch(formHTML); len(methodMatch) > 1 {
				form.Method = strings.ToUpper(methodMatch[1])
			}
			
			// 提取enctype属性
			enctypeRe := regexp.MustCompile(`enctype=["']([^"']+)["']`)
			if enctypeMatch := enctypeRe.FindStringSubmatch(formHTML); len(enctypeMatch) > 1 {
				form.EncType = enctypeMatch[1]
			}
			
			// 提取输入字段
			form.Inputs = c.extractInputsFromHTML(formContent)
			
			// 分析表单特征
			c.analyzeFormType(form, formContent)
			
			forms = append(forms, form)
		}
	}
	
	return forms
}

// extractInputs 提取所有输入字段
func (c *Crawler) extractInputs(htmlContent string) []*InputInfo {
	return c.extractInputsFromHTML(htmlContent)
}

// extractInputsFromHTML 从HTML中提取输入字段
func (c *Crawler) extractInputsFromHTML(htmlContent string) []*InputInfo {
	var inputs []*InputInfo
	
	// 提取input标签
	inputRe := regexp.MustCompile(`<input[^>]*>`)
	inputMatches := inputRe.FindAllString(htmlContent, -1)
	
	for _, inputHTML := range inputMatches {
		input := &InputInfo{}
		
		// 提取属性
		if nameMatch := regexp.MustCompile(`name=["']([^"']+)["']`).FindStringSubmatch(inputHTML); len(nameMatch) > 1 {
			input.Name = nameMatch[1]
		}
		if typeMatch := regexp.MustCompile(`type=["']([^"']+)["']`).FindStringSubmatch(inputHTML); len(typeMatch) > 1 {
			input.Type = typeMatch[1]
		} else {
			input.Type = "text" // 默认类型
		}
		if valueMatch := regexp.MustCompile(`value=["']([^"']*)["']`).FindStringSubmatch(inputHTML); len(valueMatch) > 1 {
			input.Value = valueMatch[1]
		}
		if placeholderMatch := regexp.MustCompile(`placeholder=["']([^"']+)["']`).FindStringSubmatch(inputHTML); len(placeholderMatch) > 1 {
			input.Placeholder = placeholderMatch[1]
		}
		input.Required = strings.Contains(inputHTML, "required")
		
		inputs = append(inputs, input)
	}
	
	// 提取textarea标签
	textareaRe := regexp.MustCompile(`<textarea[^>]*name=["']([^"']+)["'][^>]*>(.*?)</textarea>`)
	textareaMatches := textareaRe.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range textareaMatches {
		if len(match) > 1 {
			input := &InputInfo{
				Name:  match[1],
				Type:  "textarea",
				Value: strings.TrimSpace(match[2]),
			}
			inputs = append(inputs, input)
		}
	}
	
	// 提取select标签
	selectRe := regexp.MustCompile(`<select[^>]*name=["']([^"']+)["'][^>]*>(.*?)</select>`)
	selectMatches := selectRe.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range selectMatches {
		if len(match) > 1 {
			input := &InputInfo{
				Name: match[1],
				Type: "select",
			}
			inputs = append(inputs, input)
		}
	}
	
	return inputs
}

// analyzeFormType 分析表单类型
func (c *Crawler) analyzeFormType(form *FormInfo, formContent string) {
	lowerContent := strings.ToLower(formContent)
	
	// 检查文件上传
	if strings.Contains(lowerContent, `type="file"`) || 
	   strings.Contains(form.EncType, "multipart/form-data") {
		form.HasUpload = true
	}
	
	// 检查隐藏字段
	if strings.Contains(lowerContent, `type="hidden"`) {
		form.HasHidden = true
	}
	
	// 检查登录表单
	if (strings.Contains(lowerContent, "password") || strings.Contains(lowerContent, "login")) &&
	   (strings.Contains(lowerContent, "username") || strings.Contains(lowerContent, "user") || strings.Contains(lowerContent, "email")) {
		form.IsLogin = true
	}
	
	// 检查搜索表单
	if strings.Contains(lowerContent, "search") || strings.Contains(lowerContent, "query") ||
	   strings.Contains(lowerContent, "keyword") {
		form.IsSearch = true
	}
}

// analyzeAllFunctions 分析所有页面的功能
func (c *Crawler) analyzeAllFunctions() {
	for _, result := range c.results {
		c.analyzeFunctionType(result)
		c.recommendTestPlugins(result)
	}
}

// analyzeFunctionType 分析页面功能类型
func (c *Crawler) analyzeFunctionType(result *CrawlResult) {
	url := strings.ToLower(result.URL)
	title := strings.ToLower(result.Title)
	
	// 基于URL和标题的功能识别
	functionPatterns := map[string][]string{
		"login":    {"login", "signin", "auth", "登录", "sign_in"},
		"search":   {"search", "query", "find", "搜索", "查询"},
		"upload":   {"upload", "file", "上传", "文件"},
		"admin":    {"admin", "管理", "administration", "管理员"},
		"profile":  {"profile", "account", "user", "个人", "用户"},
		"comment":  {"comment", "feedback", "message", "评论", "留言"},
		"blog":     {"blog", "post", "article", "博客", "文章"},
		"register": {"register", "signup", "注册", "sign_up"},
		"contact":  {"contact", "联系", "about"},
		"cart":     {"cart", "shop", "order", "购物", "订单"},
	}
	
	for funcType, patterns := range functionPatterns {
		for _, pattern := range patterns {
			if strings.Contains(url, pattern) || strings.Contains(title, pattern) {
				result.FunctionType = funcType
				break
			}
		}
		if result.FunctionType != "" {
			break
		}
	}
	
	// 基于表单的功能识别
	for _, form := range result.Forms {
		if form.IsLogin {
			result.FunctionType = "login"
		} else if form.IsSearch {
			result.FunctionType = "search"
		} else if form.HasUpload {
			result.FunctionType = "upload"
		}
	}
	
	// 默认功能类型
	if result.FunctionType == "" {
		if len(result.Forms) > 0 || len(result.Inputs) > 0 {
			result.FunctionType = "form"
		} else {
			result.FunctionType = "static"
		}
	}
}

// recommendTestPlugins 推荐测试插件
func (c *Crawler) recommendTestPlugins(result *CrawlResult) {
	var plugins []string
	var riskLevel = "low"
	
	// 基于功能类型推荐插件
	switch result.FunctionType {
	case "login":
		plugins = append(plugins, "sqli-detector", "auth-bypass", "brute-force")
		riskLevel = "high"
	case "search":
		plugins = append(plugins, "sqli-detector", "xss-detector", "nosql-injection")
		riskLevel = "high"
	case "upload":
		plugins = append(plugins, "file-upload", "path-traversal", "malware-upload")
		riskLevel = "critical"
	case "admin":
		plugins = append(plugins, "sqli-detector", "xss-detector", "auth-bypass", "privilege-escalation")
		riskLevel = "critical"
	case "comment", "blog":
		plugins = append(plugins, "xss-detector", "sqli-detector", "csrf-detector")
		riskLevel = "medium"
	case "form":
		plugins = append(plugins, "sqli-detector", "xss-detector", "csrf-detector")
		riskLevel = "medium"
	}
	
	// 基于输入字段推荐插件
	for _, input := range result.Inputs {
		lowerName := strings.ToLower(input.Name)
		if strings.Contains(lowerName, "file") {
			plugins = append(plugins, "file-upload")
			riskLevel = "high"
		}
		if strings.Contains(lowerName, "url") || strings.Contains(lowerName, "link") {
			plugins = append(plugins, "ssrf-detector", "open-redirect")
		}
		if strings.Contains(lowerName, "cmd") || strings.Contains(lowerName, "command") {
			plugins = append(plugins, "command-injection")
			riskLevel = "critical"
		}
	}
	
	// 如果有参数，添加基础注入检测
	if len(result.Inputs) > 0 || strings.Contains(result.URL, "?") {
		if !contains(plugins, "sqli-detector") {
			plugins = append(plugins, "sqli-detector")
		}
		if !contains(plugins, "xss-detector") {
			plugins = append(plugins, "xss-detector")
		}
	}
	
	result.TestPlugins = removeDuplicates(plugins)
	result.RiskLevel = riskLevel
}

// isAllowedURL 检查URL是否允许爬取
func (c *Crawler) isAllowedURL(targetURL string) bool {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	
	// 检查域名
	allowed := false
	for _, domain := range c.config.AllowedDomains {
		if parsedURL.Host == domain {
			allowed = true
			break
		}
	}
	if !allowed {
		return false
	}
	
	// 检查文件扩展名
	path := strings.ToLower(parsedURL.Path)
	for _, ext := range c.config.ExcludeExts {
		if strings.HasSuffix(path, "."+ext) {
			return false
		}
	}
	
	return true
}

// resolveURL 解析相对URL为绝对URL
func (c *Crawler) resolveURL(baseURL, href string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	resolved, err := base.Parse(href)
	if err != nil {
		return ""
	}
	
	return resolved.String()
}

// GetResults 获取爬取结果
func (c *Crawler) GetResults() []*CrawlResult {
	c.resultsMutex.RLock()
	defer c.resultsMutex.RUnlock()
	
	// 返回副本
	results := make([]*CrawlResult, len(c.results))
	copy(results, c.results)
	return results
}

// 辅助函数
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}