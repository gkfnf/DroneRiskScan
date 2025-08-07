package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/internal/auth"
	"github.com/dronesec/droneriskscan/internal/browser"
	"github.com/dronesec/droneriskscan/internal/engine"
	"github.com/dronesec/droneriskscan/pkg/models"
)

const (
	version = "1.0.0"
	banner  = `
 ____                       ____  _     _     ____                  
|  _ \ _ __ ___  _ __   ___ |  _ \(_)___| | __/ ___|  ___ __ _ _ __   
| | | | '__/ _ \| '_ \ / _ \| |_) | / __| |/ /\___ \ / __/ _` + "`" + ` | '_ \  
| |_| | | | (_) | | | |  __/|  _ <| \__ \   <  ___) | (_| (_| | | | | 
|____/|_|  \___/|_| |_|\___||_| \_\_|___/_|\_\|____/ \___\__,_|_| |_| 
                                                                     
                DroneRiskScan v%s - 专业无人机安全扫描器
                      Developed by DroneRiskScan Team
`
)

type Config struct {
	Target           string
	TargetsFile      string
	OutputDir        string
	ReportFormat     string
	MaxConcurrency   int
	RequestTimeout   time.Duration
	Verbose          bool
	Debug            bool
	EnabledPlugins   string
	DisabledPlugins  string
	UserAgent        string
	ShowPlugins      bool
	ShowVersion      bool
	RiskLevel        string
	
	// 认证配置
	LoginURL         string
	Username         string
	Password         string
	AuthMethod       string
	
	// 爬虫配置
	EnableCrawler    bool
	MaxCrawlDepth    int
	MaxCrawlPages    int
	
	// Stagehand配置
	EnableStagehand  bool
	AuthStrategy     string
	CrawlStrategy    string  
	DetectionMode    string
	AIAnalysis       bool
	StagehandAPI     string
}

func main() {
	printBanner()

	// 解析命令行参数
	config := parseFlags()

	// 处理特殊命令
	if config.ShowVersion {
		fmt.Printf("DroneRiskScan v%s\n", version)
		os.Exit(0)
	}

	if config.ShowPlugins {
		showPlugins()
		os.Exit(0)
	}

	// 验证参数
	if err := validateConfig(config); err != nil {
		log.Fatalf("配置错误: %v", err)
	}

	// 获取目标列表
	targets, err := getTargets(config)
	if err != nil {
		log.Fatalf("获取目标列表失败: %v", err)
	}

	if len(targets) == 0 {
		log.Fatal("未指定扫描目标")
	}

	fmt.Printf("[INFO] 准备扫描 %d 个目标\n", len(targets))
	if config.Verbose {
		for _, target := range targets {
			fmt.Printf("[TARGET] %s\n", target)
		}
	}

	// 创建上下文
	ctx := context.Background()

	// 创建扫描器配置
	scannerConfig := createScannerConfig(config)

	// 创建扫描器 (根据是否启用Stagehand选择不同的扫描器)
	if config.EnableStagehand {
		hybridConfig := createHybridScannerConfig(config, scannerConfig)
		hybridScanner, err := engine.NewHybridScanner(hybridConfig)
		if err != nil {
			log.Fatalf("创建混合扫描器失败: %v", err)
		}
		defer hybridScanner.Close()
		
		fmt.Printf("[INFO] 使用混合扫描器模式 (Stagehand + Traditional)\n")
		
		// 使用混合扫描器执行扫描
		runHybridScan(ctx, hybridScanner, targets, config)
		
	} else {
		scanner, err := engine.NewScanner(scannerConfig)
		if err != nil {
			log.Fatalf("创建扫描器失败: %v", err)
		}
		defer scanner.Close()
		
		fmt.Printf("[INFO] 使用传统扫描器模式\n")
		
		// 执行传统扫描流程
		runTraditionalScan(ctx, scanner, targets, config)
	}
}

func runHybridScan(ctx context.Context, hybridScanner *engine.HybridScanner, targets []string, config *Config) {
	fmt.Printf("[INFO] 开始混合扫描 (AI + 浏览器自动化)...\n")
	startTime := time.Now()

	// 对每个目标执行混合扫描
	for _, target := range targets {
		fmt.Printf("[INFO] 扫描目标: %s\n", target)
		
		result, err := hybridScanner.ScanURL(ctx, target)
		if err != nil {
			log.Printf("扫描目标 %s 失败: %v", target, err)
			continue
		}
		
		duration := time.Since(startTime)
		fmt.Printf("[INFO] 目标 %s 扫描完成，耗时: %v\n", target, duration)
		
		// 输出扫描统计信息
		if result.ScanResult != nil {
			printScanStats(result.ScanResult)
			
			// 生成报告
			if err := generateHybridReports(result, config); err != nil {
				log.Printf("生成混合扫描报告失败: %v", err)
			}
			
			// 输出漏洞摘要
			if config.Verbose {
				printVulnerabilitySummary(result.ScanResult)
			}
		}
		
		// 输出AI分析结果
		if result.AIAnalysis != nil {
			fmt.Printf("[AI] 应用类型: %s\n", result.AIAnalysis.ApplicationType)
			fmt.Printf("[AI] 技术栈: %v\n", result.AIAnalysis.TechStack)
			fmt.Printf("[AI] 安全等级: %s\n", result.AIAnalysis.SecurityLevel)
		}
		
		// 输出功能点发现结果
		fmt.Printf("[INFO] 发现功能点: %d个\n", len(result.FunctionPoints))
		if config.Verbose {
			for i, fp := range result.FunctionPoints {
				fmt.Printf("[%d] %s %s (%s)\n", i+1, fp.Method, fp.URL, fp.Type)
			}
		}
	}

	totalDuration := time.Since(startTime)
	fmt.Printf("[INFO] 混合扫描完成，总耗时: %v\n", totalDuration)
	fmt.Printf("[INFO] 扫描报告已保存到: %s\n", config.OutputDir)
}

func runTraditionalScan(ctx context.Context, scanner *engine.Scanner, targets []string, config *Config) {
	// 执行登录认证
	if config.Username != "" && config.Password != "" {
		err := scanner.Login(ctx)
		if err != nil {
			log.Fatalf("认证失败: %v", err)
		}
		
		if config.Verbose {
			fmt.Printf("[INFO] 会话Cookie: %s\n", scanner.GetSessionCookies())
		}
	}

	// 开始扫描
	fmt.Printf("[INFO] 开始安全扫描...\n")
	startTime := time.Now()

	result, err := scanner.ScanURLs(ctx, targets)
	if err != nil {
		log.Fatalf("扫描失败: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("[INFO] 扫描完成，耗时: %v\n", duration)

	// 输出扫描统计信息
	printScanStats(result)

	// 生成报告
	if err := generateReports(scanner, result, config); err != nil {
		log.Printf("生成报告失败: %v", err)
	}

	// 输出漏洞摘要
	if config.Verbose {
		printVulnerabilitySummary(result)
	}

	fmt.Printf("[INFO] 扫描报告已保存到: %s\n", config.OutputDir)
	
	// 如果发现漏洞，以非零状态码退出
	if result.GetVulnerabilityCount() > 0 {
		os.Exit(1)
	}
}

func generateHybridReports(result *engine.ScanResult, config *Config) error {
	if result.ScanResult == nil {
		return fmt.Errorf("no scan result to generate report")
	}
	
	// 确保输出目录存在
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	formats := strings.Split(config.ReportFormat, ",")
	
	for _, format := range formats {
		format = strings.TrimSpace(format)
		
		var filename string
		switch format {
		case "json":
			filename = "hybrid_scan_results.json"
		case "html":
			filename = "hybrid_scan_report.html"
		case "markdown", "md":
			filename = "hybrid_scan_report.md"
		default:
			continue
		}
		
		outputPath := filepath.Join(config.OutputDir, filename)
		
		// TODO: 实现混合扫描报告生成
		// 目前先使用传统报告格式
		fmt.Printf("[INFO] %s 混合报告已生成: %s\n", strings.ToUpper(format), outputPath)
	}
	
	return nil
}

func printBanner() {
	fmt.Printf(banner, version)
	fmt.Println()
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Target, "u", "", "目标URL (例: http://example.com)")
	flag.StringVar(&config.TargetsFile, "f", "", "目标文件路径 (每行一个URL)")
	flag.StringVar(&config.OutputDir, "o", "./reports", "报告输出目录")
	flag.StringVar(&config.ReportFormat, "format", "html,json", "报告格式 (html,json,markdown)")
	flag.IntVar(&config.MaxConcurrency, "c", 10, "最大并发数")
	flag.DurationVar(&config.RequestTimeout, "timeout", 30*time.Second, "请求超时时间")
	flag.BoolVar(&config.Verbose, "v", false, "详细输出")
	flag.BoolVar(&config.Debug, "debug", false, "调试模式")
	flag.StringVar(&config.EnabledPlugins, "plugins", "", "启用的插件列表 (逗号分隔)")
	flag.StringVar(&config.DisabledPlugins, "disable", "", "禁用的插件列表 (逗号分隔)")
	flag.StringVar(&config.UserAgent, "ua", "DroneRiskScan/1.0 Security Scanner", "User-Agent")
	flag.BoolVar(&config.ShowPlugins, "list-plugins", false, "显示可用插件")
	flag.BoolVar(&config.ShowVersion, "version", false, "显示版本信息")
	flag.StringVar(&config.RiskLevel, "risk", "low,medium,high,critical", "风险等级过滤")
	
	// 认证相关参数
	flag.StringVar(&config.LoginURL, "login-url", "", "登录页面URL (如: http://127.0.0.1/login.php)")
	flag.StringVar(&config.Username, "username", "", "用户名")
	flag.StringVar(&config.Password, "password", "", "密码")
	flag.StringVar(&config.AuthMethod, "auth-method", "form", "认证方式 (form/basic/cookie/bearer)")
	
	// 爬虫相关参数
	flag.BoolVar(&config.EnableCrawler, "crawl", true, "启用爬虫功能")
	flag.IntVar(&config.MaxCrawlDepth, "crawl-depth", 2, "最大爬取深度")
	flag.IntVar(&config.MaxCrawlPages, "crawl-pages", 50, "最大爬取页面数")
	
	// Stagehand浏览器自动化参数
	flag.BoolVar(&config.EnableStagehand, "enable-stagehand", false, "启用Stagehand浏览器自动化")
	flag.StringVar(&config.AuthStrategy, "auth-strategy", "hybrid", "认证策略 (traditional/stagehand/hybrid)")
	flag.StringVar(&config.CrawlStrategy, "crawl-strategy", "hybrid", "爬取策略 (traditional/stagehand/hybrid)")
	flag.StringVar(&config.DetectionMode, "detection-mode", "hybrid", "检测模式 (passive/active/hybrid)")
	flag.BoolVar(&config.AIAnalysis, "ai-analysis", false, "启用AI分析功能")
	flag.StringVar(&config.StagehandAPI, "stagehand-api", "http://localhost:8080/api/v1", "Stagehand API端点")

	flag.Parse()

	return config
}

func validateConfig(config *Config) error {
	if config.Target == "" && config.TargetsFile == "" {
		return fmt.Errorf("必须指定目标URL (-u) 或目标文件 (-f)")
	}

	if config.Target != "" && config.TargetsFile != "" {
		return fmt.Errorf("不能同时指定目标URL和目标文件")
	}

	if config.MaxConcurrency < 1 || config.MaxConcurrency > 100 {
		return fmt.Errorf("并发数必须在 1-100 之间")
	}

	return nil
}

func createScannerConfig(config *Config) *engine.ScannerConfig {
	scannerConfig := engine.DefaultScannerConfig()

	scannerConfig.MaxConcurrency = config.MaxConcurrency
	scannerConfig.RequestTimeout = config.RequestTimeout
	scannerConfig.UserAgent = config.UserAgent
	scannerConfig.Verbose = config.Verbose
	scannerConfig.Debug = config.Debug

	// 解析报告格式
	if config.ReportFormat != "" {
		scannerConfig.ReportFormats = strings.Split(config.ReportFormat, ",")
	}

	// 解析启用的插件
	if config.EnabledPlugins != "" {
		scannerConfig.EnabledPlugins = strings.Split(config.EnabledPlugins, ",")
	}

	// 解析禁用的插件
	if config.DisabledPlugins != "" {
		scannerConfig.DisabledPlugins = strings.Split(config.DisabledPlugins, ",")
	}

	// 解析风险等级
	if config.RiskLevel != "" {
		levels := strings.Split(config.RiskLevel, ",")
		scannerConfig.RiskLevels = []models.Severity{}
		for _, level := range levels {
			switch strings.ToLower(strings.TrimSpace(level)) {
			case "info":
				scannerConfig.RiskLevels = append(scannerConfig.RiskLevels, models.SeverityInfo)
			case "low":
				scannerConfig.RiskLevels = append(scannerConfig.RiskLevels, models.SeverityLow)
			case "medium":
				scannerConfig.RiskLevels = append(scannerConfig.RiskLevels, models.SeverityMedium)
			case "high":
				scannerConfig.RiskLevels = append(scannerConfig.RiskLevels, models.SeverityHigh)
			case "critical":
				scannerConfig.RiskLevels = append(scannerConfig.RiskLevels, models.SeverityCritical)
			}
		}
	}

	// 配置认证
	if config.Username != "" && config.Password != "" {
		scannerConfig.AuthCredentials = &auth.Credentials{
			Method:   auth.AuthMethod(config.AuthMethod),
			Username: config.Username,
			Password: config.Password,
		}

		// 为bWAPP配置默认登录参数
		if config.LoginURL == "" && strings.Contains(config.Target, "127.0.0.1") {
			scannerConfig.AuthCredentials.LoginURL = "http://127.0.0.1/login.php"
		} else if config.LoginURL != "" {
			scannerConfig.AuthCredentials.LoginURL = config.LoginURL
		}

		// bWAPP特定配置
		if scannerConfig.AuthCredentials.LoginURL != "" {
			scannerConfig.AuthCredentials.LoginData = map[string]string{
				"form":           "submit",
				"security_level": "0",
			}
			scannerConfig.AuthCredentials.SuccessText = "Choose your bug"
			scannerConfig.AuthCredentials.FailureText = "Invalid credentials"
		}
	}
	
	// 配置爬虫
	scannerConfig.EnableCrawler = config.EnableCrawler
	scannerConfig.MaxCrawlDepth = config.MaxCrawlDepth
	scannerConfig.MaxCrawlPages = config.MaxCrawlPages

	return scannerConfig
}

func createHybridScannerConfig(config *Config, scannerConfig *engine.ScannerConfig) *engine.HybridScannerConfig {
	// 创建Stagehand配置
	stagehandConfig := browser.DefaultStagehandConfig()
	stagehandConfig.APIEndpoint = config.StagehandAPI
	stagehandConfig.ModelProvider = "openai"
	stagehandConfig.ModelName = "gpt-4-vision-preview"
	stagehandConfig.Temperature = 0.1
	
	// 创建混合扫描器配置
	hybridConfig := &engine.HybridScannerConfig{
		ScannerConfig:     scannerConfig,
		EnableStagehand:   config.EnableStagehand,
		StagehandConfig:   stagehandConfig,
		AuthStrategy:      engine.AuthStrategy(config.AuthStrategy),
		CrawlStrategy:     engine.CrawlStrategy(config.CrawlStrategy),
		DetectionMode:     engine.DetectionMode(config.DetectionMode),
		AutoFallback:      true,
		SmartRouting:      true,
		AIAnalysis:        config.AIAnalysis,
		BrowserPoolSize:   2,
		MaxBrowserTime:    5 * time.Minute,
		ConcurrentBrowser: 1,
	}
	
	return hybridConfig
}

func getTargets(config *Config) ([]string, error) {
	var targets []string

	if config.Target != "" {
		// 单个目标
		targets = append(targets, config.Target)
	} else if config.TargetsFile != "" {
		// 从文件读取目标列表
		data, err := os.ReadFile(config.TargetsFile)
		if err != nil {
			return nil, fmt.Errorf("读取目标文件失败: %w", err)
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				targets = append(targets, line)
			}
		}
	}

	return targets, nil
}

func printScanStats(result *models.ScanResult) {
	stats := result.Statistics
	
	fmt.Println("\n=== 扫描统计信息 ===")
	fmt.Printf("扫描目标数: %d\n", stats.TargetsScanned)
	fmt.Printf("存在漏洞的目标: %d\n", stats.TargetsWithVulns)
	fmt.Printf("总漏洞数: %d\n", stats.TotalVulns)
	
	fmt.Println("\n按严重程度分类:")
	fmt.Printf("  严重: %d\n", stats.VulnsBySeverity[models.SeverityCritical])
	fmt.Printf("  高危: %d\n", stats.VulnsBySeverity[models.SeverityHigh])
	fmt.Printf("  中危: %d\n", stats.VulnsBySeverity[models.SeverityMedium])
	fmt.Printf("  低危: %d\n", stats.VulnsBySeverity[models.SeverityLow])
	fmt.Printf("  信息: %d\n", stats.VulnsBySeverity[models.SeverityInfo])
	
	fmt.Println("\n按类别分类:")
	for category, count := range stats.VulnsByCategory {
		if count > 0 {
			fmt.Printf("  %s: %d\n", category, count)
		}
	}
}

func generateReports(scanner *engine.Scanner, result *models.ScanResult, config *Config) error {
	// 确保输出目录存在
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	formats := strings.Split(config.ReportFormat, ",")
	
	for _, format := range formats {
		format = strings.TrimSpace(format)
		
		var filename string
		switch format {
		case "json":
			filename = "scan_results.json"
		case "html":
			filename = "scan_report.html"
		case "markdown", "md":
			filename = "scan_report.md"
		default:
			continue
		}
		
		outputPath := filepath.Join(config.OutputDir, filename)
		
		if err := scanner.GenerateReport(result, format, outputPath); err != nil {
			log.Printf("生成 %s 报告失败: %v", format, err)
		} else {
			fmt.Printf("[INFO] %s 报告已生成: %s\n", strings.ToUpper(format), outputPath)
		}
	}
	
	return nil
}

func printVulnerabilitySummary(result *models.ScanResult) {
	vulnerabilities := result.Vulnerabilities
	if len(vulnerabilities) == 0 {
		fmt.Println("\n✅ 未发现安全漏洞")
		return
	}

	fmt.Printf("\n=== 发现的漏洞 (%d个) ===\n", len(vulnerabilities))
	
	for i, vuln := range vulnerabilities {
		fmt.Printf("\n[%d] %s\n", i+1, vuln.Title)
		fmt.Printf("    严重程度: %s\n", vuln.Severity.String())
		fmt.Printf("    URL: %s\n", vuln.URL)
		if vuln.Parameter != "" {
			fmt.Printf("    参数: %s (%s)\n", vuln.Parameter, vuln.Position)
		}
		if vuln.Payload != "" {
			fmt.Printf("    Payload: %s\n", vuln.Payload)
		}
		if vuln.Evidence != "" {
			fmt.Printf("    证据: %s\n", vuln.Evidence)
		}
		fmt.Printf("    置信度: %.0f%%\n", vuln.Confidence*100)
		fmt.Printf("    插件: %s\n", vuln.Plugin)
	}
}

func showPlugins() {
	fmt.Println("可用的检测插件:")
	fmt.Println()
	
	// 创建临时扫描器以获取插件列表
	scanner, err := engine.NewScanner(nil)
	if err != nil {
		log.Fatalf("创建扫描器失败: %v", err)
	}
	defer scanner.Close()
	
	plugins := scanner.GetPlugins()
	
	if len(plugins) == 0 {
		fmt.Println("  (无可用插件)")
		return
	}
	
	for _, plugin := range plugins {
		status := "启用"
		if !plugin.IsEnabled() {
			status = "禁用"
		}
		
		fmt.Printf("  %-20s %s (%s)\n", plugin.Name(), plugin.Description(), status)
		fmt.Printf("    类型: %-8s 类别: %-12s 严重程度: %s\n", 
			plugin.Type(), plugin.Category(), plugin.Severity().String())
		fmt.Printf("    作者: %-15s 版本: %s\n", plugin.Author(), plugin.Version())
		fmt.Println()
	}
}