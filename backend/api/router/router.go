package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"alice/api/handler"
	chathdl "alice/api/handler/chat"
	"alice/api/middleware"
	"alice/application"
	"alice/infra/config"
)

type Router struct {
	userHandler       *handler.UserHandler
	appUserHandler    *handler.AppUserHandler
	roleHandler       *handler.RoleHandler
	permissionHandler *handler.PermissionHandler
	menuHandler       *handler.MenuHandler
	chatHub           *chathdl.Hub
	storageHandler    *handler.StorageHandler
	momentHandler     *handler.MomentHandler
}

func NewRouter(
	userHandler *handler.UserHandler,
	appUserHandler *handler.AppUserHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	menuHandler *handler.MenuHandler,
) *Router {
	// 初始化聊天 Hub（基于应用层 ChatSvc）
	hub := chathdl.NewHub(application.ChatSvc)
	storageHandler := handler.NewStorageHandler()
	momentHandler := handler.NewMomentHandler(application.MomentSvc)
	return &Router{
		userHandler:       userHandler,
		appUserHandler:    appUserHandler,
		roleHandler:       roleHandler,
		permissionHandler: permissionHandler,
		menuHandler:       menuHandler,
		chatHub:           hub,
		storageHandler:    storageHandler,
		momentHandler:     momentHandler,
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

		// 存储相关路由
		storage := protected.Group("/storage")
		{
			buckets := storage.Group("/buckets")
			{
				buckets.GET("", r.storageHandler.ListBuckets)
				buckets.POST(":bucket", r.storageHandler.CreateBucket)
				buckets.DELETE(":bucket", r.storageHandler.DeleteBucket)
				buckets.GET(":bucket/objects", r.storageHandler.ListObjects)
				buckets.POST(":bucket/objects", r.storageHandler.UploadObject)
				buckets.DELETE(":bucket/objects/:object", r.storageHandler.DeleteObject)
				buckets.GET(":bucket/objects/:object/url", r.storageHandler.GetObjectPresigned)
				buckets.POST(":bucket/public", r.storageHandler.SetBucketPublic)
			}
		}

	}

	// ===== 移动端 App 路由（独立于后台鉴权） =====
	app := v1.Group("/app")
	{
		app.POST("/register", r.appUserHandler.AppRegister)
		app.POST("/login", r.appUserHandler.AppLogin)

		appProtected := app.Group("")
		appProtected.Use(middleware.AppJWTAuth())
		{
			appProtected.GET("/profile", r.appUserHandler.AppProfile)
			appProtected.PUT("/profile", r.appUserHandler.AppUpdateProfile)
			appProtected.POST("/profile/avatar", r.appUserHandler.AppUploadAvatar)
			appProtected.POST("/friends/request", r.appUserHandler.RequestFriend)
			appProtected.GET("/friends", r.appUserHandler.ListFriends)
			appProtected.GET("/friends/requests", r.appUserHandler.ListPendingRequests)
			appProtected.POST("/friends/requests/:request_id/accept", r.appUserHandler.AcceptFriendRequest)
			appProtected.POST("/friends/requests/:request_id/decline", r.appUserHandler.DeclineFriendRequest)

			// Moments
			appProtected.POST("/moments", r.momentHandler.PostMoment)
			appProtected.GET("/moments", r.momentHandler.ListMoments)
			appProtected.DELETE("/moments/:moment_id", r.momentHandler.DeleteMoment)
			appProtected.POST("/moments/images", r.momentHandler.UploadImage)
			appProtected.GET("/users/:user_id/moments", r.momentHandler.ListUserMoments)

			// Chat routes
			chat := appProtected.Group("/chat")
			{
				chat.GET("/ws", r.chatHub.WS)
				chat.GET("/history/:peer_id", r.chatHub.History)
				chat.POST("/read", r.chatHub.MarkRead)
				chat.GET("/conversations", r.chatHub.Conversations)
			}
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
