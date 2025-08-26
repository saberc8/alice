package router

import (
	"alice/api/handler"
	"alice/api/middleware"
	"alice/application"

	"github.com/gin-gonic/gin"
)

// SetupRBACRoutes 设置RBAC路由
func SetupRBACRoutes(
	r *gin.RouterGroup,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	menuHandler *handler.MenuHandler,
) {
	// 直接在传入的受保护组下注册 (调用方确保已加需要的中间件)
	v1 := r // 保持原命名语义

	// 角色管理路由
	roles := v1.Group("/roles")
	{
		roles.POST("", middleware.RequirePerm(application.PermissionSvc, "system:role:create"), roleHandler.CreateRole)
		roles.GET("/:id", middleware.RequirePerm(application.PermissionSvc, "system:role:get"), roleHandler.GetRole)
		roles.GET("", middleware.RequirePerm(application.PermissionSvc, "system:role:list"), roleHandler.ListRoles)
		roles.PUT("/:id", middleware.RequirePerm(application.PermissionSvc, "system:role:update"), roleHandler.UpdateRole)
		roles.DELETE("/:id", middleware.RequirePerm(application.PermissionSvc, "system:role:delete"), roleHandler.DeleteRole)
	}

	// 用户角色管理路由
	userRoles := v1.Group("/users")
	{
		userRoles.GET("/:user_id/roles", middleware.RequirePerm(application.PermissionSvc, "system:user:roles:get"), roleHandler.GetUserRoles)
		userRoles.POST("/:user_id/roles", middleware.RequirePerm(application.PermissionSvc, "system:user:roles:assign"), roleHandler.AssignRolesToUser)
		userRoles.DELETE("/:user_id/roles", middleware.RequirePerm(application.PermissionSvc, "system:user:roles:remove"), roleHandler.RemoveRolesFromUser)
	}

	// 权限管理路由
	permissions := v1.Group("/permissions")
	{
		permissions.POST("", middleware.RequirePerm(application.PermissionSvc, "system:permission:create"), permissionHandler.CreatePermission)
		permissions.GET("/:id", middleware.RequirePerm(application.PermissionSvc, "system:permission:get"), permissionHandler.GetPermission)
		permissions.GET("", middleware.RequirePerm(application.PermissionSvc, "system:permission:list"), permissionHandler.ListPermissions)
		permissions.PUT("/:id", middleware.RequirePerm(application.PermissionSvc, "system:permission:update"), permissionHandler.UpdatePermission)
		permissions.DELETE("/:id", middleware.RequirePerm(application.PermissionSvc, "system:permission:delete"), permissionHandler.DeletePermission)
	}

	// 角色权限管理路由
	rolePermissions := v1.Group("/roles")
	{
		rolePermissions.GET("/:id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:role:permissions:get"), permissionHandler.GetRolePermissions)
		rolePermissions.POST("/:id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:role:permissions:assign"), permissionHandler.AssignPermissionsToRole)
		rolePermissions.DELETE("/:id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:role:permissions:remove"), permissionHandler.RemovePermissionsFromRole)
	}

	// 用户权限查询路由
	userPermissions := v1.Group("/users")
	{
		userPermissions.GET("/:user_id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:user:permissions:get"), permissionHandler.GetUserPermissions)
		userPermissions.GET("/:user_id/permissions/check", middleware.RequirePerm(application.PermissionSvc, "system:user:permissions:check"), permissionHandler.CheckUserPermission)
	}

	// 菜单管理路由
	menus := v1.Group("/menus")
	{
		menus.POST("", middleware.RequirePerm(application.PermissionSvc, "system:menu:create"), menuHandler.CreateMenu)
		menus.GET("/:id", middleware.RequirePerm(application.PermissionSvc, "system:menu:get"), menuHandler.GetMenu)
		menus.GET("", middleware.RequirePerm(application.PermissionSvc, "system:menu:list"), menuHandler.ListMenus)
		menus.GET("/tree", middleware.RequirePerm(application.PermissionSvc, "system:menu:list"), menuHandler.GetMenuTree)
		menus.PUT("/:id", middleware.RequirePerm(application.PermissionSvc, "system:menu:update"), menuHandler.UpdateMenu)
		menus.DELETE("/:id", middleware.RequirePerm(application.PermissionSvc, "system:menu:delete"), menuHandler.DeleteMenu)
		menus.GET("/:id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:menu:list"), menuHandler.ListMenuPermissions)
		menus.POST("/:id/permissions", middleware.RequirePerm(application.PermissionSvc, "system:permission:create"), menuHandler.CreateMenuPermission)
	}

	// 角色菜单管理路由
	roleMenus := v1.Group("/roles")
	{
		roleMenus.GET("/:id/menus", middleware.RequirePerm(application.PermissionSvc, "system:role:menus:list"), menuHandler.GetRoleMenus)
		roleMenus.POST("/:id/menus", middleware.RequirePerm(application.PermissionSvc, "system:role:menus:assign"), menuHandler.AssignMenusToRole)
		roleMenus.DELETE("/:id/menus", middleware.RequirePerm(application.PermissionSvc, "system:role:menus:remove"), menuHandler.RemoveMenusFromRole)
		roleMenus.GET("/:id/menus/tree", middleware.RequirePerm(application.PermissionSvc, "system:role:menus:list"), menuHandler.GetRoleMenuTree)
	}

	// 用户菜单查询路由
	userMenus := v1.Group("/users")
	{
		userMenus.GET("/:user_id/menus", menuHandler.GetUserMenus)         // 获取用户菜单
		userMenus.GET("/:user_id/menus/tree", menuHandler.GetUserMenuTree) // 获取用户菜单树
	}
}
