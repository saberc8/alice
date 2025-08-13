package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

		// RBAC表
		&rbacEntity.Role{},
		&rbacEntity.Permission{},
		&rbacEntity.Menu{},
		&rbacEntity.UserRole{},
		&rbacEntity.RolePermission{},
		&rbacEntity.RoleMenu{},
	)
}
