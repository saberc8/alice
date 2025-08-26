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

	// --- Pre-migration compatibility fix ---
	// 之前版本 permissions.menu_id 可能存储为 uuid/text，现在模型改为 *uint (bigint)。
	// GORM 在 AutoMigrate 时会执行: ALTER TABLE "permissions" ALTER COLUMN "menu_id" TYPE bigint USING "menu_id"::bigint
	// 若列里存在 uuid 字符串将导致 invalid input syntax for type bigint 报错。
	// 处理策略（开发/测试环境安全）：检测当前列数据类型，若不是 bigint，则先将无法转换的值置空，避免迁移失败。
	// 若需要保留旧的 uuid -> 新菜单 ID 映射，请在此逻辑前实现映射策略。
	type columnInfo struct{ DataType string }
	var ci columnInfo
	if err := db.Raw(`SELECT data_type FROM information_schema.columns WHERE table_name='permissions' AND column_name='menu_id'`).Scan(&ci).Error; err == nil {
		if ci.DataType != "bigint" && ci.DataType != "integer" && ci.DataType != "bigserial" && ci.DataType != "" { // 非期望类型且列存在
			// 将所有非纯数字的 menu_id 置为 NULL，保证后续 ::bigint 转换成功
			// 使用正则筛选仅数字的行保留，其余置空（Postgres ~* 为不区分大小写）
			if err := db.Exec("UPDATE permissions SET menu_id = NULL WHERE menu_id IS NOT NULL AND menu_id !~ '^[0-9]+$'").Error; err != nil {
				logger.Warn("预处理 permissions.menu_id 失败", "error", err)
			} else {
				logger.Info("已清理无法转换为 bigint 的 permissions.menu_id 旧数据")
			}
		}
	}

	// 清理 menus.parent_id
	ci = columnInfo{}
	if err := db.Raw(`SELECT data_type FROM information_schema.columns WHERE table_name='menus' AND column_name='parent_id'`).Scan(&ci).Error; err == nil {
		if ci.DataType != "bigint" && ci.DataType != "integer" && ci.DataType != "bigserial" && ci.DataType != "" {
			if err := db.Exec("UPDATE menus SET parent_id = NULL WHERE parent_id IS NOT NULL AND parent_id !~ '^[0-9]+$'").Error; err != nil {
				logger.Warn("预处理 menus.parent_id 失败", "error", err)
			} else {
				logger.Info("已清理无法转换为 bigint 的 menus.parent_id 旧数据")
			}
		}
	}

	// 若检测到旧的 UUID/文本主键结构，直接删表重建（开发/测试环境策略）
	ci = columnInfo{}
	if err := db.Raw(`SELECT data_type FROM information_schema.columns WHERE table_name='menus' AND column_name='id'`).Scan(&ci).Error; err == nil {
		if ci.DataType == "uuid" || ci.DataType == "text" || ci.DataType == "character varying" {
			logger.Warn("检测到旧的 menus.id 类型为", "data_type", ci.DataType, "action", "drop & recreate rbac tables")
			// 依赖关系：关联表先删
			_ = db.Migrator().DropTable("role_menus")
			_ = db.Migrator().DropTable("role_permissions")
			_ = db.Migrator().DropTable("user_roles")
			_ = db.Migrator().DropTable("permissions")
			_ = db.Migrator().DropTable("menus")
			_ = db.Migrator().DropTable("roles")
			logger.Info("旧 RBAC 表已删除, 将由 AutoMigrate 重建")
		}
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
