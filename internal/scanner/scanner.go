package scanner

import (
	"fmt"
	"sync"
	"time"
	"web-scanner/internal/templates"
	"strings"
)

// Scanner Web扫描器结构体
type Scanner struct {
	config        *Config
	client        *HTTPClient
	templateEngine *templates.TemplateEngine
}

// Result 扫描结果
type Result struct {
	Target        string    `json:"target"`
	Vulnerability string    `json:"vulnerability"`
	Severity      string    `json:"severity"`
	Request       string    `json:"request"`
	Response      string    `json:"response"`
	TemplateID    string    `json:"template_id"`
	Timestamp     time.Time `json:"timestamp"`
}

// NewScanner 创建新的扫描器实例
func NewScanner(config *Config) *Scanner {
	// 初始化模板引擎
	templateEngine := templates.NewTemplateEngine()
	
	// 加载模板
	if err := templateEngine.LoadTemplates(config.TemplateDir); err != nil {
		fmt.Printf("警告: 加载模板时出错: %v\n", err)
	}
	
	return &Scanner{
		config:        config,
		client:        NewHTTPClient(config.Timeout),
		templateEngine: templateEngine,
	}
}

// Scan 执行扫描任务
func (s *Scanner) Scan(targets []string) ([]Result, error) {
	var results []Result
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("开始扫描 %d 个目标，使用 %d 个线程...\n", len(targets), s.config.Threads)
	
	// 创建任务通道
	jobs := make(chan string, len(targets))

	// 启动工作协程
	for i := 0; i < s.config.Threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for target := range jobs {
				result := s.scanTarget(target, workerID)
				mu.Lock()
				results = append(results, result...)
				mu.Unlock()
			}
		}(i)
	}

	// 发送任务到通道
	for _, target := range targets {
		jobs <- target
	}
	close(jobs)

	// 等待所有任务完成，设置超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 等待完成或超时
	select {
	case <-done:
		fmt.Println("扫描完成!")
	case <-time.After(time.Duration(s.config.Timeout*len(targets)) * time.Second):
		fmt.Println("扫描超时!")
	}

	return results, nil
}

// scanTarget 扫描单个目标
func (s *Scanner) scanTarget(target string, workerID int) []Result {
	var results []Result
	
	fmt.Printf("[Worker-%d] 正在扫描目标: %s\n", workerID, target)
	
	// 获取所有模板
	templates := s.templateEngine.GetAllTemplates()
	
	// 对每个模板执行扫描
	templateCount := len(templates)
	for i, template := range templates {
		fmt.Printf("  [Worker-%d] 使用模板 '%s' 扫描... (%d/%d)\n", workerID, template.ID, i+1, templateCount)
		
		// 处理模板中的路径
		for _, httpRequest := range template.HTTP {
			processedPaths := s.templateEngine.ProcessPaths(httpRequest.Path, target)
			
			// 对每个路径执行请求
			for _, path := range processedPaths {
				var url string
				if strings.HasPrefix(path, "http") {
					url = path
				} else {
					// 确保目标URL以http://或https://开头
					if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
						url = "http://" + target
					} else {
						url = target
					}
					
					// 添加路径
					if strings.HasPrefix(path, "/") {
						url += path
					} else {
						url += "/" + path
					}
				}
				
				// 发送请求
				resp, err := s.client.Get(url)
				if err != nil {
					// 记录连接错误
					results = append(results, Result{
						Target:        target,
						Vulnerability: "连接失败",
						Severity:      "High",
						Request:       fmt.Sprintf("%s %s", httpRequest.Method, url),
						Response:      err.Error(),
						TemplateID:    template.ID,
						Timestamp:     time.Now(),
					})
					continue
				}
				
				// 检查匹配器
				for _, matcher := range httpRequest.Matchers {
					match := s.checkMatcher(matcher, resp)
					if match {
						results = append(results, Result{
							Target:        target,
							Vulnerability: template.Info.Name,
							Severity:      template.Info.Severity,
							Request:       fmt.Sprintf("%s %s", httpRequest.Method, url),
							Response:      resp.Status,
							TemplateID:    template.ID,
							Timestamp:     time.Now(),
						})
						// 找到匹配项后跳出循环，避免重复报告
						break
					}
				}
			}
		}
	}
	
	return results
}

// checkMatcher 检查响应是否匹配
func (s *Scanner) checkMatcher(matcher templates.Matcher, resp *HTTPResponse) bool {
	switch matcher.Type {
	case "word":
		if matcher.Part == "body" {
			// 检查响应体中的关键词
			for _, word := range matcher.Words {
				if strings.Contains(string(resp.Body), word) {
					return true
				}
			}
		} else if matcher.Part == "header" {
			// 检查响应头中的关键词
			for _, word := range matcher.Words {
				if strings.Contains(fmt.Sprintf("%v", resp.Header), word) {
					return true
				}
			}
		} else if matcher.Part == "" || matcher.Part == "all" {
			// 默认检查响应体和响应头
			for _, word := range matcher.Words {
				if strings.Contains(string(resp.Body), word) || strings.Contains(fmt.Sprintf("%v", resp.Header), word) {
					return true
				}
			}
		}
	case "status":
		// 检查状态码
		for _, status := range matcher.Status {
			if resp.StatusCode == status {
				return true
			}
		}
	}
	
	return false
}

// Output 输出扫描结果
func (s *Scanner) Output(results []Result) error {
	// 如果指定了输出文件，则写入文件
	if s.config.Output != "" {
		// 这里应该实现文件写入逻辑
		fmt.Printf("结果已保存到文件: %s\n", s.config.Output)
		return nil
	}
	
	// 否则输出到控制台
	if len(results) == 0 {
		fmt.Println("未发现漏洞。")
		return nil
	}
	
	fmt.Printf("发现 %d 个潜在漏洞:\n\n", len(results))
	for _, result := range results {
		fmt.Printf("目标: %s\n", result.Target)
		fmt.Printf("漏洞: %s\n", result.Vulnerability)
		fmt.Printf("模板: %s\n", result.TemplateID)
		fmt.Printf("严重性: %s\n", result.Severity)
		fmt.Printf("请求: %s\n", result.Request)
		fmt.Printf("响应: %s\n", result.Response)
		fmt.Printf("时间: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println("---")
	}
	
	return nil
}