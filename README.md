# DroneRiskScan - 专业Web安全扫描引擎

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

## 🚀 项目简介

DroneRiskScan 是一款基于 Go 语言开发的现代化 Web 安全扫描引擎，采用浏览器自动化技术，能够准确检测各类 Web 应用安全漏洞。该工具集成了智能爬虫、动态内容分析和多种漏洞检测插件，为安全研究人员和渗透测试工程师提供强大的自动化扫描能力。

### 核心特性

- 🌐 **浏览器自动化引擎** - 基于 Playwright 的动态页面渲染和交互
- 🔍 **智能漏洞检测** - 支持 SQL 注入、XSS、命令注入等多种漏洞类型
- 🤖 **AI 辅助扫描** - 集成 Stagehand AI 进行智能页面分析
- 📊 **并发扫描架构** - 高效的任务调度和并发控制
- 🔐 **认证管理** - 支持多种认证方式和会话保持
- 📝 **丰富的报告格式** - HTML、JSON 等多种输出格式

## 📁 项目结构

```
DroneRiskScan/
├── cmd/
│   └── dronescan/          # 主程序入口
├── internal/
│   ├── auth/               # 认证管理模块
│   ├── browser/            # 浏览器自动化引擎
│   │   ├── playwright.go   # Playwright 集成
│   │   └── stagehand.go    # Stagehand AI 集成
│   ├── crawler/            # 智能爬虫模块
│   ├── detector/           # 漏洞检测器
│   │   ├── base.go         # 检测器基类
│   │   └── injection/      # 注入类漏洞检测
│   │       ├── sqli.go     # SQL 注入检测
│   │       └── sqli_enhanced.go  # 增强 SQL 注入检测
│   ├── engine/             # 扫描引擎核心
│   │   ├── scanner.go      # 扫描器主逻辑
│   │   └── hybrid.go       # 混合扫描引擎
│   ├── reporter/           # 报告生成器
│   ├── scheduler/          # 任务调度器
│   └── transport/          # HTTP 传输层
├── pkg/
│   └── models/             # 数据模型
│       ├── scan.go         # 扫描任务模型
│       └── vulnerability.go # 漏洞数据模型
├── scripts/
│   └── stagehand_auth.py   # Stagehand 认证脚本
├── reports/                # 扫描报告输出目录
├── docker-compose.yml      # Docker 编排配置
└── test_targets.txt        # 测试目标列表
```

## 🛠️ 安装部署

### 环境要求

- Go 1.18+
- Python 3.8+ (用于 Stagehand AI)
- Docker & Docker Compose (可选)
- Playwright 浏览器驱动

### 快速安装

1. **克隆项目**
```bash
git clone https://github.com/gkfnf/DroneRiskScan.git
cd DroneRiskScan
```

2. **安装依赖**
```bash
# 安装 Go 依赖
go mod download

# 安装 Python 依赖（如使用 Stagehand）
pip install -r requirements.txt

# 安装 Playwright 浏览器
playwright install chromium
```

3. **编译项目**
```bash
go build -o dronescan ./cmd/dronescan
```

### Docker 部署

```bash
# 使用 Docker Compose 启动
docker-compose up -d

# 查看运行状态
docker-compose ps
```

## 📖 使用指南

### 基础扫描

```bash
# 扫描单个目标
./dronescan -target https://example.com

# 从文件批量扫描
./dronescan -targets-file test_targets.txt

# 指定输出目录
./dronescan -target https://example.com -output reports/
```

### 高级选项

```bash
# 启用调试模式
./dronescan -target https://example.com -debug

# 设置并发数
./dronescan -target https://example.com -concurrency 10

# 指定风险等级
./dronescan -target https://example.com -risk-level high

# 启用特定插件
./dronescan -target https://example.com -enable-plugins sqli,xss

# 使用认证扫描
./dronescan -target https://example.com \
    -login-url https://example.com/login \
    -username admin \
    -password secret
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-target` | 扫描目标 URL | - |
| `-targets-file` | 目标列表文件 | - |
| `-output` | 报告输出目录 | `./reports` |
| `-report-format` | 报告格式 (html/json) | `html` |
| `-concurrency` | 并发扫描数 | `5` |
| `-timeout` | 请求超时时间 | `30s` |
| `-risk-level` | 风险等级 (low/medium/high/critical) | `medium` |
| `-debug` | 调试模式 | `false` |
| `-verbose` | 详细输出 | `false` |
| `-enable-plugins` | 启用插件列表 | all |
| `-disable-plugins` | 禁用插件列表 | - |
| `-show-plugins` | 显示可用插件 | - |
| `-version` | 显示版本信息 | - |

## 🔌 插件系统

### 已支持的漏洞类型

- **SQL 注入** - 多种 SQL 注入检测技术
  - 布尔盲注
  - 时间盲注
  - 错误回显
  - 联合查询
- **XSS** - 跨站脚本攻击
  - 反射型 XSS
  - 存储型 XSS
  - DOM XSS
- **命令注入** - OS 命令执行
- **文件包含** - 本地/远程文件包含
- **LDAP 注入** - 目录服务注入
- **XXE** - XML 外部实体注入
- **SSRF** - 服务器端请求伪造
- **路径遍历** - 目录穿越攻击

### 自定义插件开发

创建自定义检测器，实现 `Detector` 接口：

```go
type Detector interface {
    Name() string
    Detect(target string, params map[string]string) (*Vulnerability, error)
    GetRiskLevel() string
}
```

## 🧪 测试环境

项目包含了用于测试的靶场环境配置：

```bash
# 启动 bWAPP 测试环境
docker run -d -p 8081:80 raesene/bwapp

# 运行测试扫描
./dronescan -targets-file test_targets.txt
```

## 📊 扫描报告

扫描完成后会在 `reports/` 目录生成详细报告：

- **HTML 报告** - 可视化展示扫描结果
- **JSON 报告** - 结构化数据，便于集成
- **日志文件** - 详细的扫描过程记录

报告包含：
- 漏洞详情和风险等级
- 复现步骤和 Payload
- 修复建议
- 扫描统计信息

## 🔧 配置说明

### Stagehand AI 配置

如需使用 AI 辅助扫描功能：

```bash
# 运行认证脚本
python scripts/stagehand_auth.py

# 配置 API 密钥
export STAGEHAND_API_KEY="your-api-key"
```

### 代理配置

```bash
# HTTP 代理
export HTTP_PROXY="http://proxy:8080"

# HTTPS 代理  
export HTTPS_PROXY="http://proxy:8080"
```

## 🚀 开发计划

- [ ] 支持更多漏洞类型检测
- [ ] 增强 JavaScript 动态分析能力
- [ ] 添加 API 安全扫描功能
- [ ] 实现分布式扫描架构
- [ ] 集成更多 AI 分析能力
- [ ] 支持自定义扫描策略
- [ ] 添加 Web UI 界面
- [ ] 支持扫描任务管理
- [ ] 增加漏洞验证功能
- [ ] 优化内存使用和性能

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## ⚠️ 免责声明

本工具仅供安全研究和授权测试使用。使用者需遵守当地法律法规，对未授权目标进行扫描可能触犯法律。开发者不对任何非法使用承担责任。

## 📮 联系方式

- GitHub: [https://github.com/gkfnf/DroneRiskScan](https://github.com/gkfnf/DroneRiskScan)
- Issues: [https://github.com/gkfnf/DroneRiskScan/issues](https://github.com/gkfnf/DroneRiskScan/issues)

---

**DroneRiskScan** - 让 Web 安全扫描更智能、更高效！ 🛡️