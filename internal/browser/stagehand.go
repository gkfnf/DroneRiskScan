package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/dronesec/droneriskscan/internal/transport"
)

// StagehandConfig Stagehand配置
type StagehandConfig struct {
	// Stagehand API配置
	APIEndpoint   string        `json:"api_endpoint"`
	APIKey        string        `json:"api_key,omitempty"`
	Timeout       time.Duration `json:"timeout"`
	
	// 浏览器配置
	Headless      bool   `json:"headless"`
	UserAgent     string `json:"user_agent"`
	Viewport      struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
	
	// AI模型配置
	ModelProvider string `json:"model_provider"` // openai, anthropic, local
	ModelName     string `json:"model_name"`
	Temperature   float32 `json:"temperature"`
	
	// 认证模板
	AuthTemplates map[string]*AuthTemplate `json:"auth_templates"`
}

// AuthTemplate 认证模板
type AuthTemplate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Steps       []*InteractionStep `json:"steps"`
	Success     *SuccessIndicator `json:"success"`
	Fallback    *FallbackStrategy `json:"fallback,omitempty"`
}

// InteractionStep 交互步骤
type InteractionStep struct {
	Type        string            `json:"type"` // navigate, fill, click, wait, extract
	Target      string            `json:"target"` // CSS选择器或描述
	Value       string            `json:"value,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	AIPrompt    string            `json:"ai_prompt,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
}

// SuccessIndicator 成功指示器
type SuccessIndicator struct {
	Type     string `json:"type"` // url, element, text, cookie
	Value    string `json:"value"`
	AIPrompt string `json:"ai_prompt,omitempty"`
}

// FallbackStrategy 降级策略
type FallbackStrategy struct {
	Type   string `json:"type"` // retry, skip, traditional_auth
	Config interface{} `json:"config,omitempty"`
}

// InteractionResult 交互结果
type InteractionResult struct {
	Success      bool                     `json:"success"`
	Message      string                   `json:"message"`
	Cookies      []*http.Cookie           `json:"cookies"`
	LocalStorage map[string]string        `json:"local_storage"`
	SessionData  map[string]interface{}   `json:"session_data"`
	Screenshots  []string                 `json:"screenshots,omitempty"`
	NetworkLogs  []*NetworkLog            `json:"network_logs,omitempty"`
	FunctionPoints []*FunctionPoint       `json:"function_points,omitempty"`
	Duration     time.Duration            `json:"duration"`
	Error        string                   `json:"error,omitempty"`
}

// NetworkLog 网络日志
type NetworkLog struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	RequestBody string           `json:"request_body,omitempty"`
	ResponseBody string          `json:"response_body,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
}

// FunctionPoint 功能点
type FunctionPoint struct {
	URL         string               `json:"url"`
	Type        string               `json:"type"` // form, link, api, websocket
	Method      string               `json:"method"`
	Parameters  map[string]*ParamInfo `json:"parameters"`
	Description string               `json:"description"`
	Selector    string               `json:"selector,omitempty"`
	AIAnalysis  string               `json:"ai_analysis,omitempty"`
}

// ParamInfo 参数信息
type ParamInfo struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // string, number, boolean, file
	Required    bool        `json:"required"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Validation  string      `json:"validation,omitempty"`
	Injectable  bool        `json:"injectable"` // 是否可注入测试
}

// StagehandManager Stagehand管理器
type StagehandManager struct {
	config     *StagehandConfig
	httpClient transport.HTTPClient
	process    *exec.Cmd
	sessionID  string
}

// NewStagehandManager 创建Stagehand管理器
func NewStagehandManager(config *StagehandConfig, httpClient transport.HTTPClient) *StagehandManager {
	if config == nil {
		config = DefaultStagehandConfig()
	}
	
	return &StagehandManager{
		config:     config,
		httpClient: httpClient,
	}
}

// DefaultStagehandConfig 默认配置
func DefaultStagehandConfig() *StagehandConfig {
	config := &StagehandConfig{
		APIEndpoint:   "http://localhost:8080/api/v1",
		Timeout:       60 * time.Second,
		Headless:      true,
		UserAgent:     "DroneRiskScan/1.0 (Stagehand)",
		ModelProvider: "openai",
		ModelName:     "gpt-4-vision-preview",
		Temperature:   0.1,
		AuthTemplates: make(map[string]*AuthTemplate),
	}
	
	config.Viewport.Width = 1920
	config.Viewport.Height = 1080
	
	// 添加默认认证模板
	config.AuthTemplates["bwapp"] = &AuthTemplate{
		Name:        "bWAPP Authentication",
		Description: "bWAPP (Buggy Web Application) login flow",
		Steps: []*InteractionStep{
			{
				Type:     "navigate",
				Target:   "http://127.0.0.1/login.php",
				Timeout:  10 * time.Second,
			},
			{
				Type:     "fill",
				Target:   "input[name='login']",
				Value:    "${username}",
				AIPrompt: "Find the username input field on this login page",
				Timeout:  5 * time.Second,
			},
			{
				Type:     "fill",
				Target:   "input[name='password']",
				Value:    "${password}",
				AIPrompt: "Find the password input field on this login page",
				Timeout:  5 * time.Second,
			},
			{
				Type:     "click",
				Target:   "input[type='submit'], button[type='submit'], .login-button",
				AIPrompt: "Click the login/submit button to authenticate",
				Timeout:  5 * time.Second,
			},
			{
				Type:    "wait",
				Target:  "2000", // wait 2 seconds
				Timeout: 5 * time.Second,
			},
		},
		Success: &SuccessIndicator{
			Type:     "text",
			Value:    "Choose your bug",
			AIPrompt: "Look for successful login indicators like 'Choose your bug' text or dashboard elements",
		},
		Fallback: &FallbackStrategy{
			Type: "traditional_auth",
			Config: map[string]interface{}{
				"login_url": "http://127.0.0.1/login.php",
				"method":    "POST",
			},
		},
	}
	
	// 通用Web应用认证模板
	config.AuthTemplates["generic_web"] = &AuthTemplate{
		Name:        "Generic Web Application",
		Description: "AI-driven authentication for generic web applications",
		Steps: []*InteractionStep{
			{
				Type:     "navigate",
				Target:   "${target_url}",
				AIPrompt: "Navigate to the target URL and analyze the page",
				Timeout:  15 * time.Second,
			},
			{
				Type:     "fill",
				Target:   "detect",
				Value:    "${username}",
				AIPrompt: "Analyze the page and find the username/email input field. Look for common patterns like 'username', 'email', 'login', or similar labels.",
				Timeout:  10 * time.Second,
			},
			{
				Type:     "fill",
				Target:   "detect",
				Value:    "${password}",
				AIPrompt: "Find the password input field on this page. Look for password type inputs or fields labeled 'password', 'pass', etc.",
				Timeout:  10 * time.Second,
			},
			{
				Type:     "click",
				Target:   "detect",
				AIPrompt: "Find and click the login/signin/submit button. Look for buttons with text like 'Login', 'Sign In', 'Submit', or similar.",
				Timeout:  10 * time.Second,
			},
			{
				Type:    "wait",
				Target:  "3000",
				Timeout: 10 * time.Second,
			},
		},
		Success: &SuccessIndicator{
			Type:     "ai_analysis",
			AIPrompt: "Analyze if login was successful. Look for: 1) URL change to dashboard/home, 2) Welcome messages, 3) User profile elements, 4) Logout buttons, 5) Protected content. Return true if authentication appears successful.",
		},
		Fallback: &FallbackStrategy{
			Type: "retry",
			Config: map[string]interface{}{
				"max_attempts": 2,
				"delay":        "5s",
			},
		},
	}
	
	return config
}

// Start 启动Stagehand服务
func (sm *StagehandManager) Start(ctx context.Context) error {
	// TODO: 启动Stagehand服务进程
	// 这里需要根据Stagehand的实际安装和启动方式来实现
	
	fmt.Println("[INFO] Stagehand browser automation service starting...")
	
	// 创建会话
	sessionReq := map[string]interface{}{
		"headless":   sm.config.Headless,
		"user_agent": sm.config.UserAgent,
		"viewport":   sm.config.Viewport,
		"model": map[string]interface{}{
			"provider":    sm.config.ModelProvider,
			"name":        sm.config.ModelName,
			"temperature": sm.config.Temperature,
		},
	}
	
	session, err := sm.createSession(ctx, sessionReq)
	if err != nil {
		return fmt.Errorf("failed to create Stagehand session: %w", err)
	}
	
	sm.sessionID = session["session_id"].(string)
	fmt.Printf("[INFO] Stagehand session created: %s\n", sm.sessionID)
	
	return nil
}

// PerformAuthentication 执行认证
func (sm *StagehandManager) PerformAuthentication(ctx context.Context, templateName string, credentials map[string]string) (*InteractionResult, error) {
	fmt.Printf("[INFO] Performing authentication using Stagehand AI browser automation\n")
	
	startTime := time.Now()
	
	// 获取认证参数
	targetURL := credentials["target_url"]
	username := credentials["username"]
	password := credentials["password"]
	
	if targetURL == "" || username == "" || password == "" {
		return nil, fmt.Errorf("missing required credentials: target_url, username, password")
	}
	
	// 调用Python Stagehand脚本
	pythonPath := sm.getPythonPath()
	scriptPath := "./scripts/stagehand_auth.py"
	
	// 构建命令
	cmd := exec.Command(pythonPath, scriptPath, "auth", targetURL, username, password, templateName)
	
	// 设置超时
	cmdCtx, cancel := context.WithTimeout(ctx, sm.config.Timeout)
	defer cancel()
	cmd = exec.CommandContext(cmdCtx, pythonPath, scriptPath, "auth", targetURL, username, password, templateName)
	
	fmt.Printf("[DEBUG] Executing Stagehand command: %s %s auth %s %s ***\n", pythonPath, scriptPath, targetURL, username)
	
	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		// 尝试获取错误输出
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("stagehand script failed: %s, stderr: %s", err, string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute stagehand script: %w", err)
	}
	
	// 解析Python脚本的JSON输出
	var pythonResult struct {
		Success     bool                     `json:"success"`
		Error       string                   `json:"error"`
		Cookies     []map[string]interface{} `json:"cookies"`
		SessionData map[string]interface{}   `json:"session_data"`
		Screenshots []string                 `json:"screenshots"`
		NetworkLogs []map[string]interface{} `json:"network_logs"`
	}
	
	if err := json.Unmarshal(output, &pythonResult); err != nil {
		return nil, fmt.Errorf("failed to parse stagehand script output: %w, output: %s", err, string(output))
	}
	
	// 构建结果
	result := &InteractionResult{
		Success:      pythonResult.Success,
		Message:      "Stagehand authentication completed",
		SessionData:  pythonResult.SessionData,
		Screenshots:  pythonResult.Screenshots,
		Duration:     time.Since(startTime),
		Error:        pythonResult.Error,
		NetworkLogs:  make([]*NetworkLog, 0),
		Cookies:      make([]*http.Cookie, 0),
		LocalStorage: make(map[string]string),
	}
	
	// 转换Cookies
	for _, cookieData := range pythonResult.Cookies {
		cookie := &http.Cookie{
			Name:     getString(cookieData, "name"),
			Value:    getString(cookieData, "value"),
			Domain:   getString(cookieData, "domain"),
			Path:     getString(cookieData, "path"),
			Secure:   getBool(cookieData, "secure"),
			HttpOnly: getBool(cookieData, "httpOnly"),
		}
		result.Cookies = append(result.Cookies, cookie)
	}
	
	if result.Success {
		result.Message = "Stagehand authentication successful"
		fmt.Printf("[INFO] ✅ Stagehand authentication successful, extracted %d cookies\n", len(result.Cookies))
		
		// 打印Cookie信息用于调试
		for _, cookie := range result.Cookies {
			fmt.Printf("[DEBUG] Cookie: %s=%s (Domain: %s, Path: %s)\n", cookie.Name, cookie.Value, cookie.Domain, cookie.Path)
		}
	} else {
		fmt.Printf("[ERROR] ❌ Stagehand authentication failed: %s\n", result.Error)
	}
	
	return result, nil
}

// DiscoverFunctionPoints 发现功能点
func (sm *StagehandManager) DiscoverFunctionPoints(ctx context.Context, targetURL string, authenticated bool) ([]*FunctionPoint, error) {
	fmt.Printf("[INFO] Discovering function points for: %s (authenticated: %t)\n", targetURL, authenticated)
	
	// 导航到目标页面
	_, err := sm.executeStep(ctx, &InteractionStep{
		Type:    "navigate",
		Target:  targetURL,
		Timeout: 15 * time.Second,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to target: %w", err)
	}
	
	// AI分析页面功能点
	analysisPrompt := `
	Analyze this web page and identify all potential security testing points:
	
	1. Forms (especially with user inputs)
	2. Links with parameters
	3. API endpoints (from JavaScript/AJAX calls)
	4. File upload features
	5. Search functionality
	6. Admin/privileged areas
	7. Database interaction points
	
	For each function point, provide:
	- URL/endpoint
	- HTTP method
	- Parameters (name, type, required)
	- Potential vulnerability types (SQLi, XSS, etc.)
	- Injectable parameters for testing
	
	Return structured JSON data.
	`
	
	functionPoints, err := sm.analyzePageWithAI(ctx, analysisPrompt)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}
	
	fmt.Printf("[INFO] Discovered %d function points\n", len(functionPoints))
	
	return functionPoints, nil
}

// ExecuteInteraction 执行页面交互
func (sm *StagehandManager) ExecuteInteraction(ctx context.Context, interaction *InteractionStep) (*InteractionResult, error) {
	return sm.executeStep(ctx, interaction, nil)
}

// Close 关闭Stagehand管理器
func (sm *StagehandManager) Close() error {
	if sm.sessionID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		sm.closeSession(ctx, sm.sessionID)
	}
	
	if sm.process != nil {
		return sm.process.Process.Kill()
	}
	
	return nil
}

// Private methods

func (sm *StagehandManager) createSession(ctx context.Context, config map[string]interface{}) (map[string]interface{}, error) {
	// 实际上，我们通过Python脚本来管理Stagehand会话
	// 这里只是创建一个会话标识符
	return map[string]interface{}{
		"session_id": fmt.Sprintf("stagehand_%d", time.Now().UnixNano()),
		"status":     "created",
	}, nil
}

func (sm *StagehandManager) executeStep(ctx context.Context, step *InteractionStep, credentials map[string]string) (*InteractionResult, error) {
	// TODO: 实现具体的步骤执行逻辑
	// 这里需要调用Stagehand的API来执行实际的浏览器操作
	
	result := &InteractionResult{
		Success:     true,
		Message:     fmt.Sprintf("Executed %s step", step.Type),
		NetworkLogs: make([]*NetworkLog, 0),
	}
	
	// 模拟执行时间
	time.Sleep(time.Millisecond * 500)
	
	return result, nil
}

func (sm *StagehandManager) verifyAuthSuccess(ctx context.Context, indicator *SuccessIndicator) (bool, error) {
	// TODO: 实现认证成功验证逻辑
	// 根据indicator的类型进行不同的验证
	return true, nil
}

func (sm *StagehandManager) extractCookies(ctx context.Context) ([]*http.Cookie, error) {
	// TODO: 从浏览器会话中提取cookies
	return []*http.Cookie{}, nil
}

func (sm *StagehandManager) extractLocalStorage(ctx context.Context) (map[string]string, error) {
	// TODO: 从浏览器会话中提取localStorage
	return map[string]string{}, nil
}

func (sm *StagehandManager) executeFallback(ctx context.Context, fallback *FallbackStrategy, credentials map[string]string) (*InteractionResult, error) {
	// TODO: 实现降级策略
	return &InteractionResult{Success: false}, nil
}

func (sm *StagehandManager) analyzePageWithAI(ctx context.Context, prompt string) ([]*FunctionPoint, error) {
	// TODO: 使用AI分析页面功能点
	// 这里需要调用Stagehand的AI分析API
	return []*FunctionPoint{}, nil
}

func (sm *StagehandManager) closeSession(ctx context.Context, sessionID string) error {
	// TODO: 关闭Stagehand会话
	return nil
}

// getPythonPath 获取Python解释器路径
func (sm *StagehandManager) getPythonPath() string {
	// 使用虚拟环境的Python
	if _, err := os.Stat("./stagehand_env/bin/python"); err == nil {
		return "./stagehand_env/bin/python"
	}
	
	// 降级到系统Python
	return "python3"
}

// getString 从map中获取字符串值
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getBool 从map中获取布尔值
func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}