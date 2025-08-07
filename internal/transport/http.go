package transport

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dronesec/droneriskscan/pkg/models"
)

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	SetProxy(proxyURL string) error
	SetTimeout(timeout time.Duration)
	SetUserAgent(userAgent string)
	Close() error
}

// Client HTTP客户端实现
type Client struct {
	client    *http.Client
	userAgent string
	headers   map[string]string
	options   *ClientOptions
}

// ClientOptions 客户端选项
type ClientOptions struct {
	Timeout          time.Duration
	MaxRedirects     int
	InsecureSkipTLS  bool
	UserAgent        string
	Proxy            *models.ProxyConfig
	Headers          map[string]string
	MaxIdleConns     int
	MaxConnsPerHost  int
	IdleConnTimeout  time.Duration
	DisableKeepAlive bool
}

// DefaultClientOptions 默认客户端选项
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout:          30 * time.Second,
		MaxRedirects:     5,
		InsecureSkipTLS:  true,
		UserAgent:        "DroneRiskScan/1.0",
		Headers:          make(map[string]string),
		MaxIdleConns:     100,
		MaxConnsPerHost:  10,
		IdleConnTimeout:  90 * time.Second,
		DisableKeepAlive: false,
	}
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(options *ClientOptions) *Client {
	if options == nil {
		options = DefaultClientOptions()
	}

	// 创建传输配置
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.InsecureSkipTLS,
		},
		MaxIdleConns:        options.MaxIdleConns,
		MaxIdleConnsPerHost: options.MaxConnsPerHost,
		IdleConnTimeout:     options.IdleConnTimeout,
		DisableKeepAlives:   options.DisableKeepAlive,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	// 设置代理
	if options.Proxy != nil && options.Proxy.Enabled {
		proxyURL, err := url.Parse(options.Proxy.Address)
		if err == nil {
			if options.Proxy.Username != "" && options.Proxy.Password != "" {
				proxyURL.User = url.UserPassword(options.Proxy.Username, options.Proxy.Password)
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   options.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= options.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", options.MaxRedirects)
			}
			return nil
		},
	}

	client := &Client{
		client:    httpClient,
		userAgent: options.UserAgent,
		headers:   make(map[string]string),
		options:   options,
	}

	// 设置默认头部
	if options.Headers != nil {
		for k, v := range options.Headers {
			client.headers[k] = v
		}
	}

	return client
}

// Do 执行HTTP请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// 添加默认头部
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	
	for key, value := range c.headers {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}

	// 添加默认头部以模拟真实浏览器
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	}
	if req.Header.Get("Accept-Encoding") == "" {
		req.Header.Set("Accept-Encoding", "gzip, deflate")
	}
	if req.Header.Get("Connection") == "" && !c.options.DisableKeepAlive {
		req.Header.Set("Connection", "keep-alive")
	}

	return c.client.Do(req)
}

// Get 发送GET请求
func (c *Client) Get(targetURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建GET请求失败: %w", err)
	}
	
	return c.Do(req)
}

// Post 发送POST请求
func (c *Client) Post(targetURL, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建POST请求失败: %w", err)
	}
	
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	
	return c.Do(req)
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxyURL string) error {
	if proxyURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("解析代理URL失败: %w", err)
	}

	transport := c.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(parsedURL)
	
	return nil
}

// SetTimeout 设置超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetUserAgent 设置User-Agent
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// SetHeader 设置请求头
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// Close 关闭客户端
func (c *Client) Close() error {
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}

// RequestBuilder HTTP请求构建器
type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
	ctx     context.Context
}

// NewRequestBuilder 创建请求构建器
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		headers: make(map[string]string),
		ctx:     context.Background(),
	}
}

// Method 设置请求方法
func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.method = strings.ToUpper(method)
	return rb
}

// URL 设置请求URL
func (rb *RequestBuilder) URL(url string) *RequestBuilder {
	rb.url = url
	return rb
}

// Header 添加请求头
func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

// Headers 批量设置请求头
func (rb *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		rb.headers[k] = v
	}
	return rb
}

// Body 设置请求体
func (rb *RequestBuilder) Body(body io.Reader) *RequestBuilder {
	rb.body = body
	return rb
}

// Context 设置上下文
func (rb *RequestBuilder) Context(ctx context.Context) *RequestBuilder {
	rb.ctx = ctx
	return rb
}

// Build 构建HTTP请求
func (rb *RequestBuilder) Build() (*http.Request, error) {
	if rb.method == "" {
		rb.method = "GET"
	}
	if rb.url == "" {
		return nil, fmt.Errorf("URL不能为空")
	}

	req, err := http.NewRequestWithContext(rb.ctx, rb.method, rb.url, rb.body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 添加请求头
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// ResponseHelper 响应助手
type ResponseHelper struct{}

// NewResponseHelper 创建响应助手
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// ReadBody 读取响应体
func (rh *ResponseHelper) ReadBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("响应或响应体为空")
	}
	
	// 不要在这里关闭resp.Body，让调用者处理
	var reader io.Reader = resp.Body
	
	// 检查是否为gzip压缩
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("创建gzip读取器失败: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}
	
	// 读取全部数据
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	
	return data, nil
}

// GetContentType 获取Content-Type
func (rh *ResponseHelper) GetContentType(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	return resp.Header.Get("Content-Type")
}

// IsTextResponse 判断是否为文本响应
func (rh *ResponseHelper) IsTextResponse(resp *http.Response) bool {
	contentType := rh.GetContentType(resp)
	textTypes := []string{
		"text/html",
		"text/plain",
		"application/json",
		"application/xml",
		"application/javascript",
		"text/css",
	}
	
	for _, textType := range textTypes {
		if strings.Contains(strings.ToLower(contentType), textType) {
			return true
		}
	}
	return false
}

// GetResponseSize 获取响应大小
func (rh *ResponseHelper) GetResponseSize(resp *http.Response) int64 {
	if resp == nil {
		return 0
	}
	return resp.ContentLength
}

// HasHeader 检查响应头是否存在
func (rh *ResponseHelper) HasHeader(resp *http.Response, headerName string) bool {
	if resp == nil {
		return false
	}
	return resp.Header.Get(headerName) != ""
}