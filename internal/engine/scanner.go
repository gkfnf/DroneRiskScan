package engine

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dronesec/droneriskscan/internal/auth"
	"github.com/dronesec/droneriskscan/internal/crawler"
	"github.com/dronesec/droneriskscan/internal/detector"
	"github.com/dronesec/droneriskscan/internal/detector/injection"
	"github.com/dronesec/droneriskscan/internal/reporter"
	"github.com/dronesec/droneriskscan/internal/scheduler"
	"github.com/dronesec/droneriskscan/internal/transport"
	"github.com/dronesec/droneriskscan/pkg/models"
)

// Scanner 漏洞扫描引擎
type Scanner struct {
	httpClient     transport.HTTPClient
	scheduler      scheduler.TaskScheduler
	reporter       reporter.ReportGenerator
	plugins        map[string]detector.Plugin
	config         *ScannerConfig
	sessionManager *auth.SessionManager
	crawler        *crawler.Crawler
	mutex          sync.RWMutex
}

// ScannerConfig 扫描器配置
type ScannerConfig struct {
	MaxConcurrency   int
	RequestTimeout   time.Duration
	MaxRedirects     int
	UserAgent        string
	EnabledPlugins   []string
	DisabledPlugins  []string
	RiskLevels       []models.Severity
	ReportFormats    []string
	Verbose          bool
	Debug            bool
	
	// 认证配置
	AuthCredentials  *auth.Credentials
	
	// 爬虫配置
	EnableCrawler    bool
	MaxCrawlDepth    int
	MaxCrawlPages    int
}

// NewScanner 创建新的扫描器实例
func NewScanner(config *ScannerConfig) (*Scanner, error) {
	if config == nil {
		config = DefaultScannerConfig()
	}

	// 创建HTTP客户端
	clientOptions := &transport.ClientOptions{
		Timeout:         config.RequestTimeout,
		MaxRedirects:    config.MaxRedirects,
		UserAgent:       config.UserAgent,
		InsecureSkipTLS: true,
	}
	httpClient := transport.NewHTTPClient(clientOptions)

	// 创建任务调度器
	taskScheduler := scheduler.NewTaskScheduler(&scheduler.Config{
		MaxWorkers:    config.MaxConcurrency,
		QueueSize:     1000,
		RetryAttempts: 3,
		RetryDelay:    time.Second,
	})

	// 创建报告生成器
	reportGenerator := reporter.NewReportGenerator(&reporter.Config{
		Formats: config.ReportFormats,
		Debug:   config.Debug,
	})

	scanner := &Scanner{
		httpClient: httpClient,
		scheduler:  taskScheduler,
		reporter:   reportGenerator,
		plugins:    make(map[string]detector.Plugin),
		config:     config,
	}

	// 初始化会话管理器
	if config.AuthCredentials != nil {
		scanner.sessionManager = auth.NewSessionManager(httpClient, config.AuthCredentials)
	}

	// 初始化爬虫
	if config.EnableCrawler {
		crawlerConfig := &crawler.CrawlerConfig{
			MaxDepth:       config.MaxCrawlDepth,
			MaxPages:       config.MaxCrawlPages,
			RequestTimeout: config.RequestTimeout,
			Delay:          100 * time.Millisecond,
			UserAgent:      config.UserAgent,
			Verbose:        config.Verbose,
		}
		scanner.crawler = crawler.NewCrawler(httpClient, scanner.sessionManager, crawlerConfig)
	}

	// 注册默认插件
	if err := scanner.registerDefaultPlugins(); err != nil {
		return nil, fmt.Errorf("注册默认插件失败: %w", err)
	}

	return scanner, nil
}

// DefaultScannerConfig 返回默认配置
func DefaultScannerConfig() *ScannerConfig {
	return &ScannerConfig{
		MaxConcurrency: 10,
		RequestTimeout: 30 * time.Second,
		MaxRedirects:   5,
		UserAgent:      "DroneRiskScan/1.0 Security Scanner",
		RiskLevels: []models.Severity{
			models.SeverityLow,
			models.SeverityMedium,
			models.SeverityHigh,
			models.SeverityCritical,
		},
		ReportFormats: []string{"json"},
		Verbose:       false,
		Debug:         false,
		EnableCrawler: true,
		MaxCrawlDepth: 2,
		MaxCrawlPages: 50,
	}
}

// registerDefaultPlugins 注册默认插件
func (s *Scanner) registerDefaultPlugins() error {
	// 注册SQL注入检测器
	sqliDetector := injection.NewEnhancedSQLiDetector(s.httpClient)
	if err := s.RegisterPlugin(sqliDetector); err != nil {
		return fmt.Errorf("注册SQL注入检测器失败: %w", err)
	}

	// 这里可以注册更多检测器...
	// xssDetector := xss.NewXSSDetector(s.httpClient)
	// s.RegisterPlugin(xssDetector)

	return nil
}

// RegisterPlugin 注册检测插件
func (s *Scanner) RegisterPlugin(plugin detector.Plugin) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if plugin == nil {
		return fmt.Errorf("插件不能为空")
	}

	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("插件名称不能为空")
	}

	// 检查是否被禁用
	for _, disabled := range s.config.DisabledPlugins {
		if disabled == name {
			if s.config.Debug {
				fmt.Printf("[DEBUG] 插件 %s 被禁用，跳过注册\n", name)
			}
			return nil
		}
	}

	// 检查启用列表（如果配置了）
	if len(s.config.EnabledPlugins) > 0 {
		enabled := false
		for _, enabledName := range s.config.EnabledPlugins {
			if enabledName == name {
				enabled = true
				break
			}
		}
		if !enabled {
			if s.config.Debug {
				fmt.Printf("[DEBUG] 插件 %s 不在启用列表中，跳过注册\n", name)
			}
			return nil
		}
	}

	// 检查严重程度过滤
	pluginSeverity := plugin.Severity()
	severityAllowed := false
	for _, allowedSeverity := range s.config.RiskLevels {
		if allowedSeverity == pluginSeverity {
			severityAllowed = true
			break
		}
	}
	if !severityAllowed {
		if s.config.Debug {
			fmt.Printf("[DEBUG] 插件 %s 严重程度 %s 不在允许列表中\n", name, pluginSeverity)
		}
		return nil
	}

	s.plugins[name] = plugin

	if s.config.Verbose {
		fmt.Printf("[INFO] 已注册插件: %s (类型: %s, 严重程度: %s)\n",
			name, plugin.Type(), plugin.Severity())
	}

	return nil
}

// ScanURL 扫描单个URL
func (s *Scanner) ScanURL(ctx context.Context, targetURL string) (*models.ScanResult, error) {
	return s.ScanURLs(ctx, []string{targetURL})
}

// ScanURLs 扫描多个URL
func (s *Scanner) ScanURLs(ctx context.Context, targetURLs []string) (*models.ScanResult, error) {
	if len(targetURLs) == 0 {
		return nil, fmt.Errorf("目标URL列表不能为空")
	}

	// 创建扫描结果
	scanID := generateScanID()
	result := models.NewScanResult(scanID)
	result.SetRunning()

	if s.config.Verbose {
		fmt.Printf("[INFO] 开始扫描 %d 个目标\n", len(targetURLs))
	}
	
	// 直接使用原始目标
	allTargets := targetURLs

	// 启动任务调度器
	if err := s.scheduler.Start(ctx); err != nil {
		return nil, fmt.Errorf("启动任务调度器失败: %w", err)
	}
	defer s.scheduler.Stop()

	// 为每个目标创建扫描任务
	var wg sync.WaitGroup
	resultChan := make(chan *models.Vulnerability, 100)
	errorChan := make(chan error, len(allTargets))

	// 启动结果收集器
	go s.collectResults(result, resultChan, errorChan)

	for _, targetURL := range allTargets {
		wg.Add(1)
		
		// 创建扫描任务
		task := &scheduler.Task{
			ID:       fmt.Sprintf("scan_%s_%d", extractHostFromURL(targetURL), time.Now().UnixNano()),
			Type:     scheduler.TaskTypeScan,
			Priority: scheduler.PriorityNormal,
			Payload: map[string]interface{}{
				"url":    targetURL,
				"result": result,
			},
			CreatedAt: time.Now(),
		}

		// 提交任务
		s.scheduler.Submit(task, func(ctx context.Context, task *scheduler.Task) error {
			defer wg.Done()
			
			targetURL := task.Payload["url"].(string)
			scanResult := task.Payload["result"].(*models.ScanResult)
			
			return s.scanSingleTarget(ctx, targetURL, scanResult, resultChan, errorChan)
		})
	}

	// 等待所有任务完成
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// 等待结果收集完成
	time.Sleep(100 * time.Millisecond)

	result.SetCompleted()

	if s.config.Verbose {
		fmt.Printf("[INFO] 扫描完成，共发现 %d 个漏洞\n", result.GetVulnerabilityCount())
	}

	return result, nil
}

// scanSingleTarget 扫描单个目标
func (s *Scanner) scanSingleTarget(ctx context.Context, targetURL string, result *models.ScanResult, resultChan chan<- *models.Vulnerability, errorChan chan<- error) error {
	// 解析URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		errorChan <- fmt.Errorf("解析URL失败 %s: %w", targetURL, err)
		return err
	}

	// 创建目标结果记录
	targetResult := &models.TargetResult{
		URL:    targetURL,
		Status: models.TargetStatusScanning,
	}
	result.AddTarget(targetResult)

	// 发送初始请求获取基准响应
	startTime := time.Now()
	resp, err := s.httpClient.Get(targetURL)
	responseTime := time.Since(startTime)

	targetResult.ResponseTime = responseTime

	if err != nil {
		targetResult.Status = models.TargetStatusFailed
		targetResult.Errors = []string{err.Error()}
		result.UpdateTarget(targetURL, models.TargetStatusFailed)
		errorChan <- fmt.Errorf("请求目标失败 %s: %w", targetURL, err)
		return err
	}
	defer resp.Body.Close()

	// 更新目标信息
	targetResult.StatusCode = resp.StatusCode
	targetResult.ContentType = resp.Header.Get("Content-Type")
	if resp.ContentLength >= 0 {
		targetResult.ContentSize = resp.ContentLength
	}

	// 读取响应体
	helper := transport.NewResponseHelper()
	body, err := helper.ReadBody(resp)
	if err != nil {
		errorChan <- fmt.Errorf("读取响应体失败 %s: %w", targetURL, err)
		return err
	}

	// 创建扫描目标
	scanTarget := &detector.ScanTarget{
		URL:        parsedURL,
		Method:     "GET",
		Headers:    make(map[string]string),
		Parameters: make(map[string][]string),
		Cookies:    make(map[string]string),
		Metadata:   make(map[string]interface{}),
		
		BaselineResponse: resp,
		BaselineBody:     body,
	}

	// 提取参数
	if parsedURL.RawQuery != "" {
		for key, values := range parsedURL.Query() {
			scanTarget.Parameters[key] = values
		}
	}

	// 提取Cookie
	for _, cookie := range resp.Cookies() {
		scanTarget.Cookies[cookie.Name] = cookie.Value
	}

	// 复制响应头到请求头（某些需要的头部）
	for key, values := range resp.Header {
		if len(values) > 0 && shouldCopyHeader(key) {
			scanTarget.Headers[key] = values[0]
		}
	}

	// 执行所有启用的检测插件
	s.mutex.RLock()
	plugins := make([]detector.Plugin, 0, len(s.plugins))
	for _, plugin := range s.plugins {
		if plugin.IsEnabled() {
			plugins = append(plugins, plugin)
		}
	}
	s.mutex.RUnlock()

	for _, plugin := range plugins {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if s.config.Debug {
			fmt.Printf("[DEBUG] 执行插件 %s 对目标 %s\n", plugin.Name(), targetURL)
		}

		// 如果存在会话管理器，为插件设置会话Cookie
		if s.sessionManager != nil && s.sessionManager.IsLoggedIn() {
			if sqliPlugin, ok := plugin.(*injection.SQLiDetector); ok {
				sqliPlugin.SetSessionCookies(s.sessionManager.GetCookies())
			}
		}

		// 执行检测
		detectionResult, err := plugin.Execute(ctx, scanTarget)
		if err != nil {
			if s.config.Debug {
				fmt.Printf("[DEBUG] 插件 %s 执行失败: %v\n", plugin.Name(), err)
			}
			continue
		}

		// 处理检测结果
		if detectionResult != nil && detectionResult.IsVulnerable {
			for _, vuln := range detectionResult.Vulnerabilities {
				// 设置发现时间
				vuln.Timestamp = time.Now()
				
				// 通过通道发送漏洞
				select {
				case resultChan <- vuln:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}

	targetResult.Status = models.TargetStatusCompleted
	result.UpdateTarget(targetURL, models.TargetStatusCompleted)

	return nil
}

// collectResults 收集扫描结果
func (s *Scanner) collectResults(result *models.ScanResult, resultChan <-chan *models.Vulnerability, errorChan <-chan error) {
	for {
		select {
		case vuln, ok := <-resultChan:
			if !ok {
				return
			}
			if vuln != nil {
				result.AddVulnerability(vuln)
				if s.config.Verbose {
					fmt.Printf("[FOUND] %s: %s (参数: %s)\n", 
						vuln.Severity.String(), vuln.Title, vuln.Parameter)
				}
			}
		case err, ok := <-errorChan:
			if !ok {
				continue
			}
			if err != nil && s.config.Debug {
				fmt.Printf("[ERROR] %v\n", err)
			}
		}
	}
}

// GenerateReport 生成扫描报告
func (s *Scanner) GenerateReport(result *models.ScanResult, format string, outputPath string) error {
	return s.reporter.GenerateReport(result, format, outputPath)
}

// GetPlugins 获取已注册的插件列表
func (s *Scanner) GetPlugins() []detector.Plugin {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	plugins := make([]detector.Plugin, 0, len(s.plugins))
	for _, plugin := range s.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// Login 执行登录认证
func (s *Scanner) Login(ctx context.Context) error {
	if s.sessionManager == nil {
		return fmt.Errorf("未配置认证凭据")
	}

	if s.config.Verbose {
		fmt.Printf("[INFO] 正在执行登录认证...\n")
	}

	err := s.sessionManager.Login(ctx)
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	if s.config.Verbose {
		fmt.Printf("[INFO] 登录成功，会话ID: %s\n", s.sessionManager.GetSessionID())
	}

	return nil
}

// IsAuthenticated 检查是否已认证
func (s *Scanner) IsAuthenticated() bool {
	if s.sessionManager == nil {
		return false
	}
	return s.sessionManager.IsLoggedIn()
}

// GetSessionCookies 获取会话Cookie
func (s *Scanner) GetSessionCookies() string {
	if s.sessionManager == nil {
		return ""
	}

	var cookieStrs []string
	for _, cookie := range s.sessionManager.GetCookies() {
		cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	return strings.Join(cookieStrs, "; ")
}

// Close 关闭扫描器，释放资源
func (s *Scanner) Close() error {
	// 登出会话
	if s.sessionManager != nil && s.sessionManager.IsLoggedIn() {
		ctx := context.Background()
		s.sessionManager.Logout(ctx)
	}

	if s.httpClient != nil {
		s.httpClient.Close()
	}
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
	return nil
}

// 辅助函数

// generateScanID 生成扫描ID
func generateScanID() string {
	return fmt.Sprintf("scan_%d", time.Now().UnixNano())
}

// extractHostFromURL 从URL中提取主机名
func extractHostFromURL(rawURL string) string {
	if u, err := url.Parse(rawURL); err == nil {
		return u.Host
	}
	return "unknown"
}

// crawlAndDiscoverTargets 爬取并发现目标
func (s *Scanner) crawlAndDiscoverTargets(ctx context.Context, initialURLs []string) []string {
	var allTargets []string
	
	for _, initialURL := range initialURLs {
		if s.config.Verbose {
			fmt.Printf("[INFO] 开始爬取: %s\n", initialURL)
		}
		
		// 执行爬取
		crawlResults, err := s.crawler.Crawl(ctx, initialURL)
		if err != nil {
			if s.config.Debug {
				fmt.Printf("[ERROR] 爬取失败: %v\n", err)
			}
			// 如果爬取失败，至少添加原始 URL
			allTargets = append(allTargets, initialURL)
			continue
		}
		
		if s.config.Verbose {
			fmt.Printf("[INFO] 爬取完成，发现 %d 个页面\n", len(crawlResults))
		}
		
		// 分析爬取结果并生成目标 URL
		for _, crawlResult := range crawlResults {
			// 基本 URL
			allTargets = append(allTargets, crawlResult.URL)
			
			// 为每个表单生成测试 URL
			for _, form := range crawlResult.Forms {
				if form.Method == "GET" && form.Action != "" {
					// 为 GET 表单生成测试 URL
					formURL := s.generateFormTestURL(crawlResult.URL, form)
					if formURL != "" {
						allTargets = append(allTargets, formURL)
					}
				}
			}
			
			// 打印功能分析结果
			if s.config.Debug {
				fmt.Printf("[DEBUG] 发现页面: %s\n", crawlResult.URL)
			}
		}
	}
	
	return removeDuplicateURLs(allTargets)
}

// generateFormTestURL 为表单生成测试URL
func (s *Scanner) generateFormTestURL(baseURL string, form *crawler.FormInfo) string {
	if form.Action == "" {
		return ""
	}
	
	// 解析基础URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	// 解析表单action
	actionURL, err := base.Parse(form.Action)
	if err != nil {
		return ""
	}
	
	// 为表单输入添加默认测试参数
	params := url.Values{}
	for _, input := range form.Inputs {
		if input.Name != "" && input.Type != "hidden" && input.Type != "submit" {
			// 为不同类型的输入设置默认测试值
			testValue := s.getTestValueForInput(input)
			params.Add(input.Name, testValue)
		}
	}
	
	// 构建最终URL
	actionURL.RawQuery = params.Encode()
	return actionURL.String()
}

// getTestValueForInput 为输入字段获取测试值
func (s *Scanner) getTestValueForInput(input *crawler.InputInfo) string {
	lowerName := strings.ToLower(input.Name)
	
	// 基于字段名称返回合适的测试值
	switch {
	case strings.Contains(lowerName, "search") || strings.Contains(lowerName, "query") || strings.Contains(lowerName, "keyword"):
		return "test"
	case strings.Contains(lowerName, "id") || strings.Contains(lowerName, "user_id") || strings.Contains(lowerName, "uid"):
		return "1"
	case strings.Contains(lowerName, "name") || strings.Contains(lowerName, "username"):
		return "test"
	case strings.Contains(lowerName, "email"):
		return "test@example.com"
	case strings.Contains(lowerName, "url") || strings.Contains(lowerName, "link"):
		return "http://example.com"
	case strings.Contains(lowerName, "file"):
		return "test.txt"
	default:
		// 基于输入类型返回默认值
		switch input.Type {
		case "number":
			return "1"
		case "email":
			return "test@example.com"
		case "url":
			return "http://example.com"
		default:
			return "test"
		}
	}
}

// shouldCopyHeader 判断是否应该复制响应头到请求头
func shouldCopyHeader(headerName string) bool {
	// 只复制某些特定的头部
	copyHeaders := map[string]bool{
		"Set-Cookie": false, // Cookie单独处理
		"Server":     false, // 不需要复制的头部
		"Date":       false,
	}
	
	if shouldCopy, exists := copyHeaders[headerName]; exists {
		return shouldCopy
	}
	return false // 默认不复制
}

// removeDuplicateURLs 移除重复URL
func removeDuplicateURLs(urls []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, url := range urls {
		if !seen[url] {
			seen[url] = true
			result = append(result, url)
		}
	}
	
	return result
}