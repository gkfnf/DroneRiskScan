# 多阶段构建
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY . .

# 构建前端应用
RUN npm run build

# Go 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制 Go 模块文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY backend/ .

# 构建 Go 应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# 从构建阶段复制文件
COPY --from=frontend-builder /app/.next ./.next
COPY --from=frontend-builder /app/public ./public
COPY --from=frontend-builder /app/package.json ./
COPY --from=frontend-builder /app/node_modules ./node_modules
COPY --from=backend-builder /app/main ./backend/

# 创建数据目录
RUN mkdir -p /root/data /root/logs

# 暴露端口
EXPOSE 3000 8080

# 创建启动脚本
RUN echo '#!/bin/sh' > start.sh && \
    echo './backend/main &' >> start.sh && \
    echo 'npm start' >> start.sh && \
    chmod +x start.sh

# 启动应用
CMD ["./start.sh"]
