package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"web-scanner/internal/crawler"
)

func main() {
	var (
		url         string
		maxDepth    int
		concurrency int
		timeout     int
	)

	flag.StringVar(&url, "url", "", "起始URL")
	flag.IntVar(&maxDepth, "depth", 2, "最大爬取深度")
	flag.IntVar(&concurrency, "concurrency", 5, "并发数")
	flag.IntVar(&timeout, "timeout", 10, "请求超时时间(秒)")

	flag.Parse()

	if url == "" {
		log.Fatal("请提供起始URL: -url http://example.com")
	}

	// 创建爬虫配置
	config := &crawler.Config{
		MaxDepth:    maxDepth,
		Concurrency: concurrency,
		Timeout:     timeout,
	}

	// 创建爬虫实例
	c := crawler.NewCrawler(config)

	fmt.Printf("开始爬取: %s (最大深度: %d)\n", url, maxDepth)

	// 开始爬取
	results := c.Crawl(url)

	// 处理结果
	for result := range results {
		fmt.Printf("\nURL: %s\n", result.URL)
		fmt.Printf("状态码: %d\n", result.StatusCode)
		fmt.Printf("标题: %s\n", result.Title)
		
		if len(result.Links) > 0 {
			fmt.Printf("发现 %d 个链接:\n", len(result.Links))
			for _, link := range result.Links {
				fmt.Printf("  - %s\n", link)
			}
		}
		
		if len(result.Forms) > 0 {
			fmt.Printf("发现 %d 个表单:\n", len(result.Forms))
			for i, form := range result.Forms {
				fmt.Printf("  表单 %d:\n", i+1)
				fmt.Printf("    动作: %s\n", form.Action)
				fmt.Printf("    方法: %s\n", form.Method)
				for _, input := range form.Inputs {
					fmt.Printf("      输入: %s (%s)\n", input.Name, input.Type)
				}
			}
		}
		
		fmt.Println(strings.Repeat("-", 50))
	}

	fmt.Println("爬取完成!")
}