package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	userHandler := handler.NewUserHandler(application.UserSvc, application.RoleSvc)
	appUserHandler := handler.NewAppUserHandler(application.AppUserSvc, application.FriendSvc)
	roleHandler := handler.NewRoleHandler(application.RoleSvc)
	permissionHandler := handler.NewPermissionHandler(application.PermissionSvc)
	menuHandler := handler.NewMenuHandler(application.MenuSvc, application.PermissionSvc)

	// 初始化路由
	apiRouter := router.NewRouter(userHandler, appUserHandler, roleHandler, permissionHandler, menuHandler)
	r := apiRouter.SetupRoutes()

	// 启动HTTP服务器
	// 兼容：若 Port 字段已包含冒号或完整地址，则优先使用；否则拼接 Host+":"+Port
	addr := cfg.Server.Port
	if addr == "" {
		addr = ":8090"
	}
	// 如果 addr 形如 "8090"，补冒号；如果不包含冒号且 Host 存在，则拼接 Host
	if cfg.Server.Host != "" {
		// 如果 addr 已经包含冒号，认为是完整的 ":port" 或 "host:port"，否则与 Host 拼接
		if addr[0] == ':' {
			addr = cfg.Server.Host + addr
		} else if !containsColon(addr) {
			addr = cfg.Server.Host + ":" + addr
		}
	}

	// 友好输出：列出可通过 IPv4 访问的地址
	if p := extractPort(addr); p != "" {
		for _, ip := range getIPv4Addrs() {
			logger.Infof("Accessible at: http://%s:%s", ip, p)
		}
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动服务器
	go func() {
		logger.Infof("Server starting on %s", addr)
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

// containsColon 判断字符串中是否包含冒号
func containsColon(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return true
		}
	}
	return false
}

// extractPort 从监听地址中提取端口，兼容 ":8090" 或 "0.0.0.0:8090"
func extractPort(addr string) string {
	if addr == "" {
		return ""
	}
	// 尝试使用标准解析
	if _, p, err := net.SplitHostPort(addr); err == nil {
		return p
	}
	// 回退：若是以冒号开头，例如 ":8090"
	return strings.TrimPrefix(addr, ":")
}

// getIPv4Addrs 返回本机所有非回环 IPv4 地址
func getIPv4Addrs() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		// 忽略未启用或回环接口
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 仅保留 IPv4
			}
			ips = append(ips, ip.String())
		}
	}
	return ips
}
