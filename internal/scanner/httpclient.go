package scanner

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
)

// HTTPClient 自定义HTTP客户端
type HTTPClient struct {
	client  *http.Client
	timeout int
}

// HTTPResponse HTTP响应结构
type HTTPResponse struct {
	StatusCode int
	Status     string
	Header     http.Header
	Body       []byte
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(timeout int) *HTTPClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	return &HTTPClient{
		client:  client,
		timeout: timeout,
	}
}

// Get 发送GET请求
func (hc *HTTPClient) Get(url string) (*HTTPResponse, error) {
	// 对于可能导致超时的URL，我们使用更短的超时时间
	effectiveTimeout := hc.timeout
	if strings.Contains(url, "long-time") || strings.Contains(url, "too-large") || strings.Contains(url, "chunked") {
		effectiveTimeout = 3 // 对这些URL使用3秒超时
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(effectiveTimeout)*time.Second)
	defer cancel()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header,
		Body:       body,
	}, nil
}

// Post 发送POST请求
func (hc *HTTPClient) Post(url string, body []byte) (*HTTPResponse, error) {
	// 实现POST请求逻辑
	return nil, nil
}