package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"alice/api/handler"
	"alice/api/middleware"
	"alice/application"
	"alice/infra/config"
)

type Router struct {
	userHandler       *handler.UserHandler
	roleHandler       *handler.RoleHandler
	permissionHandler *handler.PermissionHandler
	menuHandler       *handler.MenuHandler
}

func NewRouter(
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	menuHandler *handler.MenuHandler,
) *Router {
	return &Router{
		userHandler:       userHandler,
		roleHandler:       roleHandler,
		permissionHandler: permissionHandler,
		menuHandler:       menuHandler,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// API路由组 (添加公共与受保护子组)
	v1 := router.Group("/api/v1")

	// 用户认证相关 (仅注册/登录无需 token)
	userAuth := v1.Group("/auth")
	{
		userAuth.POST("/register", r.userHandler.Register)
		userAuth.POST("/login", r.userHandler.Login)

		// 其余用户接口需要认证
		userProtected := userAuth.Group("")
		userProtected.Use(middleware.JWTAuth())
		{
			userProtected.GET("/profile", r.userHandler.GetProfile)
			userProtected.PUT("/profile", r.userHandler.UpdateProfile)
		}
	}

	// 受保护的业务功能路由 (除上面开放的登录注册外全部要求 token)
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth())
	{
		// 设置RBAC路由 (已处于受保护组中)
		SetupRBACRoutes(protected, r.roleHandler, r.permissionHandler, r.menuHandler)

		// 用户管理 CRUD
		users := protected.Group("/users")
		{
			users.POST("", middleware.RequirePerm(application.PermissionSvc, "system:user:create"), r.userHandler.CreateUser)
			users.GET("", middleware.RequirePerm(application.PermissionSvc, "system:user:list"), r.userHandler.ListUsers)
			users.GET("/:user_id", middleware.RequirePerm(application.PermissionSvc, "system:user:get"), r.userHandler.GetUser)
			users.PUT("/:user_id", middleware.RequirePerm(application.PermissionSvc, "system:user:update"), r.userHandler.UpdateUser)
			users.DELETE("/:user_id", middleware.RequirePerm(application.PermissionSvc, "system:user:delete"), r.userHandler.DeleteUser)
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger 文档路由 (根据配置开关)
	cfg := config.Load() // 简单方式(注意: 若频繁调用可考虑依赖注入避免重复解析)
	if cfg.Server.EnableSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
