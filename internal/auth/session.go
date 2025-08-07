package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/internal/transport"
)

// AuthMethod 认证方式
type AuthMethod string

const (
	AuthMethodBasic  AuthMethod = "basic"
	AuthMethodForm   AuthMethod = "form"
	AuthMethodCookie AuthMethod = "cookie"
	AuthMethodBearer AuthMethod = "bearer"
)

// Credentials 认证凭据
type Credentials struct {
	Method   AuthMethod `json:"method"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Token    string     `json:"token,omitempty"`
	
	// 表单认证相关
	LoginURL    string            `json:"login_url,omitempty"`
	LoginData   map[string]string `json:"login_data,omitempty"`
	SuccessText string            `json:"success_text,omitempty"`
	FailureText string            `json:"failure_text,omitempty"`
	
	// Cookie会话
	Cookies map[string]string `json:"cookies,omitempty"`
}

// SessionManager 会话管理器
type SessionManager struct {
	httpClient  transport.HTTPClient
	credentials *Credentials
	cookies     []*http.Cookie
	sessionID   string
	isLoggedIn  bool
	loginTime   time.Time
}

// NewSessionManager 创建会话管理器
func NewSessionManager(httpClient transport.HTTPClient, credentials *Credentials) *SessionManager {
	return &SessionManager{
		httpClient:  httpClient,
		credentials: credentials,
		cookies:     make([]*http.Cookie, 0),
	}
}

// Login 执行登录
func (sm *SessionManager) Login(ctx context.Context) error {
	if sm.credentials == nil {
		return fmt.Errorf("未配置认证凭据")
	}

	switch sm.credentials.Method {
	case AuthMethodForm:
		return sm.loginWithForm(ctx)
	case AuthMethodBasic:
		return sm.loginWithBasic(ctx)
	case AuthMethodCookie:
		return sm.loginWithCookie(ctx)
	case AuthMethodBearer:
		return sm.loginWithBearer(ctx)
	default:
		return fmt.Errorf("不支持的认证方式: %s", sm.credentials.Method)
	}
}

// loginWithForm 表单登录
func (sm *SessionManager) loginWithForm(ctx context.Context) error {
	if sm.credentials.LoginURL == "" {
		return fmt.Errorf("表单登录需要指定LoginURL")
	}

	// 第一步：先访问登录页面获取初始会话
	fmt.Printf("[DEBUG] 首先访问登录页面获取初始会话\n")
	getReq, err := http.NewRequestWithContext(ctx, "GET", sm.credentials.LoginURL, nil)
	if err != nil {
		return fmt.Errorf("创建GET登录页面请求失败: %w", err)
	}
	getReq.Header.Set("User-Agent", "DroneRiskScan/1.0")
	
	getResp, err := sm.httpClient.Do(getReq)
	if err != nil {
		return fmt.Errorf("访问登录页面失败: %w", err)
	}
	defer getResp.Body.Close()
	
	// 保存初始Cookie
	initialCookies := getResp.Cookies()
	fmt.Printf("[DEBUG] 初始访问获得 %d 个Cookie\n", len(initialCookies))
	for _, cookie := range initialCookies {
		fmt.Printf("[DEBUG] 初始Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}

	// 准备登录数据
	loginData := url.Values{}
	
	// 添加用户名密码 (bWAPP使用 login/password 字段)
	if sm.credentials.Username != "" {
		loginData.Set("login", sm.credentials.Username)
	}
	if sm.credentials.Password != "" {
		loginData.Set("password", sm.credentials.Password)
	}
	
	// 添加其他必要的字段
	for key, value := range sm.credentials.LoginData {
		loginData.Set(key, value)
	}
	
	// 调试：打印POST数据
	fmt.Printf("[DEBUG] 登录POST数据: %s\n", loginData.Encode())
	fmt.Printf("[DEBUG] 登录POST URL: %s\n", sm.credentials.LoginURL)

	// 发送登录请求
	req, err := http.NewRequestWithContext(ctx, "POST", sm.credentials.LoginURL, strings.NewReader(loginData.Encode()))
	if err != nil {
		return fmt.Errorf("创建登录请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "DroneRiskScan/1.0")
	
	// 添加初始访问获得的Cookie
	for _, cookie := range initialCookies {
		req.AddCookie(cookie)
	}

	resp, err := sm.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 保存Cookie并修复Domain问题
	loginResponseCookies := resp.Cookies()
	fmt.Printf("[DEBUG] 登录响应返回了 %d 个Cookie\n", len(loginResponseCookies))
	
	// 合并初始Cookie和登录响应Cookie
	allCookies := make([]*http.Cookie, 0)
	
	// 首先添加初始Cookie
	for _, cookie := range initialCookies {
		allCookies = append(allCookies, cookie)
	}
	
	// 然后添加或更新登录响应的Cookie
	for _, loginCookie := range loginResponseCookies {
		// 检查是否已经存在同名Cookie，如果存在则更新
		updated := false
		for i, existing := range allCookies {
			if existing.Name == loginCookie.Name {
				allCookies[i] = loginCookie
				updated = true
				break
			}
		}
		// 如果不存在，则添加新Cookie
		if !updated {
			allCookies = append(allCookies, loginCookie)
		}
	}
	
	sm.cookies = allCookies
	fmt.Printf("[DEBUG] 登录后总共保存了 %d 个Cookie:\n", len(sm.cookies))
	for _, cookie := range sm.cookies {
		fmt.Printf("[DEBUG]   %s=%s (Domain: %s, Path: %s)\n", cookie.Name, cookie.Value, cookie.Domain, cookie.Path)
		
		// 修复Cookie Domain问题
		if cookie.Domain == "" {
			if strings.Contains(sm.credentials.LoginURL, "127.0.0.1") {
				cookie.Domain = "127.0.0.1"
			} else if strings.Contains(sm.credentials.LoginURL, "localhost") {
				cookie.Domain = "localhost"
			}
			fmt.Printf("[DEBUG] 修复Cookie Domain: %s -> %s\n", cookie.Name, cookie.Domain)
		}
		
		// 确保Path设置
		if cookie.Path == "" {
			cookie.Path = "/"
		}
	}
	
	// 为bWAPP添加security_level=0 Cookie
	if strings.Contains(sm.credentials.LoginURL, "127.0.0.1") || strings.Contains(sm.credentials.LoginURL, "bwapp") {
		securityLevelCookie := &http.Cookie{
			Name:   "security_level",
			Value:  "0",
			Path:   "/",
			Domain: "127.0.0.1", // 明确设置Domain
		}
		sm.cookies = append(sm.cookies, securityLevelCookie)
		fmt.Printf("[DEBUG] 为bWAPP添加security_level=0 Cookie (Domain: 127.0.0.1)\n")
	}

	// 读取响应内容（处理gzip压缩）
	helper := transport.NewResponseHelper()
	body, err := helper.ReadBody(resp)
	if err != nil {
		return fmt.Errorf("读取登录响应失败: %w", err)
	}

	responseText := string(body)

	// 调试：保存登录响应内容
	fmt.Printf("[DEBUG] 登录响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("[DEBUG] 登录响应长度: %d bytes\n", len(responseText))
	if len(responseText) > 200 {
		fmt.Printf("[DEBUG] 登录响应内容开头: %s...\n", responseText[:200])
	}
	
	// 检查登录是否成功
	if sm.credentials.SuccessText != "" {
		if strings.Contains(responseText, sm.credentials.SuccessText) {
			fmt.Printf("[DEBUG] 找到成功登录标志: %s\n", sm.credentials.SuccessText)
			sm.isLoggedIn = true
			sm.loginTime = time.Now()
			return nil
		} else {
			fmt.Printf("[DEBUG] 未找到成功登录标志: %s\n", sm.credentials.SuccessText)
		}
	}

	// 检查登录失败标志
	if sm.credentials.FailureText != "" {
		if strings.Contains(responseText, sm.credentials.FailureText) {
			fmt.Printf("[DEBUG] 发现登录失败标志: %s\n", sm.credentials.FailureText)
			return fmt.Errorf("登录失败: 在响应中发现失败标志")
		}
	}

	// 检查是否包含登录表单（如果包含，说明登录失败了）
	if strings.Contains(responseText, "<title>bWAPP - Login</title>") || strings.Contains(responseText, "name=\"login\"") {
		fmt.Printf("[DEBUG] 登录响应包含登录表单，登录可能失败\n")
		return fmt.Errorf("登录失败: 响应仍然是登录页面")
	}

	// 如果有Cookie且状态码正常，认为登录成功
	if len(sm.cookies) > 0 && resp.StatusCode < 400 {
		fmt.Printf("[DEBUG] 基于Cookie和状态码判断登录成功\n")
		sm.isLoggedIn = true
		sm.loginTime = time.Now()
		return nil
	}

	return fmt.Errorf("登录失败: 状态码 %d", resp.StatusCode)
}

// loginWithBasic HTTP Basic认证
func (sm *SessionManager) loginWithBasic(ctx context.Context) error {
	// HTTP Basic认证在每个请求中都会发送，无需预先登录
	sm.isLoggedIn = true
	sm.loginTime = time.Now()
	return nil
}

// loginWithCookie 使用现有Cookie
func (sm *SessionManager) loginWithCookie(ctx context.Context) error {
	if sm.credentials.Cookies == nil || len(sm.credentials.Cookies) == 0 {
		return fmt.Errorf("Cookie认证需要提供cookies")
	}

	// 将字符串Cookie转换为http.Cookie
	for name, value := range sm.credentials.Cookies {
		cookie := &http.Cookie{
			Name:  name,
			Value: value,
		}
		sm.cookies = append(sm.cookies, cookie)
	}

	sm.isLoggedIn = true
	sm.loginTime = time.Now()
	return nil
}

// loginWithBearer Bearer Token认证
func (sm *SessionManager) loginWithBearer(ctx context.Context) error {
	if sm.credentials.Token == "" {
		return fmt.Errorf("Bearer认证需要提供token")
	}

	sm.isLoggedIn = true
	sm.loginTime = time.Now()
	return nil
}

// ApplyAuth 为请求应用认证信息
func (sm *SessionManager) ApplyAuth(req *http.Request) error {
	if !sm.isLoggedIn {
		return fmt.Errorf("未登录，无法应用认证信息")
	}

	switch sm.credentials.Method {
	case AuthMethodBasic:
		req.SetBasicAuth(sm.credentials.Username, sm.credentials.Password)
	case AuthMethodBearer:
		req.Header.Set("Authorization", "Bearer "+sm.credentials.Token)
	case AuthMethodForm, AuthMethodCookie:
		// 添加Cookie
		for _, cookie := range sm.cookies {
			req.AddCookie(cookie)
		}
	}

	return nil
}

// IsLoggedIn 检查是否已登录
func (sm *SessionManager) IsLoggedIn() bool {
	return sm.isLoggedIn
}

// GetCookies 获取当前Cookie
func (sm *SessionManager) GetCookies() []*http.Cookie {
	return sm.cookies
}

// GetSessionID 获取会话ID
func (sm *SessionManager) GetSessionID() string {
	// 从Cookie中查找常见的会话ID
	for _, cookie := range sm.cookies {
		if strings.Contains(strings.ToLower(cookie.Name), "session") ||
			strings.Contains(strings.ToLower(cookie.Name), "phpsessid") ||
			strings.Contains(strings.ToLower(cookie.Name), "jsessionid") {
			return cookie.Value
		}
	}
	return sm.sessionID
}

// Logout 登出
func (sm *SessionManager) Logout(ctx context.Context) error {
	sm.isLoggedIn = false
	sm.cookies = make([]*http.Cookie, 0)
	sm.sessionID = ""
	sm.loginTime = time.Time{}
	return nil
}

// GetLoginDuration 获取登录持续时间
func (sm *SessionManager) GetLoginDuration() time.Duration {
	if !sm.isLoggedIn {
		return 0
	}
	return time.Since(sm.loginTime)
}