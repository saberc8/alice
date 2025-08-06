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

package application

import (
	"context"

	rbacService "alice/domain/rbac/service"
	"alice/domain/user/service"
	"alice/infra/config"
	"alice/infra/database"
	"alice/infra/repository"
	"alice/pkg/logger"
)

var (
	// UserSvc 用户服务实例
	UserSvc service.UserService

	// RBAC 服务实例
	RoleSvc       rbacService.RoleService
	PermissionSvc rbacService.PermissionService
	MenuSvc       rbacService.MenuService
)

// Init 初始化应用
func Init(ctx context.Context, cfg *config.Config) error {
	// 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		return err
	}

	// 初始化用户仓储
	userRepo := repository.NewUserRepository(db)

	// 初始化RBAC仓储
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	menuRepo := repository.NewMenuRepository(db)

	// 初始化用户服务
	UserSvc = service.NewUserService(userRepo)

	// 初始化RBAC服务
	RoleSvc = rbacService.NewRoleService(roleRepo)
	PermissionSvc = rbacService.NewPermissionService(permissionRepo)
	MenuSvc = rbacService.NewMenuService(menuRepo)

	logger.Info("Application initialized successfully")
	return nil
}
