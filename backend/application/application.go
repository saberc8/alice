package application

import (
	"context"

	appfriendservice "alice/domain/appfriend/service"
	appuserservice "alice/domain/appuser/service"
	chatservice "alice/domain/chat/service"
	momentservice "alice/domain/moment/service"
	rbacService "alice/domain/rbac/service"
	"alice/domain/user/service"
	"alice/infra/config"
	"alice/infra/database"
	"alice/infra/repository"
	chatrepo "alice/infra/repository/chat"
	"alice/infra/storage"
	"alice/pkg/logger"
)

var (
	// UserSvc 用户服务实例
	UserSvc service.UserService

	// App 端服务实例
	AppUserSvc appuserservice.AppUserService
	FriendSvc  appfriendservice.FriendService
	ChatSvc    chatservice.ChatService
	MomentSvc  momentservice.MomentService

	// RBAC 服务实例
	RoleSvc       rbacService.RoleService
	PermissionSvc rbacService.PermissionService
	MenuSvc       rbacService.MenuService

	// 对象存储
	ObjectStore storage.ObjectStorage
)

// Init 初始化应用
func Init(ctx context.Context, cfg *config.Config) error {
	// 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		return err
	}

	// 初始化仓储
	userRepo := repository.NewUserRepository(db)
	appUserRepo := repository.NewAppUserRepository(db)
	friendRepo := repository.NewFriendRepository(db)
	momentRepo := repository.NewMomentRepository(db)
	msgRepo := chatrepo.NewMessageRepository(db)

	// 初始化RBAC仓储
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	menuRepo := repository.NewMenuRepository(db)

	// 初始化服务
	UserSvc = service.NewUserService(userRepo)
	AppUserSvc = appuserservice.NewAppUserService(appUserRepo)
	FriendSvc = appfriendservice.NewFriendService(appUserRepo, friendRepo)
	ChatSvc = chatservice.NewChatService(msgRepo, friendRepo)
	MomentSvc = momentservice.NewMomentService(momentRepo)

	// 初始化RBAC服务
	RoleSvc = rbacService.NewRoleService(roleRepo)
	PermissionSvc = rbacService.NewPermissionService(permissionRepo)
	MenuSvc = rbacService.NewMenuService(menuRepo, permissionRepo)

	// 初始化对象存储（MinIO）
	if cfg.Minio.Endpoint != "" {
		minioCli, err := storage.NewMinio(cfg.Minio)
		if err != nil {
			logger.Errorf("init minio failed: %v", err)
		} else {
			ObjectStore = minioCli
			_ = minioCli.HealthCheck(ctx)
			logger.Info("MinIO storage initialized")
		}
	}

	logger.Info("Application initialized successfully")
	return nil
}
