package main

import (
	"context"
	"fmt"
	"log"

	"alice/domain/rbac/entity"
	rbacService "alice/domain/rbac/service"
	userEntity "alice/domain/user/entity"
	userServicePkg "alice/domain/user/service"
	"alice/infra/config"
	"alice/infra/database"
	"alice/infra/repository"
	"alice/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化仓储
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 初始化服务
	roleService := rbacService.NewRoleService(roleRepo)
	permissionService := rbacService.NewPermissionService(permissionRepo)
	menuService := rbacService.NewMenuService(menuRepo, permissionRepo)
	userService := userServicePkg.NewUserService(userRepo)

	ctx := context.Background()

	// 清理现有数据
	if err := cleanExistingData(db); err != nil {
		log.Fatal("Failed to clean existing data:", err)
	}

	// 初始化数据
	if err := initRoles(ctx, roleService); err != nil {
		log.Fatal("Failed to init roles:", err)
	}

	// 先初始化菜单（以便为权限绑定 menu_id）
	if err := initMenus(ctx, db, menuService); err != nil {
		log.Fatal("Failed to init menus:", err)
	}

	// 再初始化权限（三段式并绑定到对应菜单）
	if err := initPermissions(ctx, db, permissionService); err != nil {
		log.Fatal("Failed to init permissions:", err)
	}

	// 创建 admin 超级管理员并分配所有角色/权限/菜单
	if err := initAdminUser(ctx, db, userService, roleService, permissionService, menuService); err != nil {
		log.Fatal("Failed to init admin user:", err)
	}

	fmt.Println("初始化数据完成!")
}

// initRoles 初始化角色
func initRoles(ctx context.Context, roleService rbacService.RoleService) error {
	roles := []rbacService.CreateRoleRequest{
		{Name: "超级管理员", Code: "super_admin", Status: entity.RoleStatusActive},
		{Name: "管理员", Code: "admin", Status: entity.RoleStatusActive},
		{Name: "普通用户", Code: "user", Status: entity.RoleStatusActive},
	}

	for _, req := range roles {
		_, err := roleService.CreateRole(ctx, &req)
		if err != nil {
			fmt.Printf("创建角色 %s 失败: %v\n", req.Name, err)
		} else {
			fmt.Printf("创建角色 %s 成功\n", req.Name)
		}
	}

	return nil
}

// initAdminUser 创建一个 admin 超级管理员账号并授予全部角色/权限/菜单
func initAdminUser(
	ctx context.Context,
	db *gorm.DB,
	userService userServicePkg.UserService,
	roleService rbacService.RoleService,
	permissionService rbacService.PermissionService,
	menuService rbacService.MenuService,
) error {
	const (
		adminUsername = "admin"
		adminPassword = "123456"
		adminEmail    = "admin@example.com"
	)

	// 如果已存在则跳过
	var existing userEntity.User
	if err := db.Where("username = ?", adminUsername).First(&existing).Error; err == nil {
		fmt.Println("admin 用户已存在，跳过创建")
		return nil
	}

	// 手动创建（避免 Register 逻辑修改 email 唯一校验外的流程变化）
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}
	user := &userEntity.User{Username: adminUsername, Email: adminEmail, PasswordHash: string(hash), Status: userEntity.UserStatusActive}
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("创建 admin 用户失败: %w", err)
	}
	fmt.Println("创建 admin 用户成功, ID=", user.ID)

	// 读取 super_admin 角色 ID
	superRole, err := repository.NewRoleRepository(db).GetByCode(ctx, "super_admin")
	if err != nil || superRole == nil {
		return fmt.Errorf("获取 super_admin 角色失败: %v", err)
	}

	// 给 admin 分配 super_admin 角色
	if err := repository.NewRoleRepository(db).AssignToUser(ctx, fmt.Sprintf("%d", user.ID), []string{superRole.ID}); err != nil {
		return fmt.Errorf("为 admin 分配 super_admin 角色失败: %w", err)
	}
	fmt.Println("已为 admin 分配 super_admin 角色")

	// 获取全部权限并分配给 super_admin 角色
	perms, _, err := repository.NewPermissionRepository(db).List(ctx, 0, 1000)
	if err != nil {
		return fmt.Errorf("获取权限列表失败: %w", err)
	}
	var permIDs []string
	for _, p := range perms {
		permIDs = append(permIDs, p.ID)
	}
	if err := repository.NewPermissionRepository(db).AssignToRole(ctx, superRole.ID, permIDs); err != nil {
		return fmt.Errorf("为 super_admin 分配全部权限失败: %w", err)
	}
	fmt.Println("已为 super_admin 角色分配全部权限 (", len(permIDs), ")")

	// 获取全部菜单并分配给 super_admin 角色
	menus, err := repository.NewMenuRepository(db).List(ctx)
	if err != nil {
		return fmt.Errorf("获取菜单列表失败: %w", err)
	}
	var menuIDs []string
	for _, m := range menus {
		menuIDs = append(menuIDs, m.ID)
	}
	if err := repository.NewMenuRepository(db).AssignToRole(ctx, superRole.ID, menuIDs); err != nil {
		return fmt.Errorf("为 super_admin 分配全部菜单失败: %w", err)
	}
	fmt.Println("已为 super_admin 角色分配全部菜单 (", len(menuIDs), ")")

	return nil
}

// initPermissions 初始化权限
func initPermissions(ctx context.Context, db *gorm.DB, permissionService rbacService.PermissionService) error {
	// 通过菜单 code 查 ID
	getMenuID := func(code string) *string {
		m, err := repository.NewMenuRepository(db).GetByCode(ctx, code)
		if err != nil || m == nil {
			return nil
		}
		return &m.ID
	}

	// 三段式权限码，绑定到对应菜单
	permissions := []rbacService.CreatePermissionRequest{
		// 角色管理 (system:roles)
		{Name: "角色-查询", Code: "system:role:list", MenuID: getMenuID("system:roles"), Resource: "role", Action: "list", Status: entity.PermissionStatusActive},
		{Name: "角色-详情", Code: "system:role:get", MenuID: getMenuID("system:roles"), Resource: "role", Action: "get", Status: entity.PermissionStatusActive},
		{Name: "角色-创建", Code: "system:role:create", MenuID: getMenuID("system:roles"), Resource: "role", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "角色-更新", Code: "system:role:update", MenuID: getMenuID("system:roles"), Resource: "role", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "角色-删除", Code: "system:role:delete", MenuID: getMenuID("system:roles"), Resource: "role", Action: "delete", Status: entity.PermissionStatusActive},
		{Name: "角色-分配菜单", Code: "system:role:menus:assign", MenuID: getMenuID("system:roles"), Resource: "role_menus", Action: "assign", Status: entity.PermissionStatusActive},
		{Name: "角色-移除菜单", Code: "system:role:menus:remove", MenuID: getMenuID("system:roles"), Resource: "role_menus", Action: "remove", Status: entity.PermissionStatusActive},
		{Name: "角色-菜单查询", Code: "system:role:menus:list", MenuID: getMenuID("system:roles"), Resource: "role_menus", Action: "list", Status: entity.PermissionStatusActive},
		{Name: "角色-分配权限", Code: "system:role:permissions:assign", MenuID: getMenuID("system:roles"), Resource: "role_permissions", Action: "assign", Status: entity.PermissionStatusActive},
		{Name: "角色-移除权限", Code: "system:role:permissions:remove", MenuID: getMenuID("system:roles"), Resource: "role_permissions", Action: "remove", Status: entity.PermissionStatusActive},
		{Name: "角色-权限查询", Code: "system:role:permissions:get", MenuID: getMenuID("system:roles"), Resource: "role_permissions", Action: "get", Status: entity.PermissionStatusActive},

		// 用户管理 (system:users)
		{Name: "用户-角色查询", Code: "system:user:roles:get", MenuID: getMenuID("system:users"), Resource: "user_roles", Action: "get", Status: entity.PermissionStatusActive},
		{Name: "用户-分配角色", Code: "system:user:roles:assign", MenuID: getMenuID("system:users"), Resource: "user_roles", Action: "assign", Status: entity.PermissionStatusActive},
		{Name: "用户-移除角色", Code: "system:user:roles:remove", MenuID: getMenuID("system:users"), Resource: "user_roles", Action: "remove", Status: entity.PermissionStatusActive},
		{Name: "用户-权限查询", Code: "system:user:permissions:get", MenuID: getMenuID("system:users"), Resource: "user_permissions", Action: "get", Status: entity.PermissionStatusActive},
		{Name: "用户-权限校验", Code: "system:user:permissions:check", MenuID: getMenuID("system:users"), Resource: "user_permissions", Action: "check", Status: entity.PermissionStatusActive},

		// 权限管理 (system:permissions)
		{Name: "权限-查询", Code: "system:permission:list", MenuID: getMenuID("system:permissions"), Resource: "permission", Action: "list", Status: entity.PermissionStatusActive},
		{Name: "权限-详情", Code: "system:permission:get", MenuID: getMenuID("system:permissions"), Resource: "permission", Action: "get", Status: entity.PermissionStatusActive},
		{Name: "权限-创建", Code: "system:permission:create", MenuID: getMenuID("system:permissions"), Resource: "permission", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "权限-更新", Code: "system:permission:update", MenuID: getMenuID("system:permissions"), Resource: "permission", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "权限-删除", Code: "system:permission:delete", MenuID: getMenuID("system:permissions"), Resource: "permission", Action: "delete", Status: entity.PermissionStatusActive},

		// 菜单管理 (system:menus)
		{Name: "菜单-查询", Code: "system:menu:list", MenuID: getMenuID("system:menus"), Resource: "menu", Action: "list", Status: entity.PermissionStatusActive},
		{Name: "菜单-详情", Code: "system:menu:get", MenuID: getMenuID("system:menus"), Resource: "menu", Action: "get", Status: entity.PermissionStatusActive},
		{Name: "菜单-创建", Code: "system:menu:create", MenuID: getMenuID("system:menus"), Resource: "menu", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "菜单-更新", Code: "system:menu:update", MenuID: getMenuID("system:menus"), Resource: "menu", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "菜单-删除", Code: "system:menu:delete", MenuID: getMenuID("system:menus"), Resource: "menu", Action: "delete", Status: entity.PermissionStatusActive},
	}

	for _, req := range permissions {
		_, err := permissionService.CreatePermission(ctx, &req)
		if err != nil {
			fmt.Printf("创建权限 %s 失败: %v\n", req.Name, err)
		} else {
			fmt.Printf("创建权限 %s 成功\n", req.Name)
		}
	}

	return nil
}

// initMenus 初始化菜单
func initMenus(ctx context.Context, db *gorm.DB, menuService rbacService.MenuService) error {
	// ================= 精简版菜单结构 =================

	// 顶层分组：仪表板
	dashboardGroup, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		Name:   "仪表板",
		Code:   "dashboard",
		Type:   entity.MenuTypeGroup,
		Order:  1,
		Status: entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建仪表板分组失败: %w", err)
	}

	// 顶层分组：系统设置
	systemGroup, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		Name:   "系统设置",
		Code:   "system",
		Type:   entity.MenuTypeGroup,
		Order:  2,
		Status: entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建系统设置分组失败: %w", err)
	}

	// 修正分组 type（见原注释说明）
	if err := db.Model(&entity.Menu{}).Where("id IN ?", []string{dashboardGroup.ID, systemGroup.ID}).Update("type", entity.MenuTypeGroup).Error; err != nil {
		return fmt.Errorf("修正分组菜单类型失败: %w", err)
	}

	// 工作台菜单
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &dashboardGroup.ID,
		Name:     "工作台",
		Code:     "workbench",
		Path:     stringPtr("/workbench"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon:      stringPtr("local:ic-workbench"),
			Component: stringPtr("views/dashboard/workbench"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建工作台菜单失败: %w", err)
	}

	// 系统设置下的四个核心管理菜单
	// 菜单管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &systemGroup.ID,
		Name:     "菜单管理",
		Code:     "system:menus",
		Path:     stringPtr("/system/menus"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta:     entity.MenuMeta{Component: stringPtr("views/management/rbac/MenuManagement")},
	})
	if err != nil {
		return fmt.Errorf("创建菜单管理失败: %w", err)
	}

	// 角色管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &systemGroup.ID,
		Name:     "角色管理",
		Code:     "system:roles",
		Path:     stringPtr("/system/roles"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta:     entity.MenuMeta{Component: stringPtr("views/management/rbac/RoleManagement")},
	})
	if err != nil {
		return fmt.Errorf("创建角色管理失败: %w", err)
	}

	// 用户管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &systemGroup.ID,
		Name:     "用户管理",
		Code:     "system:users",
		Path:     stringPtr("/system/users"),
		Type:     entity.MenuTypeMenu,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta:     entity.MenuMeta{Component: stringPtr("views/management/rbac/UserManagement")},
	})
	if err != nil {
		return fmt.Errorf("创建用户管理失败: %w", err)
	}

	// 权限管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &systemGroup.ID,
		Name:     "权限管理",
		Code:     "system:permissions",
		Path:     stringPtr("/system/permissions"),
		Type:     entity.MenuTypeMenu,
		Order:    4,
		Status:   entity.MenuStatusActive,
		Meta:     entity.MenuMeta{Component: stringPtr("views/management/rbac/PermissionManagement")},
	})
	if err != nil {
		return fmt.Errorf("创建权限管理失败: %w", err)
	}

	fmt.Println("初始化菜单完成(精简版)")
	return nil
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// cleanExistingData 清理现有数据
func cleanExistingData(db *gorm.DB) error {
	fmt.Println("开始清理现有数据...")

	// 删除菜单数据（由于外键约束，需要先删除子菜单）
	if err := db.Exec("DELETE FROM menus").Error; err != nil {
		return fmt.Errorf("清理菜单数据失败: %w", err)
	}

	// 删除权限数据
	if err := db.Exec("DELETE FROM permissions").Error; err != nil {
		return fmt.Errorf("清理权限数据失败: %w", err)
	}

	// 删除角色数据
	if err := db.Exec("DELETE FROM roles").Error; err != nil {
		return fmt.Errorf("清理角色数据失败: %w", err)
	}

	fmt.Println("清理现有数据完成")
	return nil
}
