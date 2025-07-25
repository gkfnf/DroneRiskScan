package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
	"context"
	"github.com/PuerkitoBio/goquery"
	"web-scanner/internal/scanner"
)

// Crawler 爬虫结构体
type Crawler struct {
	client      *scanner.HTTPClient
	maxDepth    int
	visited     map[string]bool
	visitedLock sync.RWMutex
	links       chan string
	results     chan CrawlResult
	concurrency int
	timeout     int
}

// CrawlResult 爬取结果
type CrawlResult struct {
	URL        string
	StatusCode int
	Links      []string
	Forms      []Form
	Title      string
}

// Form 表单信息
type Form struct {
	Action string
	Method string
	Inputs []Input
}

// Input 表单输入字段
type Input struct {
	Name  string
	Type  string
	Value string
}

// Config 爬虫配置
type Config struct {
	MaxDepth    int
	Concurrency int
	Timeout     int
}

// NewCrawler 创建新的爬虫实例
func NewCrawler(config *Config) *Crawler {
	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	maxDepth := config.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 2
	}

	return &Crawler{
		client:      scanner.NewHTTPClient(timeout),
		maxDepth:    maxDepth,
		visited:     make(map[string]bool),
		links:       make(chan string, 100),
		results:     make(chan CrawlResult, 100),
		concurrency: concurrency,
		timeout:     timeout,
	}
}

// Crawl 开始爬取
func (c *Crawler) Crawl(startURL string) <-chan CrawlResult {
	// 添加起始URL
	c.links <- fmt.Sprintf("%s|%d", startURL, 0)
	
	// 启动工作者协程
	var wg sync.WaitGroup
	for i := 0; i < c.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 为每个worker设置一个总的超时时间，防止永久阻塞
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout*10)*time.Second)
			defer cancel()
			c.workerWithContext(ctx)
		}()
	}

	// 启动结果收集协程
	resultChan := make(chan CrawlResult, 100)
	go func() {
		defer close(resultChan)
		wg.Wait()
		close(c.links)
	}()

	// 启动结果转发协程
	go func() {
		for result := range c.results {
			select {
			case resultChan <- result:
			case <-time.After(1 * time.Second):
				// 防止结果通道阻塞
			}
		}
	}()

	return resultChan
}

// workerWithContext 带上下文的爬虫工作协程
func (c *Crawler) workerWithContext(ctx context.Context) {
	for {
		select {
		case link, ok := <-c.links:
			if !ok {
				return // 通道已关闭
			}
			
			parts := strings.SplitN(link, "|", 2)
			if len(parts) != 2 {
				continue
			}

			url := parts[0]
			var depth int
			fmt.Sscanf(parts[1], "%d", &depth)

			// 检查深度限制
			if depth > c.maxDepth {
				continue
			}

			// 检查是否已访问
			if c.isVisited(url) {
				continue
			}

			// 标记为已访问
			c.markVisited(url)

			// 爬取页面
			result, err := c.crawlPage(url)
			if err != nil {
				// 对于超时错误，我们记录但不中断整个爬取过程
				fmt.Printf("爬取 %s 时出错: %v\n", url, err)
				continue
			}

			// 发送结果
			select {
			case c.results <- *result:
			case <-time.After(1 * time.Second):
				// 如果结果通道已满且超时，跳过结果
			}

			// 添加新发现的链接到队列
			if depth < c.maxDepth {
				for _, newLink := range result.Links {
					// 将新链接加入队列，深度加1
					// 避免添加已访问的链接
					if !c.isVisited(newLink) {
						select {
						case c.links <- fmt.Sprintf("%s|%d", newLink, depth+1):
						case <-time.After(1 * time.Second):
							// 如果链接通道已满且超时，跳过链接
						}
					}
				}
			}
		case <-ctx.Done():
			// 上下文超时或取消
			return
		}
	}
}

// crawlPage 爬取单个页面
func (c *Crawler) crawlPage(pageURL string) (*CrawlResult, error) {
	// 发送HTTP请求
	httpResp, err := c.client.Get(pageURL)
	if err != nil {
		return nil, err
	}

	// 解析HTML并提取信息
	result, err := c.parseHTML(string(httpResp.Body), pageURL)
	if err != nil {
		return nil, err
	}

	result.URL = pageURL
	result.StatusCode = httpResp.StatusCode

	return result, nil
}

// isVisited 检查URL是否已访问
func (c *Crawler) isVisited(url string) bool {
	c.visitedLock.RLock()
	defer c.visitedLock.RUnlock()
	return c.visited[url]
}

// markVisited 标记URL为已访问
func (c *Crawler) markVisited(url string) {
	c.visitedLock.Lock()
	defer c.visitedLock.Unlock()
	c.visited[url] = true
}

// parseHTML 解析HTML文档并提取信息
func (c *Crawler) parseHTML(htmlContent string, baseURL string) (*CrawlResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	result := &CrawlResult{
		Links: []string{},
		Forms: []Form{},
	}

	// 提取页面标题
	result.Title = strings.TrimSpace(doc.Find("title").Text())

	// 提取所有链接
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// 解析并规范化URL
		resolvedURL := c.resolveURL(baseURL, href)
		if resolvedURL != "" {
			result.Links = append(result.Links, resolvedURL)
		}
	})

	// 提取所有表单
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		form := Form{}
		
		// 获取表单动作和方法
		form.Action, _ = s.Attr("action")
		form.Method, _ = s.Attr("method")
		if form.Method == "" {
			form.Method = "GET"
		}
		
		// 解析表单动作URL
		if form.Action != "" {
			form.Action = c.resolveURL(baseURL, form.Action)
		} else {
			// 如果没有指定action，则使用当前页面URL作为action
			form.Action = baseURL
		}
		
		// 提取表单输入字段
		s.Find("input").Each(func(j int, input *goquery.Selection) {
			inputField := Input{}
			inputField.Name, _ = input.Attr("name")
			inputField.Type, _ = input.Attr("type")
			inputField.Value, _ = input.Attr("value")
			
			if inputField.Type == "" {
				inputField.Type = "text"
			}
			
			form.Inputs = append(form.Inputs, inputField)
		})
		
		result.Forms = append(result.Forms, form)
	})

	return result, nil
}

// resolveURL 解析并规范化URL
func (c *Crawler) resolveURL(baseURL, href string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	resolved := base.ResolveReference(ref)
	
	// 只处理HTTP/HTTPS协议
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return ""
	}
	
	// 只处理同域URL，避免爬取外部网站
	baseHost := base.Host
	refHost := resolved.Host
	
	// 如果是相对路径或者相同主机名
	if refHost == "" || refHost == baseHost {
		return resolved.String()
	}
	
	// 对于外部链接，我们也可以选择不处理
	return ""
}