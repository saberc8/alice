/*
 * Copyright 2025 alice Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
)

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
