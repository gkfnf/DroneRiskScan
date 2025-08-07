package browser

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// PlaywrightManager manages Playwright browser automation
type PlaywrightManager struct {
	pw         *playwright.Playwright
	browser    playwright.Browser
	context    playwright.BrowserContext
	page       playwright.Page
	config     *StagehandConfig
	cookies    []*http.Cookie
	sessionID  string
}

// NewPlaywrightManager creates a new Playwright manager
func NewPlaywrightManager(config *StagehandConfig) *PlaywrightManager {
	if config == nil {
		config = DefaultStagehandConfig()
	}
	
	return &PlaywrightManager{
		config:  config,
		cookies: make([]*http.Cookie, 0),
	}
}

// Start initializes Playwright and launches browser
func (pm *PlaywrightManager) Start(ctx context.Context) error {
	// Install Playwright browsers if needed
	err := playwright.Install()
	if err != nil {
		return fmt.Errorf("failed to install Playwright: %w", err)
	}
	
	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to start Playwright: %w", err)
	}
	pm.pw = pw
	
	// Launch browser
	browserOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(pm.config.Headless),
	}
	
	browser, err := pw.Chromium.Launch(browserOptions)
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}
	pm.browser = browser
	
	// Create browser context
	contextOptions := playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(pm.config.UserAgent),
		Viewport: &playwright.Size{
			Width:  pm.config.Viewport.Width,
			Height: pm.config.Viewport.Height,
		},
		IgnoreHttpsErrors: playwright.Bool(true),
	}
	
	context, err := browser.NewContext(contextOptions)
	if err != nil {
		return fmt.Errorf("failed to create browser context: %w", err)
	}
	pm.context = context
	
	// Create page
	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	pm.page = page
	
	fmt.Println("[INFO] Playwright browser automation started successfully")
	return nil
}

// PerformBWAPPAuthentication performs authentication for bWAPP
func (pm *PlaywrightManager) PerformBWAPPAuthentication(ctx context.Context, credentials map[string]string) (*InteractionResult, error) {
	result := &InteractionResult{
		SessionData: make(map[string]interface{}),
		NetworkLogs: make([]*NetworkLog, 0),
	}
	
	startTime := time.Now()
	
	// Navigate to login page
	loginURL := credentials["login_url"]
	if loginURL == "" {
		loginURL = "http://127.0.0.1/login.php"
	}
	
	fmt.Printf("[DEBUG] Navigating to login page: %s\n", loginURL)
	_, err := pm.page.Goto(loginURL)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to navigate to login page: %v", err)
		return result, err
	}
	
	// Wait for page to load
	time.Sleep(1 * time.Second)
	
	// Fill username
	fmt.Printf("[DEBUG] Filling username: %s\n", credentials["username"])
	err = pm.page.Fill("input[name='login']", credentials["username"])
	if err != nil {
		// Try alternative selectors
		err = pm.page.Fill("#login", credentials["username"])
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to fill username: %v", err)
			return result, err
		}
	}
	
	// Fill password
	fmt.Printf("[DEBUG] Filling password\n")
	err = pm.page.Fill("input[name='password']", credentials["password"])
	if err != nil {
		// Try alternative selectors
		err = pm.page.Fill("#password", credentials["password"])
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to fill password: %v", err)
			return result, err
		}
	}
	
	// Set security level to 0 (low) if it's a select element
	selects, _ := pm.page.QuerySelectorAll("select[name='security_level']")
	if len(selects) > 0 {
		fmt.Println("[DEBUG] Setting security level to 0 (low)")
		pm.page.SelectOption("select[name='security_level']", playwright.SelectOptionValues{Values: &[]string{"0"}})
	}
	
	// Click login button
	fmt.Println("[DEBUG] Clicking login button")
	err = pm.page.Click("input[type='submit'][value='Login']")
	if err != nil {
		// Try alternative selectors
		err = pm.page.Click("button[type='submit']")
		if err != nil {
			err = pm.page.Click("input[name='form']")
			if err != nil {
				result.Success = false
				result.Error = fmt.Sprintf("failed to click login button: %v", err)
				return result, err
			}
		}
	}
	
	// Wait for navigation
	time.Sleep(2 * time.Second)
	
	// Check if login was successful
	currentURL := pm.page.URL()
	pageContent, _ := pm.page.Content()
	
	fmt.Printf("[DEBUG] Current URL after login: %s\n", currentURL)
	
	// Check for success indicators
	if strings.Contains(pageContent, "Choose your bug") || 
	   strings.Contains(pageContent, "Portal") ||
	   strings.Contains(currentURL, "portal.php") {
		result.Success = true
		result.Message = "Authentication successful"
		fmt.Println("[INFO] bWAPP authentication successful")
		
		// Extract cookies
		cookies, err := pm.extractCookies(ctx)
		if err == nil {
			result.Cookies = cookies
			pm.cookies = cookies
			fmt.Printf("[DEBUG] Extracted %d cookies\n", len(cookies))
			for _, cookie := range cookies {
				fmt.Printf("[DEBUG] Cookie: %s=%s\n", cookie.Name, cookie.Value)
			}
		}
	} else if strings.Contains(pageContent, "Invalid credentials") ||
	          strings.Contains(pageContent, "Login failed") {
		result.Success = false
		result.Message = "Authentication failed: Invalid credentials"
		fmt.Println("[WARN] Authentication failed: Invalid credentials")
	} else {
		// Uncertain result, but proceed
		result.Success = true
		result.Message = "Authentication completed (status uncertain)"
		fmt.Println("[INFO] Authentication completed, extracting cookies")
		
		// Extract cookies anyway
		cookies, err := pm.extractCookies(ctx)
		if err == nil {
			result.Cookies = cookies
			pm.cookies = cookies
		}
	}
	
	result.Duration = time.Since(startTime)
	return result, nil
}

// extractCookies extracts cookies from the browser context
func (pm *PlaywrightManager) extractCookies(ctx context.Context) ([]*http.Cookie, error) {
	if pm.context == nil {
		return nil, fmt.Errorf("browser context not initialized")
	}
	
	playwrightCookies, err := pm.context.Cookies()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}
	
	httpCookies := make([]*http.Cookie, 0, len(playwrightCookies))
	for _, pwCookie := range playwrightCookies {
		httpCookie := &http.Cookie{
			Name:     pwCookie.Name,
			Value:    pwCookie.Value,
			Path:     pwCookie.Path,
			Domain:   pwCookie.Domain,
			Secure:   pwCookie.Secure,
			HttpOnly: pwCookie.HttpOnly,
		}
		
		if pwCookie.Expires != -1 {
			httpCookie.Expires = time.Unix(int64(pwCookie.Expires), 0)
		}
		
		// Set SameSite (simplified approach)
		httpCookie.SameSite = http.SameSiteDefaultMode
		
		httpCookies = append(httpCookies, httpCookie)
	}
	
	// Always add security_level=0 for bWAPP
	httpCookies = append(httpCookies, &http.Cookie{
		Name:  "security_level",
		Value: "0",
		Path:  "/",
	})
	
	return httpCookies, nil
}

// DiscoverFunctionPoints discovers function points on the current page
func (pm *PlaywrightManager) DiscoverFunctionPoints(ctx context.Context, targetURL string) ([]*FunctionPoint, error) {
	fmt.Printf("[INFO] Discovering function points for: %s\n", targetURL)
	
	// Check if Playwright is initialized
	if pm.page == nil {
		fmt.Println("[DEBUG] Playwright not initialized, starting...")
		err := pm.Start(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to start Playwright: %w", err)
		}
	}
	
	// Set existing cookies to maintain session
	if len(pm.cookies) > 0 {
		fmt.Printf("[DEBUG] Setting %d cookies for session\n", len(pm.cookies))
		for _, cookie := range pm.cookies {
			playwrightCookie := playwright.OptionalCookie{
				Name:  cookie.Name,
				Value: cookie.Value,
				URL:   &targetURL,
			}
			pm.context.AddCookies([]playwright.OptionalCookie{playwrightCookie})
		}
	}
	
	// Navigate to target URL
	_, err := pm.page.Goto(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to target: %w", err)
	}
	
	// Wait for page to load
	time.Sleep(1 * time.Second)
	
	functionPoints := make([]*FunctionPoint, 0)
	
	// Find all forms
	forms, err := pm.page.QuerySelectorAll("form")
	if err == nil {
		for _, form := range forms {
			action, _ := form.GetAttribute("action")
			method, _ := form.GetAttribute("method")
			if method == "" {
				method = "GET"
			}
			
			// Get form inputs
			inputs, _ := form.QuerySelectorAll("input, select, textarea")
			params := make(map[string]*ParamInfo)
			
			for _, input := range inputs {
				name, _ := input.GetAttribute("name")
				inputType, _ := input.GetAttribute("type")
				required, _ := input.GetAttribute("required")
				
				if name != "" {
					params[name] = &ParamInfo{
						Name:       name,
						Type:       inputType,
						Required:   required == "required" || required == "true",
						Injectable: true, // Mark all params as injectable for testing
					}
				}
			}
			
			if len(params) > 0 {
				fp := &FunctionPoint{
					URL:         action,
					Type:        "form",
					Method:      strings.ToUpper(method),
					Parameters:  params,
					Description: fmt.Sprintf("Form with %d parameters", len(params)),
				}
				functionPoints = append(functionPoints, fp)
				fmt.Printf("[DEBUG] Found form: %s %s with %d params\n", fp.Method, fp.URL, len(params))
			}
		}
	}
	
	// Find all links with parameters
	links, err := pm.page.QuerySelectorAll("a[href*='?']")
	if err == nil {
		for _, link := range links {
			href, _ := link.GetAttribute("href")
			if href != "" && strings.Contains(href, "?") {
				// Parse URL parameters
				parts := strings.Split(href, "?")
				if len(parts) > 1 {
					params := make(map[string]*ParamInfo)
					paramPairs := strings.Split(parts[1], "&")
					for _, pair := range paramPairs {
						kv := strings.Split(pair, "=")
						if len(kv) > 0 {
							params[kv[0]] = &ParamInfo{
								Name:       kv[0],
								Type:       "string",
								Injectable: true,
							}
						}
					}
					
					if len(params) > 0 {
						fp := &FunctionPoint{
							URL:         href,
							Type:        "link",
							Method:      "GET",
							Parameters:  params,
							Description: fmt.Sprintf("Link with %d parameters", len(params)),
						}
						functionPoints = append(functionPoints, fp)
						fmt.Printf("[DEBUG] Found link: %s with %d params\n", href, len(params))
					}
				}
			}
		}
	}
	
	fmt.Printf("[INFO] Discovered %d function points\n", len(functionPoints))
	return functionPoints, nil
}

// GetCookies returns the extracted cookies
func (pm *PlaywrightManager) GetCookies() []*http.Cookie {
	return pm.cookies
}

// Close closes the browser and cleans up resources
func (pm *PlaywrightManager) Close() error {
	if pm.page != nil {
		pm.page.Close()
	}
	if pm.context != nil {
		pm.context.Close()
	}
	if pm.browser != nil {
		pm.browser.Close()
	}
	if pm.pw != nil {
		pm.pw.Stop()
	}
	return nil
}