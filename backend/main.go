package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"drone-security-scanner/internal/api"
	"drone-security-scanner/internal/config"
	"drone-security-scanner/internal/database"
	"drone-security-scanner/internal/grpc/server"
	"drone-security-scanner/internal/nuclei"
	"drone-security-scanner/internal/scanner"
	"drone-security-scanner/internal/services"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.Initialize(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 初始化服务
	assetService := services.NewAssetService(db)
	scanService := services.NewScanService(db)
	vulnerabilityService := services.NewVulnerabilityService(db)
	rfService := services.NewRFService(db)

	// 初始化扫描器
	nucleiScanner := nuclei.NewScanner(cfg.Nuclei.TemplatesPath)
	scannerManager := scanner.NewManager(nucleiScanner, scanService, vulnerabilityService)

	// 启动 gRPC 服务器
	go startGRPCServer(cfg, assetService, scanService, vulnerabilityService, rfService)

	// 启动 HTTP API 服务器
	startHTTPServer(cfg, assetService, scanService, vulnerabilityService, rfService, scannerManager)
}

func startGRPCServer(cfg *config.Config, assetService *services.AssetService, scanService *services.ScanService, vulnerabilityService *services.VulnerabilityService, rfService *services.RFService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	s := grpc.NewServer()
	
	// 注册服务
	server.RegisterAssetService(s, assetService)
	server.RegisterScanService(s, scanService)
	server.RegisterVulnerabilityService(s, vulnerabilityService)
	server.RegisterRFService(s, rfService)

	// 启用反射（开发环境）
	reflection.Register(s)

	log.Printf("gRPC server listening on port %d", cfg.GRPC.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(cfg *config.Config, assetService *services.AssetService, scanService *services.ScanService, vulnerabilityService *services.VulnerabilityService, rfService *services.RFService, scannerManager *scanner.Manager) {
	// 设置 Gin 模式
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	// API 路由
	apiGroup := r.Group("/api/v1")
	api.SetupRoutes(apiGroup, assetService, scanService, vulnerabilityService, rfService, scannerManager)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		log.Printf("HTTP server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
