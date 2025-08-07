# Web扫描引擎开发指南

## 简介

使用Go语言开发Web漏洞扫描引擎的指南。本项目参考了Nuclei等开源工具的设计理念，采用模块化设计，支持高并发扫描和基于模板的漏洞检测。

## 项目结构

```
web-scanner/
├── cmd/
│   ├── scanner/
│   │   └── main.go
│   ├── crawler/
│   │   └── main.go
│   └── autoscan/
│       └── main.go
├── internal/
│   ├── scanner/
│   │   ├── config.go
│   │   ├── scanner.go
│   │   └── httpclient.go
│   ├── templates/
│   │   └── templates.go
│   └── crawler/
│       └── crawler.go
├── pkg/
├── templates/
│   ├── example.yaml
│   ├── xss.yaml
│   └── sqli.yaml
├── go.mod
└── README.md
```

## 核心组件

1. **HTTP客户端** - 负责发送HTTP请求和处理响应
2. **扫描引擎** - 控制扫描流程和并发执行
3. **模板引擎** - 解析和执行扫描模板
4. **爬虫引擎** - 发现网站结构和内容
5. **结果处理** - 处理和输出扫描结果

## 开发步骤

1. 初始化Go模块
2. 实现HTTP客户端
3. 构建扫描引擎
4. 添加模板支持
5. 实现并发处理
6. 添加结果输出功能

## 使用方法

### 编译

```bash
# 编译扫描器
go build -o web-scanner cmd/scanner/main.go

# 编译爬虫
go build -o web-crawler cmd/crawler/main.go

# 编译自动扫描器
go build -o web-autoscan cmd/autoscan/main.go
```

### 运行扫描器

```bash
# 扫描单个目标
./web-scanner -targets http://example.com

# 扫描多个目标
./web-scanner -targets http://example.com,http://test.com

# 设置并发线程数
./web-scanner -targets http://example.com -threads 20

# 设置超时时间
./web-scanner -targets http://example.com -timeout 15

# 指定模板目录
./web-scanner -targets http://example.com -templates ./templates

# 输出到文件
./web-scanner -targets http://example.com -output result.txt
```

### 运行爬虫

```bash
# 爬取网站
./web-crawler -url http://example.com

# 设置爬取深度
./web-crawler -url http://example.com -depth 3

# 设置并发数
./web-crawler -url http://example.com -concurrency 10

# 设置超时时间
./web-crawler -url http://example.com -timeout 15
```

### 运行自动扫描器

```bash
# 自动扫描网站（先爬取再扫描）
./web-autoscan -url http://example.com

# 设置扫描线程数
./web-autoscan -url http://example.com -threads 20

# 设置爬虫深度
./web-autoscan -url http://example.com -depth 3

# 设置超时时间
./web-autoscan -url http://example.com -timeout 15

# 指定模板目录
./web-autoscan -url http://example.com -templates ./templates

# 输出到文件
./web-autoscan -url http://example.com -output result.txt
```

## 模板格式

模板使用YAML格式定义，包含以下主要部分：

- `id`: 模板唯一标识
- `info`: 模板信息（名称、作者、严重性等）
- `http`: HTTP请求定义
- `matchers`: 响应匹配规则

示例模板:

```yaml
id: example-template

info:
  name: "示例模板"
  author: "developer"
  severity: "info"
  description: "这是一个示例模板，用于演示模板格式"
  tags: ["example", "test"]

http:
  - method: GET
    path:
      - "/test"
      - "/demo"
    headers:
      User-Agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    matchers:
      - type: word
        part: body
        words:
          - "test page"
          - "demo content"
      
      - type: status
        status:
          - 200
          - 403
```

## XSS检测

项目包含基本的XSS检测模板，可以检测常见的跨站脚本漏洞：

- 在URL参数中注入脚本标签
- 测试HTML上下文中的XSS
- 检测事件处理程序中的XSS

示例XSS模板:

```yaml
id: xss-basic

info:
  name: "基本XSS检测"
  author: "scanner"
  severity: "high"
  description: "检测基本的跨站脚本攻击漏洞"
  tags: ["xss", "client-side", "owasp-a3"]

http:
  - method: GET
    path:
      - "{{BaseURL}}/?q=%3Cscript%3Ealert%281%29%3C%2Fscript%3E"
      - "{{BaseURL}}/?search=%27%3E%3Cscript%3Ealert%281%29%3C%2Fscript%3E"
      - "{{BaseURL}}/?input=%3Cimg%20src=x%20onerror=alert%281%29%3E"
    headers:
      User-Agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    matchers:
      - type: word
        part: body
        words:
          - "<script>alert(1)</script>"
          - "'>alert(1)"
          - "<img src=x onerror=alert(1)>"
      
      - type: word
        part: header
        words:
          - "text/html"
```

## SQL注入检测

项目包含基本的SQL注入检测模板，可以检测常见的SQL注入漏洞：

- 测试单引号注入
- 检测布尔盲注
- 检测联合查询注入

示例SQL注入模板:

```yaml
id: sqli-basic

info:
  name: "基本SQL注入检测"
  author: "scanner"
  severity: "high"
  description: "检测基本的SQL注入漏洞"
  tags: ["sqli", "database", "owasp-a1"]

http:
  - method: GET
    path:
      - "{{BaseURL}}/?id=1'"
      - "{{BaseURL}}/?id=1%20AND%201=1"
      - "{{BaseURL}}/?id=1%20AND%201=2"
      - "{{BaseURL}}/?id=1%20OR%201=1"
      - "{{BaseURL}}/?id=1%20UNION%20SELECT%20NULL"
    headers:
      User-Agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    matchers:
      - type: word
        part: body
        words:
          - "SQL syntax"
          - "mysql_fetch"
          - "mysql_num_rows"
          - "ORA-01756"
          - "ORA-00933"
          - "PostgreSQL"
          - "ODBC Microsoft Access"
          - "Unclosed quotation mark"
          - "quoted string not properly terminated"
          - "You have an error in your SQL syntax"
      
      - type: status
        status:
          - 500
          - 503
```

## 爬虫功能

爬虫模块可以自动发现网站的结构和内容，包括：

- 遍历网站链接
- 提取页面标题
- 发现表单和输入字段
- 支持深度控制
- 并发爬取

爬虫的主要特性：

1. **并发爬取** - 支持多协程并发爬取，提高效率
2. **深度控制** - 可设置最大爬取深度，避免无限爬取
3. **链接发现** - 自动发现页面中的所有链接
4. **表单识别** - 识别页面中的表单及其输入字段
5. **URL规范化** - 自动解析和规范化相对URL和绝对URL

使用示例：

```bash
# 爬取网站，最大深度为3
./web-crawler -url http://example.com -depth 3
```

输出示例：
```
开始爬取: http://example.com (最大深度: 3)

URL: http://example.com
状态码: 200
标题: 示例页面
发现 3 个链接:
  - http://example.com/about
  - http://example.com/contact
  - https://example.com/external
发现 2 个表单:
  表单 1:
    动作: http://example.com/search
    方法: GET
      输入: q (text)
      输入: submit (submit)
  表单 2:
    动作: http://example.com/login
    方法: POST
      输入: username (text)
      输入: password (password)
      输入: submit (submit)
```

## 自动扫描功能

自动扫描功能集成了爬虫和漏洞扫描，可以一键完成整个安全检测流程：

1. **自动发现** - 使用爬虫自动发现网站结构和页面
2. **智能扫描** - 对发现的所有页面进行漏洞扫描
3. **结果聚合** - 汇总所有扫描结果并分类展示

使用示例：

```bash
# 自动扫描网站
./web-autoscan -url http://example.com

# 设置爬虫深度和扫描线程数
./web-autoscan -url http://example.com -depth 3 -threads 20
```

工作流程：
```
开始自动扫描网站: http://example.com
爬虫深度: 2, 扫描线程数: 10

=== 第一步：爬取网站结构 ===
发现的页面:
  http://example.com (状态码: 200)
  http://example.com/about (状态码: 200)
  http://example.com/contact (状态码: 200)

总共发现 3 个有效页面

=== 第二步：漏洞扫描 ===
开始扫描 3 个页面...

=== 第三步：扫描结果 ===
目标: http://example.com
漏洞: 基本XSS检测
模板: xss-basic
严重性: high
请求: GET http://example.com/?q=<script>alert(1)</script>
响应: 200 OK
时间: 2023-01-01 12:00:00
---

发现 1 个潜在漏洞!

按严重性分类:
  high: 1

按漏洞类型分类:
  xss-basic: 1

自动扫描完成!
```

## 扩展功能建议

1. **支持更多协议**: DNS、TCP、SSL等
2. **增强模板功能**: 支持变量、表达式、条件判断等
3. **添加插件系统**: 支持自定义插件扩展功能
4. **优化性能**: 连接池、缓存机制等
5. **增加报告功能**: 生成HTML、JSON等格式报告
6. **支持认证**: Basic Auth、JWT等认证方式
7. **改进匹配器**: 支持正则表达式、更复杂的匹配逻辑
8. **添加被动扫描**: 支持代理模式和流量拦截
9. **增强爬虫功能**: 支持JavaScript渲染、处理Cookies等
10. **智能去重**: 更好的URL去重和内容去重机制
