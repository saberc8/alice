package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	friendEntity "alice/domain/appfriend/entity"
	appEntity "alice/domain/appuser/entity"
	chatEntity "alice/domain/chat/entity"
	momentEntity "alice/domain/moment/entity"
	rbacEntity "alice/domain/rbac/entity"
	"alice/domain/user/entity"
	"alice/infra/config"
	"alice/pkg/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Database connected successfully")
	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 用户表
		&entity.User{},

		// App 端表
		&appEntity.AppUser{},
		&friendEntity.FriendRelation{},
		&friendEntity.FriendRequest{},

		// Moments
		&momentEntity.Moment{},
		&momentEntity.MomentLike{},
		&momentEntity.MomentComment{},

		// Chat
		&chatEntity.Message{},
		&chatEntity.Group{},
		&chatEntity.GroupMember{},
		&chatEntity.GroupMessage{},
		&chatEntity.GroupReadCursor{},

		// RBAC表
		&rbacEntity.Role{},
		&rbacEntity.Permission{},
		&rbacEntity.Menu{},
		&rbacEntity.UserRole{},
		&rbacEntity.RolePermission{},
		&rbacEntity.RoleMenu{},
	)
}
