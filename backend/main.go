package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alice/api/handler"
	"alice/api/router"
	"alice/application"
	"alice/infra/config"
	"alice/pkg/logger"

	_ "alice/docs" // swagger 文档 (由 swag 工具生成)
)

// @title Alice API
// @version 1.0
// @description Alice 企业级后端 API 文档。
// @BasePath /api/v1
// @schemes http
// @contact.name API Support
// @contact.url https://github.com/coze-dev/alice
// @contact.email support@example.com
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入形如: Bearer <token>

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化应用依赖
	ctx := context.Background()
	if err := application.Init(ctx, cfg); err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// 初始化处理器
	userHandler := handler.NewUserHandler(application.UserSvc)
	roleHandler := handler.NewRoleHandler(application.RoleSvc)
	permissionHandler := handler.NewPermissionHandler(application.PermissionSvc)
	menuHandler := handler.NewMenuHandler(application.MenuSvc)

	// 初始化路由
	apiRouter := router.NewRouter(userHandler, roleHandler, permissionHandler, menuHandler)
	r := apiRouter.SetupRoutes()

	// 启动HTTP服务器
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	// 启动服务器
	go func() {
		logger.Infof("Server starting on %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exited")
}
