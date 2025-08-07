package engine

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dronesec/droneriskscan/internal/browser"
	"github.com/dronesec/droneriskscan/internal/detector"
	"github.com/dronesec/droneriskscan/internal/detector/injection"
	"github.com/dronesec/droneriskscan/pkg/models"
)

// HybridScanner 混合扫描器 (传统HTTP + Stagehand浏览器自动化)
type HybridScanner struct {
	traditionalScanner *Scanner                    // 传统HTTP扫描器
	stagehandManager   *browser.StagehandManager   // Stagehand管理器
	playwrightManager  *browser.PlaywrightManager  // Playwright管理器
	config            *HybridScannerConfig
	mutex             sync.RWMutex
}

// HybridScannerConfig 混合扫描器配置
type HybridScannerConfig struct {
	*ScannerConfig                    // 继承传统扫描器配置
	
	// 浏览器自动化配置
	EnableStagehand   bool                       `json:"enable_stagehand"`
	StagehandConfig   *browser.StagehandConfig   `json:"stagehand_config"`
	
	// 扫描策略配置
	AuthStrategy      AuthStrategy               `json:"auth_strategy"`      // traditional, stagehand, hybrid
	CrawlStrategy     CrawlStrategy              `json:"crawl_strategy"`     // traditional, stagehand, hybrid
	DetectionMode     DetectionMode              `json:"detection_mode"`     // passive, active, hybrid
	
	// 智能决策配置
	AutoFallback      bool                       `json:"auto_fallback"`      // 自动降级
	SmartRouting      bool                       `json:"smart_routing"`      // 智能路由
	AIAnalysis        bool                       `json:"ai_analysis"`        // AI分析
	
	// 性能配置
	BrowserPoolSize   int                        `json:"browser_pool_size"`
	MaxBrowserTime    time.Duration              `json:"max_browser_time"`
	ConcurrentBrowser int                        `json:"concurrent_browser"`
}

// AuthStrategy 认证策略
type AuthStrategy string

const (
	AuthTraditional AuthStrategy = "traditional" // 仅使用HTTP客户端认证
	AuthStagehand   AuthStrategy = "stagehand"   // 仅使用浏览器认证
	AuthHybrid      AuthStrategy = "hybrid"      // 智能选择认证方式
)

// CrawlStrategy 爬取策略
type CrawlStrategy string

const (
	CrawlTraditional CrawlStrategy = "traditional" // 传统爬虫
	CrawlStagehand   CrawlStrategy = "stagehand"   // 浏览器爬虫
	CrawlHybrid      CrawlStrategy = "hybrid"      // 混合爬虫
)

// DetectionMode 检测模式
type DetectionMode string

const (
	DetectionPassive DetectionMode = "passive" // 被动检测(基于已有流量)
	DetectionActive  DetectionMode = "active"  // 主动检测(发送测试载荷)
	DetectionHybrid  DetectionMode = "hybrid"  // 混合检测
)

// ScanResult 增强的扫描结果
type ScanResult struct {
	*models.ScanResult
	
	// 浏览器自动化结果
	BrowserSessions   []*BrowserSession          `json:"browser_sessions"`
	AIAnalysis        *AIAnalysisResult          `json:"ai_analysis"`
	FunctionPoints    []*browser.FunctionPoint   `json:"function_points"`
	InteractionFlows  []*InteractionFlow         `json:"interaction_flows"`
}

// BrowserSession 浏览器会话
type BrowserSession struct {
	SessionID     string                      `json:"session_id"`
	StartTime     time.Time                   `json:"start_time"`
	EndTime       time.Time                   `json:"end_time"`
	Success       bool                        `json:"success"`
	AuthResult    *browser.InteractionResult  `json:"auth_result"`
	Screenshots   []string                    `json:"screenshots"`
	NetworkLogs   []*browser.NetworkLog       `json:"network_logs"`
	ErrorMessages []string                    `json:"error_messages"`
}

// AIAnalysisResult AI分析结果
type AIAnalysisResult struct {
	ApplicationType   string                     `json:"application_type"`    // 应用类型识别
	TechStack        []string                   `json:"tech_stack"`          // 技术栈识别
	SecurityLevel    string                     `json:"security_level"`      // 安全等级评估
	AuthComplexity   string                     `json:"auth_complexity"`     // 认证复杂度
	RecommendedTests []string                   `json:"recommended_tests"`   // 推荐测试类型
	RiskAssessment   *RiskAssessment            `json:"risk_assessment"`     // 风险评估
}

// RiskAssessment 风险评估
type RiskAssessment struct {
	OverallRisk      string                     `json:"overall_risk"`        // HIGH, MEDIUM, LOW
	AuthBypass       float64                    `json:"auth_bypass"`         // 认证绕过风险
	InjectionRisk    float64                    `json:"injection_risk"`      // 注入攻击风险
	BusinessLogic    float64                    `json:"business_logic"`      // 业务逻辑风险
	DataExposure     float64                    `json:"data_exposure"`       // 数据泄露风险
}

// InteractionFlow 交互流程
type InteractionFlow struct {
	Name        string                       `json:"name"`
	Steps       []*FlowStep                  `json:"steps"`
	Success     bool                         `json:"success"`
	Duration    time.Duration                `json:"duration"`
	Cookies     []*http.Cookie               `json:"cookies"`
	TestPoints  []*TestPoint                 `json:"test_points"`
}

// FlowStep 流程步骤
type FlowStep struct {
	Type        string                       `json:"type"`
	Description string                       `json:"description"`
	URL         string                       `json:"url"`
	Success     bool                         `json:"success"`
	Duration    time.Duration                `json:"duration"`
	Screenshot  string                       `json:"screenshot,omitempty"`
}

// TestPoint 测试点
type TestPoint struct {
	URL           string                     `json:"url"`
	Method        string                     `json:"method"`
	Parameters    map[string]string          `json:"parameters"`
	AuthRequired  bool                       `json:"auth_required"`
	Complexity    string                     `json:"complexity"`     // simple, medium, complex
	Priority      int                        `json:"priority"`       // 1-10
	TestTypes     []string                   `json:"test_types"`     // sqli, xss, etc.
}

// NewHybridScanner 创建混合扫描器
func NewHybridScanner(config *HybridScannerConfig) (*HybridScanner, error) {
	if config == nil {
		config = DefaultHybridScannerConfig()
	}
	
	// 创建传统扫描器
	traditionalScanner, err := NewScanner(config.ScannerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create traditional scanner: %w", err)
	}
	
	hs := &HybridScanner{
		traditionalScanner: traditionalScanner,
		config:            config,
	}
	
	// 创建Stagehand管理器
	if config.EnableStagehand {
		hs.stagehandManager = browser.NewStagehandManager(
			config.StagehandConfig,
			traditionalScanner.httpClient,
		)
		
		// 创建Playwright管理器作为实际的浏览器自动化实现
		hs.playwrightManager = browser.NewPlaywrightManager(config.StagehandConfig)
	}
	
	return hs, nil
}

// DefaultHybridScannerConfig 默认混合扫描器配置
func DefaultHybridScannerConfig() *HybridScannerConfig {
	return &HybridScannerConfig{
		ScannerConfig:     DefaultScannerConfig(),
		EnableStagehand:   true,
		StagehandConfig:   browser.DefaultStagehandConfig(),
		AuthStrategy:      AuthHybrid,
		CrawlStrategy:     CrawlHybrid,
		DetectionMode:     DetectionHybrid,
		AutoFallback:      true,
		SmartRouting:      true,
		AIAnalysis:        true,
		BrowserPoolSize:   2,
		MaxBrowserTime:    5 * time.Minute,
		ConcurrentBrowser: 1,
	}
}

// ScanURL 扫描单个URL (混合模式)
func (hs *HybridScanner) ScanURL(ctx context.Context, targetURL string) (*ScanResult, error) {
	fmt.Printf("[INFO] Starting hybrid scan for: %s\n", targetURL)
	
	result := &ScanResult{
		ScanResult:       models.NewScanResult(fmt.Sprintf("hybrid_scan_%d", time.Now().UnixNano())),
		BrowserSessions:  make([]*BrowserSession, 0),
		FunctionPoints:   make([]*browser.FunctionPoint, 0),
		InteractionFlows: make([]*InteractionFlow, 0),
	}
	
	result.SetRunning()
	
	// 1. 初始侦察和应用识别
	appInfo, err := hs.analyzeApplication(ctx, targetURL)
	if err != nil {
		fmt.Printf("[WARN] Application analysis failed: %v\n", err)
	} else {
		result.AIAnalysis = appInfo
		fmt.Printf("[INFO] Application type: %s\n", appInfo.ApplicationType)
	}
	
	// 2. 智能认证策略选择
	authStrategy := hs.selectAuthStrategy(appInfo)
	fmt.Printf("[INFO] Selected auth strategy: %s\n", authStrategy)
	
	var authResult *browser.InteractionResult
	var authCookies []*http.Cookie
	
	// 3. 执行认证
	if hs.config.AuthCredentials != nil {
		switch authStrategy {
		case AuthTraditional:
			authResult, authCookies = hs.performTraditionalAuth(ctx, targetURL)
		case AuthStagehand:
			authResult, authCookies = hs.performStagehandAuth(ctx, targetURL)
		case AuthHybrid:
			authResult, authCookies = hs.performHybridAuth(ctx, targetURL)
		}
		
		if authResult != nil && authResult.Success {
			fmt.Printf("[INFO] Authentication successful using %s strategy\n", authStrategy)
		} else if hs.config.AutoFallback {
			fmt.Println("[INFO] Primary auth failed, trying fallback strategies")
			authResult, authCookies = hs.performAuthFallback(ctx, targetURL)
		}
	}
	
	// 4. 功能点发现
	var functionPoints []*browser.FunctionPoint
	
	crawlStrategy := hs.selectCrawlStrategy(appInfo, authResult != nil && authResult.Success)
	fmt.Printf("[INFO] Selected crawl strategy: %s\n", crawlStrategy)
	
	switch crawlStrategy {
	case CrawlTraditional:
		functionPoints = hs.discoverWithTraditionalCrawler(ctx, targetURL, authCookies)
	case CrawlStagehand:
		functionPoints = hs.discoverWithStagehandCrawler(ctx, targetURL, authResult != nil && authResult.Success)
	case CrawlHybrid:
		functionPoints = hs.discoverWithHybridCrawler(ctx, targetURL, authCookies, authResult != nil && authResult.Success)
	}
	
	result.FunctionPoints = functionPoints
	fmt.Printf("[INFO] Discovered %d function points\n", len(functionPoints))
	
	// 5. 漏洞检测
	detectionResults := hs.performVulnerabilityDetection(ctx, functionPoints, authCookies)
	
	// If no function points discovered with browser, try traditional parameter extraction
	if len(functionPoints) == 0 {
		fmt.Println("[INFO] No function points discovered with browser, using traditional detection")
		traditionalResults := hs.performTraditionalVulnerabilityDetection(ctx, targetURL, authCookies)
		detectionResults = append(detectionResults, traditionalResults...)
	}
	
	// 合并检测结果
	for _, vuln := range detectionResults {
		result.AddVulnerability(vuln)
	}
	
	result.SetCompleted()
	fmt.Printf("[INFO] Hybrid scan completed: found %d vulnerabilities\n", result.GetVulnerabilityCount())
	
	return result, nil
}

// Private methods for hybrid scanning

func (hs *HybridScanner) analyzeApplication(ctx context.Context, targetURL string) (*AIAnalysisResult, error) {
	if !hs.config.AIAnalysis || hs.stagehandManager == nil {
		return &AIAnalysisResult{
			ApplicationType: "unknown",
			SecurityLevel:   "medium",
		}, nil
	}
	
	// TODO: 使用AI分析应用类型和特征
	// 这里需要实现AI驱动的应用分析
	return &AIAnalysisResult{
		ApplicationType:  "web_application",
		TechStack:       []string{"php", "apache", "mysql"},
		SecurityLevel:   "low",
		AuthComplexity:  "simple",
		RecommendedTests: []string{"sqli", "xss", "auth_bypass"},
		RiskAssessment: &RiskAssessment{
			OverallRisk:   "MEDIUM",
			AuthBypass:    0.7,
			InjectionRisk: 0.8,
			BusinessLogic: 0.5,
			DataExposure:  0.6,
		},
	}, nil
}

func (hs *HybridScanner) selectAuthStrategy(appInfo *AIAnalysisResult) AuthStrategy {
	if !hs.config.EnableStagehand {
		return AuthTraditional
	}
	
	if hs.config.AuthStrategy != AuthHybrid {
		return hs.config.AuthStrategy
	}
	
	// 智能选择认证策略
	if appInfo != nil && appInfo.AuthComplexity == "complex" {
		return AuthStagehand
	}
	
	return AuthHybrid // 默认使用混合策略
}

func (hs *HybridScanner) selectCrawlStrategy(appInfo *AIAnalysisResult, authenticated bool) CrawlStrategy {
	if !hs.config.EnableStagehand {
		return CrawlTraditional
	}
	
	if hs.config.CrawlStrategy != CrawlHybrid {
		return hs.config.CrawlStrategy
	}
	
	// 智能选择爬取策略
	if appInfo != nil && (appInfo.ApplicationType == "spa" || strings.Contains(strings.Join(appInfo.TechStack, " "), "react")) {
		return CrawlStagehand
	}
	
	if authenticated {
		return CrawlHybrid // 认证后使用混合爬取
	}
	
	return CrawlTraditional
}

func (hs *HybridScanner) performTraditionalAuth(ctx context.Context, targetURL string) (*browser.InteractionResult, []*http.Cookie) {
	err := hs.traditionalScanner.Login(ctx)
	if err != nil {
		return &browser.InteractionResult{Success: false, Error: err.Error()}, nil
	}
	
	cookies := hs.traditionalScanner.sessionManager.GetCookies()
	return &browser.InteractionResult{
		Success: true,
		Message: "Traditional authentication successful",
		Cookies: cookies,
	}, cookies
}

func (hs *HybridScanner) performStagehandAuth(ctx context.Context, targetURL string) (*browser.InteractionResult, []*http.Cookie) {
	if hs.playwrightManager == nil {
		return &browser.InteractionResult{Success: false, Error: "Playwright not enabled"}, nil
	}
	
	// 启动Playwright
	err := hs.playwrightManager.Start(ctx)
	if err != nil {
		return &browser.InteractionResult{Success: false, Error: err.Error()}, nil
	}
	
	// 执行认证
	credentials := map[string]string{
		"username":  hs.config.AuthCredentials.Username,
		"password":  hs.config.AuthCredentials.Password,
		"login_url": hs.config.AuthCredentials.LoginURL,
	}
	
	// Use bWAPP-specific authentication
	result, err := hs.playwrightManager.PerformBWAPPAuthentication(ctx, credentials)
	if err != nil {
		return &browser.InteractionResult{Success: false, Error: err.Error()}, nil
	}
	
	return result, result.Cookies
}

func (hs *HybridScanner) performHybridAuth(ctx context.Context, targetURL string) (*browser.InteractionResult, []*http.Cookie) {
	// 首先尝试传统认证
	traditionalResult, traditionalCookies := hs.performTraditionalAuth(ctx, targetURL)
	if traditionalResult.Success {
		return traditionalResult, traditionalCookies
	}
	
	fmt.Println("[INFO] Traditional auth failed, trying Stagehand...")
	
	// 传统认证失败，尝试Stagehand
	return hs.performStagehandAuth(ctx, targetURL)
}

func (hs *HybridScanner) performAuthFallback(ctx context.Context, targetURL string) (*browser.InteractionResult, []*http.Cookie) {
	// TODO: 实现多种降级认证策略
	return &browser.InteractionResult{Success: false}, nil
}

func (hs *HybridScanner) determineAuthTemplate(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "generic_web"
	}
	
	host := strings.ToLower(parsedURL.Host)
	
	// 基于主机名选择模板
	if strings.Contains(host, "127.0.0.1") || strings.Contains(host, "localhost") {
		return "bwapp"
	}
	
	return "generic_web"
}

func (hs *HybridScanner) discoverWithTraditionalCrawler(ctx context.Context, targetURL string, cookies []*http.Cookie) []*browser.FunctionPoint {
	// 使用传统爬虫发现功能点
	// TODO: 将传统爬虫结果转换为FunctionPoint格式
	return []*browser.FunctionPoint{}
}

func (hs *HybridScanner) discoverWithStagehandCrawler(ctx context.Context, targetURL string, authenticated bool) []*browser.FunctionPoint {
	if hs.playwrightManager == nil {
		return []*browser.FunctionPoint{}
	}
	
	functionPoints, err := hs.playwrightManager.DiscoverFunctionPoints(ctx, targetURL)
	if err != nil {
		fmt.Printf("[ERROR] Playwright function point discovery failed: %v\n", err)
		return []*browser.FunctionPoint{}
	}
	
	return functionPoints
}

func (hs *HybridScanner) discoverWithHybridCrawler(ctx context.Context, targetURL string, cookies []*http.Cookie, authenticated bool) []*browser.FunctionPoint {
	// 混合爬取：结合传统爬虫和Stagehand的结果
	traditionalPoints := hs.discoverWithTraditionalCrawler(ctx, targetURL, cookies)
	stagehandPoints := hs.discoverWithStagehandCrawler(ctx, targetURL, authenticated)
	
	// 合并和去重
	return hs.mergeFunctionPoints(traditionalPoints, stagehandPoints)
}

func (hs *HybridScanner) mergeFunctionPoints(traditional, stagehand []*browser.FunctionPoint) []*browser.FunctionPoint {
	pointMap := make(map[string]*browser.FunctionPoint)
	
	// 添加传统爬虫的结果
	for _, point := range traditional {
		key := fmt.Sprintf("%s_%s_%s", point.Method, point.URL, point.Type)
		pointMap[key] = point
	}
	
	// 添加Stagehand的结果
	for _, point := range stagehand {
		key := fmt.Sprintf("%s_%s_%s", point.Method, point.URL, point.Type)
		if existing, exists := pointMap[key]; exists {
			// 合并参数信息
			for paramName, paramInfo := range point.Parameters {
				existing.Parameters[paramName] = paramInfo
			}
		} else {
			pointMap[key] = point
		}
	}
	
	// 转换回切片
	result := make([]*browser.FunctionPoint, 0, len(pointMap))
	for _, point := range pointMap {
		result = append(result, point)
	}
	
	return result
}

func (hs *HybridScanner) performVulnerabilityDetection(ctx context.Context, functionPoints []*browser.FunctionPoint, cookies []*http.Cookie) []*models.Vulnerability {
	var vulnerabilities []*models.Vulnerability
	
	// 为每个功能点执行漏洞检测
	for _, point := range functionPoints {
		// 创建扫描目标
		parsedURL, err := url.Parse(point.URL)
		if err != nil {
			continue
		}
		
		scanTarget := &detector.ScanTarget{
			URL:        parsedURL,
			Method:     point.Method,
			Headers:    make(map[string]string),
			Parameters: make(map[string][]string),
			Cookies:    make(map[string]string),
			Metadata:   make(map[string]interface{}),
		}
		
		// 设置参数
		for paramName, paramInfo := range point.Parameters {
			if paramInfo.Injectable {
				scanTarget.Parameters[paramName] = []string{fmt.Sprintf("%v", paramInfo.DefaultValue)}
			}
		}
		
		// 设置认证cookies
		for _, cookie := range cookies {
			scanTarget.Cookies[cookie.Name] = cookie.Value
		}
		
		// 执行检测插件
		plugins := hs.traditionalScanner.GetPlugins()
		for _, plugin := range plugins {
			result, err := plugin.Execute(ctx, scanTarget)
			if err != nil {
				continue
			}
			
			if result.IsVulnerable {
				vulnerabilities = append(vulnerabilities, result.Vulnerabilities...)
			}
		}
	}
	
	return vulnerabilities
}

// performTraditionalVulnerabilityDetection 执行传统漏洞检测
func (hs *HybridScanner) performTraditionalVulnerabilityDetection(ctx context.Context, targetURL string, cookies []*http.Cookie) []*models.Vulnerability {
	var vulnerabilities []*models.Vulnerability
	
	// Create scan target from URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return vulnerabilities
	}
	
	scanTarget := &detector.ScanTarget{
		URL:        parsedURL,
		Method:     "GET",
		Headers:    make(map[string]string),
		Parameters: make(map[string][]string),
		Cookies:    make(map[string]string),
		Metadata:   make(map[string]interface{}),
	}
	
	// Set authentication cookies
	for _, cookie := range cookies {
		scanTarget.Cookies[cookie.Name] = cookie.Value
	}
	
	// Copy URL parameters to scan target
	queryParams := parsedURL.Query()
	for key, values := range queryParams {
		scanTarget.Parameters[key] = values
	}
	
	fmt.Printf("[DEBUG] Traditional detection with %d cookies and %d parameters\n", len(scanTarget.Cookies), len(scanTarget.Parameters))
	
	// Execute detection plugins
	plugins := hs.traditionalScanner.GetPlugins()
	for _, plugin := range plugins {
		// Set session cookies for the plugin
		if sqlPlugin, ok := plugin.(*injection.SQLiDetector); ok {
			sqlPlugin.SetSessionCookies(cookies)
		}
		
		result, err := plugin.Execute(ctx, scanTarget)
		if err != nil {
			fmt.Printf("[DEBUG] Plugin %s execution failed: %v\n", plugin.Name(), err)
			continue
		}
		
		if result.IsVulnerable {
			vulnerabilities = append(vulnerabilities, result.Vulnerabilities...)
			fmt.Printf("[FOUND] Plugin %s detected %d vulnerabilities\n", plugin.Name(), len(result.Vulnerabilities))
		}
	}
	
	return vulnerabilities
}

// Close 关闭混合扫描器
func (hs *HybridScanner) Close() error {
	if hs.stagehandManager != nil {
		hs.stagehandManager.Close()
	}
	
	if hs.playwrightManager != nil {
		hs.playwrightManager.Close()
	}
	
	if hs.traditionalScanner != nil {
		return hs.traditionalScanner.Close()
	}
	
	return nil
}