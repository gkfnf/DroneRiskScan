# 电力无人机安全扫描系统

专业的无人机设备安全漏洞检测与评估平台，支持网络基础设施、射频安全、设备管理等功能。

## 🚀 快速开始

### 环境要求

- **Node.js**: >= 18.0.0
- **Go**: >= 1.19 (后端服务)
- **Python**: >= 3.8 (可选，用于脚本)
- **操作系统**: Windows 10+, macOS 10.15+, Ubuntu 18.04+

### 安装依赖

\`\`\`bash
# 克隆项目
git clone https://github.com/your-org/drone-security-scanner.git
cd drone-security-scanner

# 安装 Node.js 依赖
npm install

# 安装 Electron 依赖
npm run postinstall
\`\`\`

## 🖥️ 开发模式运行

### 1. Web 版本开发

\`\`\`bash
# 启动 Next.js 开发服务器
npm run dev

# 浏览器访问
open http://localhost:3000
\`\`\`

### 2. Electron 桌面应用开发

\`\`\`bash
# 同时启动 Web 服务和 Electron 应用
npm run electron:dev

# 或者分别启动
npm run dev          # 终端1: 启动 Web 服务
npm run electron     # 终端2: 启动 Electron 应用
\`\`\`

### 3. Go 后端服务

\`\`\`bash
# 进入后端目录
cd backend

# 安装 Go 依赖
go mod tidy

# 启动后端服务
go run main.go

# 或者编译后运行
go build -o main
./main
\`\`\`

## 📦 生产环境部署

### 1. 构建 Web 应用

\`\`\`bash
# 构建生产版本
npm run build

# 启动生产服务器
npm start
\`\`\`

### 2. 构建 Electron 应用

\`\`\`bash
# 构建所有平台
npm run electron:build

# 仅构建当前平台
npm run electron:dist

# 构建特定平台
npx electron-builder --win    # Windows
npx electron-builder --mac    # macOS
npx electron-builder --linux  # Linux
\`\`\`

### 3. Docker 部署

\`\`\`bash
# 构建 Docker 镜像
docker build -t drone-security-scanner .

# 运行容器
docker run -p 3000:3000 -p 8080:8080 drone-security-scanner
\`\`\`

## 🔧 配置说明

### 环境变量

创建 `.env.local` 文件：

\`\`\`env
# 数据库配置
DATABASE_URL=sqlite:./data/scanner.db

# Go 后端服务
BACKEND_URL=http://localhost:8080

# 高德地图 API Key
NEXT_PUBLIC_AMAP_KEY=your_amap_api_key

# Nuclei 配置
NUCLEI_TEMPLATES_PATH=./nuclei-templates

# 射频设备配置
RF_DEVICE_PORT=/dev/ttyUSB0
RF_SCAN_FREQUENCY_RANGE=2400-5800

# 安全配置
JWT_SECRET=your_jwt_secret
ENCRYPTION_KEY=your_encryption_key
\`\`\`

### Go 后端配置

创建 `backend/config.yaml`：

\`\`\`yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  type: "sqlite"
  path: "./data/scanner.db"

nuclei:
  templates_path: "./nuclei-templates"
  concurrency: 10
  rate_limit: 150

grpc:
  port: 9090
  tls_enabled: false

logging:
  level: "info"
  file: "./logs/scanner.log"
\`\`\`

## 🛠️ 开发工具设置

### VS Code 配置

安装推荐扩展：
- TypeScript and JavaScript Language Features
- Tailwind CSS IntelliSense
- Go
- Electron Debug
- GitLens

### 调试配置

创建 `.vscode/launch.json`：

\`\`\`json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Electron Main",
      "type": "node",
      "request": "launch",
      "cwd": "${workspaceFolder}",
      "program": "${workspaceFolder}/node_modules/.bin/electron",
      "args": [".", "--remote-debugging-port=9222"],
      "outputCapture": "std"
    },
    {
      "name": "Debug Go Backend",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/main.go",
      "cwd": "${workspaceFolder}/backend"
    }
  ]
}
\`\`\`

## 📱 移动端开发

### React Native (可选)

\`\`\`bash
# 安装 React Native CLI
npm install -g @react-native-community/cli

# 创建移动端项目
npx react-native init DroneSecurityMobile

# 运行 iOS
npx react-native run-ios

# 运行 Android
npx react-native run-android
\`\`\`

## 🧪 测试

### 单元测试

\`\`\`bash
# 运行所有测试
npm test

# 运行特定测试
npm test -- --testNamePattern="射频安全"

# 生成覆盖率报告
npm run test:coverage
\`\`\`

### E2E 测试

\`\`\`bash
# 安装 Playwright
npm install -D @playwright/test

# 运行 E2E 测试
npm run test:e2e
\`\`\`

### Go 后端测试

\`\`\`bash
cd backend

# 运行单元测试
go test ./...

# 运行基准测试
go test -bench=. ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
\`\`\`

## 🔍 故障排除

### 常见问题

1. **Electron 启动失败**
   \`\`\`bash
   # 清理 node_modules 重新安装
   rm -rf node_modules package-lock.json
   npm install
   \`\`\`

2. **Go 后端连接失败**
   \`\`\`bash
   # 检查端口占用
   lsof -i :8080
   
   # 重启后端服务
   pkill -f "go run main.go"
   go run main.go
   \`\`\`

3. **地图加载失败**
   - 检查高德地图 API Key 是否正确
   - 确认网络连接正常
   - 检查域名白名单设置

4. **射频设备无法访问**
   \`\`\`bash
   # Linux 下添加用户到 dialout 组
   sudo usermod -a -G dialout $USER
   
   # 重新登录或重启
   \`\`\`

### 日志查看

\`\`\`bash
# Electron 应用日志
# Windows: %APPDATA%/drone-security-scanner/logs
# macOS: ~/Library/Logs/drone-security-scanner
# Linux: ~/.config/drone-security-scanner/logs

# Go 后端日志
tail -f backend/logs/scanner.log

# 系统日志
journalctl -u drone-security-scanner
\`\`\`

## 📊 性能优化

### 前端优化

\`\`\`bash
# 分析打包大小
npm run analyze

# 优化图片资源
npm install -D next-optimized-images

# 启用 PWA
npm install next-pwa
\`\`\`

### 后端优化

\`\`\`bash
# Go 性能分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 数据库优化
sqlite3 data/scanner.db "VACUUM; ANALYZE;"
\`\`\`

## 🚀 部署选项

### 1. 本地部署

\`\`\`bash
# 使用 PM2 管理进程
npm install -g pm2

# 启动服务
pm2 start ecosystem.config.js

# 查看状态
pm2 status
\`\`\`

### 2. 云服务部署

#### Vercel 部署
\`\`\`bash
npm install -g vercel
vercel --prod
\`\`\`

#### Docker Compose 部署
\`\`\`bash
docker-compose up -d
\`\`\`

### 3. 企业级部署

#### Kubernetes 部署
\`\`\`bash
kubectl apply -f k8s/
\`\`\`

#### 负载均衡配置
```nginx
upstream drone_scanner {
    server 127.0.0.1:3000;
    server 127.0.0.1:3001;
}

server {
    listen 80;
    server_name scanner.company.com;
    
    location / {
        proxy_pass http://drone_scanner;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
