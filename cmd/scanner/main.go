package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"web-scanner/internal/scanner"
)

func main() {
	var (
		targets     string
		threads     int
		timeout     int
		templateDir string
		output      string
	)

	flag.StringVar(&targets, "targets", "", "目标URL列表，用逗号分隔")
	flag.IntVar(&threads, "threads", 10, "并发线程数")
	flag.IntVar(&timeout, "timeout", 10, "请求超时时间(秒)")
	flag.StringVar(&templateDir, "templates", "./templates", "模板目录路径")
	flag.StringVar(&output, "output", "", "输出文件路径")

	flag.Parse()

	if targets == "" {
		log.Fatal("请提供扫描目标: -targets url1,url2,...")
	}

	targetList := strings.Split(targets, ",")
	if len(targetList) == 0 {
		log.Fatal("无效的目标列表")
	}

	// 初始化扫描器
	sc := scanner.NewScanner(&scanner.Config{
		Threads:     threads,
		Timeout:     timeout,
		TemplateDir: templateDir,
		Output:      output,
	})

	fmt.Printf("开始扫描 %d 个目标...\n", len(targetList))
	
	// 执行扫描
	results, err := sc.Scan(targetList)
	if err != nil {
		log.Fatalf("扫描过程中出现错误: %v", err)
	}

	// 输出结果
	if err := sc.Output(results); err != nil {
		log.Fatalf("输出结果时出现错误: %v", err)
	}

	fmt.Println("扫描完成!")
}