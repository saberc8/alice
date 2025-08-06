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

package router

import (
	"alice/api/handler"

	"github.com/gin-gonic/gin"
)

// SetupRBACRoutes 设置RBAC路由
func SetupRBACRoutes(
	r *gin.Engine,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	menuHandler *handler.MenuHandler,
) {
	// API版本分组
	v1 := r.Group("/api/v1")

	// 角色管理路由
	roles := v1.Group("/roles")
	{
		roles.POST("", roleHandler.CreateRole)       // 创建角色
		roles.GET("/:id", roleHandler.GetRole)       // 获取单个角色
		roles.GET("", roleHandler.ListRoles)         // 获取角色列表
		roles.PUT("/:id", roleHandler.UpdateRole)    // 更新角色
		roles.DELETE("/:id", roleHandler.DeleteRole) // 删除角色
	}

	// 用户角色管理路由
	userRoles := v1.Group("/users")
	{
		userRoles.GET("/:user_id/roles", roleHandler.GetUserRoles)           // 获取用户角色
		userRoles.POST("/:user_id/roles", roleHandler.AssignRolesToUser)     // 为用户分配角色
		userRoles.DELETE("/:user_id/roles", roleHandler.RemoveRolesFromUser) // 移除用户角色
	}

	// 权限管理路由
	permissions := v1.Group("/permissions")
	{
		permissions.POST("", permissionHandler.CreatePermission)       // 创建权限
		permissions.GET("/:id", permissionHandler.GetPermission)       // 获取单个权限
		permissions.GET("", permissionHandler.ListPermissions)         // 获取权限列表
		permissions.PUT("/:id", permissionHandler.UpdatePermission)    // 更新权限
		permissions.DELETE("/:id", permissionHandler.DeletePermission) // 删除权限
	}

	// 角色权限管理路由
	rolePermissions := v1.Group("/roles")
	{
		rolePermissions.GET("/:id/permissions", permissionHandler.GetRolePermissions)           // 获取角色权限
		rolePermissions.POST("/:id/permissions", permissionHandler.AssignPermissionsToRole)     // 为角色分配权限
		rolePermissions.DELETE("/:id/permissions", permissionHandler.RemovePermissionsFromRole) // 移除角色权限
	}

	// 用户权限查询路由
	userPermissions := v1.Group("/users")
	{
		userPermissions.GET("/:user_id/permissions", permissionHandler.GetUserPermissions)        // 获取用户权限
		userPermissions.GET("/:user_id/permissions/check", permissionHandler.CheckUserPermission) // 检查用户权限
	}

	// 菜单管理路由
	menus := v1.Group("/menus")
	{
		menus.POST("", menuHandler.CreateMenu)       // 创建菜单
		menus.GET("/:id", menuHandler.GetMenu)       // 获取单个菜单
		menus.GET("", menuHandler.ListMenus)         // 获取菜单列表
		menus.GET("/tree", menuHandler.GetMenuTree)  // 获取菜单树
		menus.PUT("/:id", menuHandler.UpdateMenu)    // 更新菜单
		menus.DELETE("/:id", menuHandler.DeleteMenu) // 删除菜单
	}

	// 角色菜单管理路由
	roleMenus := v1.Group("/roles")
	{
		roleMenus.GET("/:id/menus", menuHandler.GetRoleMenus)           // 获取角色菜单
		roleMenus.POST("/:id/menus", menuHandler.AssignMenusToRole)     // 为角色分配菜单
		roleMenus.DELETE("/:id/menus", menuHandler.RemoveMenusFromRole) // 移除角色菜单
	}

	// 用户菜单查询路由
	userMenus := v1.Group("/users")
	{
		userMenus.GET("/:user_id/menus", menuHandler.GetUserMenus)         // 获取用户菜单
		userMenus.GET("/:user_id/menus/tree", menuHandler.GetUserMenuTree) // 获取用户菜单树
	}
}
