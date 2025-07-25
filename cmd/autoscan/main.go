package main

import (
	"flag"
	"fmt"
	"log"
	"web-scanner/internal/crawler"
	"web-scanner/internal/scanner"
)

func main() {
	var (
		url         string
		threads     int
		depth       int
		timeout     int
		templateDir string
		output      string
		maxTargets  int
	)

	flag.StringVar(&url, "url", "", "目标网站URL")
	flag.IntVar(&threads, "threads", 10, "扫描线程数")
	flag.IntVar(&depth, "depth", 2, "爬虫深度")
	flag.IntVar(&timeout, "timeout", 10, "请求超时时间(秒)")
	flag.StringVar(&templateDir, "templates", "./templates", "模板目录路径")
	flag.StringVar(&output, "output", "", "输出文件路径")
	flag.IntVar(&maxTargets, "max-targets", 50, "最大扫描目标数量")

	flag.Parse()

	if url == "" {
		log.Fatal("请提供目标网站URL: -url http://example.com")
	}

	fmt.Printf("开始自动扫描网站: %s\n", url)
	fmt.Printf("爬虫深度: %d, 扫描线程数: %d, 最大目标数: %d\n", depth, threads, maxTargets)

	// 第一步：使用爬虫发现网站结构
	fmt.Println("\n=== 第一步：爬取网站结构 ===")
	crawlerConfig := &crawler.Config{
		MaxDepth:    depth,
		Concurrency: 5,
		Timeout:     timeout,
	}

	// 创建爬虫实例
	c := crawler.NewCrawler(crawlerConfig)

	// 开始爬取
	crawlResults := c.Crawl(url)

	// 收集所有发现的URL
	var allURLs []string
	urlMap := make(map[string]bool) // 用于去重

	fmt.Println("发现的页面:")
	for result := range crawlResults {
		fmt.Printf("  %s (状态码: %d)\n", result.URL, result.StatusCode)
		if result.StatusCode == 200 {
			if _, exists := urlMap[result.URL]; !exists {
				allURLs = append(allURLs, result.URL)
				urlMap[result.URL] = true
			}
		}

		// 显示表单信息
		if len(result.Forms) > 0 {
			fmt.Printf("    发现 %d 个表单:\n", len(result.Forms))
			for i, form := range result.Forms {
				fmt.Printf("      表单 %d: %s (%s)\n", i+1, form.Action, form.Method)
			}
		}
	}

	// 限制目标数量
	if len(allURLs) > maxTargets {
		fmt.Printf("\n发现 %d 个有效页面，超过最大限制 %d，仅扫描前 %d 个\n", len(allURLs), maxTargets, maxTargets)
		allURLs = allURLs[:maxTargets]
	} else {
		fmt.Printf("\n总共发现 %d 个有效页面\n", len(allURLs))
	}

	// 第二步：对发现的URL进行漏洞扫描
	fmt.Println("\n=== 第二步：漏洞扫描 ===")
	
	// 创建扫描器配置
	scannerConfig := &scanner.Config{
		Threads:     threads,
		Timeout:     timeout,
		TemplateDir: templateDir,
		Output:      output,
	}

	// 创建扫描器实例
	s := scanner.NewScanner(scannerConfig)

	// 执行扫描
	fmt.Printf("开始扫描 %d 个页面...\n", len(allURLs))
	results, err := s.Scan(allURLs)
	if err != nil {
		log.Fatalf("扫描过程中出现错误: %v", err)
	}

	// 第三步：输出结果
	fmt.Println("\n=== 第三步：扫描结果 ===")
	if err := s.Output(results); err != nil {
		log.Fatalf("输出结果时出现错误: %v", err)
	}

	// 统计结果
	vulnCount := len(results)
	if vulnCount > 0 {
		fmt.Printf("\n发现 %d 个潜在漏洞!\n", vulnCount)
		
		// 按严重性分类统计
		severityCount := make(map[string]int)
		templateCount := make(map[string]int)
		
		for _, result := range results {
			severityCount[result.Severity]++
			templateCount[result.TemplateID]++
		}
		
		fmt.Println("\n按严重性分类:")
		for severity, count := range severityCount {
			fmt.Printf("  %s: %d\n", severity, count)
		}
		
		fmt.Println("\n按漏洞类型分类:")
		for template, count := range templateCount {
			fmt.Printf("  %s: %d\n", template, count)
		}
	} else {
		fmt.Println("\n未发现明显漏洞。")
	}

	fmt.Println("\n自动扫描完成!")
}