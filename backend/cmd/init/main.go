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
	menuService := rbacService.NewMenuService(menuRepo)
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

	if err := initPermissions(ctx, permissionService); err != nil {
		log.Fatal("Failed to init permissions:", err)
	}

	if err := initMenus(ctx, db, menuService); err != nil {
		log.Fatal("Failed to init menus:", err)
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
func initPermissions(ctx context.Context, permissionService rbacService.PermissionService) error {
	permissions := []rbacService.CreatePermissionRequest{
		// 用户管理权限
		{Name: "查看用户", Code: "user:read", Resource: "user", Action: "read", Status: entity.PermissionStatusActive},
		{Name: "创建用户", Code: "user:create", Resource: "user", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "更新用户", Code: "user:update", Resource: "user", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "删除用户", Code: "user:delete", Resource: "user", Action: "delete", Status: entity.PermissionStatusActive},

		// 角色管理权限
		{Name: "查看角色", Code: "role:read", Resource: "role", Action: "read", Status: entity.PermissionStatusActive},
		{Name: "创建角色", Code: "role:create", Resource: "role", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "更新角色", Code: "role:update", Resource: "role", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "删除角色", Code: "role:delete", Resource: "role", Action: "delete", Status: entity.PermissionStatusActive},

		// 权限管理权限
		{Name: "查看权限", Code: "permission:read", Resource: "permission", Action: "read", Status: entity.PermissionStatusActive},
		{Name: "创建权限", Code: "permission:create", Resource: "permission", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "更新权限", Code: "permission:update", Resource: "permission", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "删除权限", Code: "permission:delete", Resource: "permission", Action: "delete", Status: entity.PermissionStatusActive},

		// 菜单管理权限
		{Name: "查看菜单", Code: "menu:read", Resource: "menu", Action: "read", Status: entity.PermissionStatusActive},
		{Name: "创建菜单", Code: "menu:create", Resource: "menu", Action: "create", Status: entity.PermissionStatusActive},
		{Name: "更新菜单", Code: "menu:update", Resource: "menu", Action: "update", Status: entity.PermissionStatusActive},
		{Name: "删除菜单", Code: "menu:delete", Resource: "menu", Action: "delete", Status: entity.PermissionStatusActive},
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
	// ========== 创建分组 ==========

	// 1. 仪表板分组
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

	// 2. 页面分组
	pagesGroup, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		Name:   "页面管理",
		Code:   "pages",
		Type:   entity.MenuTypeGroup,
		Order:  2,
		Status: entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建页面分组失败: %w", err)
	}

	// 3. UI组件分组
	uiGroup, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		Name:   "UI组件",
		Code:   "ui",
		Type:   entity.MenuTypeGroup,
		Order:  3,
		Status: entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建UI组件分组失败: %w", err)
	}

	// 4. 其他分组
	othersGroup, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		Name:   "其他",
		Code:   "others",
		Type:   entity.MenuTypeGroup,
		Order:  4,
		Status: entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建其他分组失败: %w", err)
	}

	// 由于 GORM 对零值字段不会生成 INSERT 列，且实体设置了 default:2，会被数据库默认覆盖。
	// 这里在创建完四个分组后，统一强制更新它们的 type=0。
	groupIDs := []string{dashboardGroup.ID, pagesGroup.ID, uiGroup.ID, othersGroup.ID}
	if err := db.Model(&entity.Menu{}).Where("id IN ?", groupIDs).Update("type", entity.MenuTypeGroup).Error; err != nil {
		return fmt.Errorf("修正分组菜单类型为0失败: %w", err)
	}

	// ========== 仪表板菜单 ==========

	// 工作台
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
			Component: stringPtr("/pages/dashboard/workbench"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建工作台菜单失败: %w", err)
	}

	// 分析页
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &dashboardGroup.ID,
		Name:     "分析页",
		Code:     "analysis",
		Path:     stringPtr("/analysis"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon:      stringPtr("local:ic-analysis"),
			Component: stringPtr("/pages/dashboard/analysis"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建分析页菜单失败: %w", err)
	}

	// ========== 页面管理 ==========

	// 系统管理目录
	managementCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &pagesGroup.ID,
		Name:     "系统管理",
		Code:     "management",
		Path:     stringPtr("/management"),
		Type:     entity.MenuTypeCatalogue,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon: stringPtr("local:ic-management"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建系统管理目录失败: %w", err)
	}

	// 用户管理目录
	userCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &managementCatalogue.ID,
		Name:     "用户管理",
		Code:     "management:user",
		Path:     stringPtr("/management/user"),
		Type:     entity.MenuTypeCatalogue,
		Order:    1,
		Status:   entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建用户管理目录失败: %w", err)
	}

	// 用户资料
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &userCatalogue.ID,
		Name:     "用户资料",
		Code:     "management:user:profile",
		Path:     stringPtr("/management/user/profile"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/user/profile"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建用户资料菜单失败: %w", err)
	}

	// 账户管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &userCatalogue.ID,
		Name:     "账户管理",
		Code:     "management:user:account",
		Path:     stringPtr("/management/user/account"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/user/account"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建账户管理菜单失败: %w", err)
	}

	// RBAC管理目录
	rbacCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &managementCatalogue.ID,
		Name:     "权限管理",
		Code:     "management:rbac",
		Path:     stringPtr("/management/rbac"),
		Type:     entity.MenuTypeCatalogue,
		Order:    2,
		Status:   entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建RBAC管理目录失败: %w", err)
	}

	// RBAC概览
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "权限概览",
		Code:     "management:rbac:overview",
		Path:     stringPtr("/management/rbac"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建权限概览菜单失败: %w", err)
	}

	// 用户管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "用户管理",
		Code:     "management:rbac:users",
		Path:     stringPtr("/management/rbac/users"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac/UserManagement"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建用户管理菜单失败: %w", err)
	}

	// 角色管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "角色管理",
		Code:     "management:rbac:roles",
		Path:     stringPtr("/management/rbac/roles"),
		Type:     entity.MenuTypeMenu,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac/RoleManagement"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建角色管理菜单失败: %w", err)
	}

	// 权限管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "权限管理",
		Code:     "management:rbac:permissions",
		Path:     stringPtr("/management/rbac/permissions"),
		Type:     entity.MenuTypeMenu,
		Order:    4,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac/PermissionManagement"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建权限管理菜单失败: %w", err)
	}

	// 菜单管理
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "菜单管理",
		Code:     "management:rbac:menus",
		Path:     stringPtr("/management/rbac/menus"),
		Type:     entity.MenuTypeMenu,
		Order:    5,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac/MenuManagement"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建菜单管理菜单失败: %w", err)
	}

	// RBAC演示
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &rbacCatalogue.ID,
		Name:     "权限演示",
		Code:     "management:rbac:demo",
		Path:     stringPtr("/management/rbac/demo"),
		Type:     entity.MenuTypeMenu,
		Order:    6,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/management/rbac/demo"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建权限演示菜单失败: %w", err)
	}

	// 多级菜单目录
	menuLevelCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &pagesGroup.ID,
		Name:     "多级菜单",
		Code:     "menu_level",
		Path:     stringPtr("/menu_level"),
		Type:     entity.MenuTypeCatalogue,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon: stringPtr("local:ic-menulevel"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建多级菜单目录失败: %w", err)
	}

	// 菜单1-a
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menuLevelCatalogue.ID,
		Name:     "菜单1-a",
		Code:     "menu_level:1a",
		Path:     stringPtr("/menu_level/1a"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/menu-level/menu-level-1a"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建菜单1-a失败: %w", err)
	}

	// 菜单1-b目录
	menu1bCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menuLevelCatalogue.ID,
		Name:     "菜单1-b",
		Code:     "menu_level:1b",
		Path:     stringPtr("/menu_level/1b"),
		Type:     entity.MenuTypeCatalogue,
		Order:    2,
		Status:   entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建菜单1-b目录失败: %w", err)
	}

	// 菜单2-a
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menu1bCatalogue.ID,
		Name:     "菜单2-a",
		Code:     "menu_level:1b:2a",
		Path:     stringPtr("/menu_level/1b/2a"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/menu-level/menu-level-1b/menu-level-2a"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建菜单2-a失败: %w", err)
	}

	// 菜单2-b目录
	menu2bCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menu1bCatalogue.ID,
		Name:     "菜单2-b",
		Code:     "menu_level:1b:2b",
		Path:     stringPtr("/menu_level/1b/2b"),
		Type:     entity.MenuTypeCatalogue,
		Order:    2,
		Status:   entity.MenuStatusActive,
	})
	if err != nil {
		return fmt.Errorf("创建菜单2-b目录失败: %w", err)
	}

	// 菜单3-a
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menu2bCatalogue.ID,
		Name:     "菜单3-a",
		Code:     "menu_level:1b:2b:3a",
		Path:     stringPtr("/menu_level/1b/2b/3a"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/menu-level/menu-level-1b/menu-level-2b/menu-level-3a"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建菜单3-a失败: %w", err)
	}

	// 菜单3-b
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &menu2bCatalogue.ID,
		Name:     "菜单3-b",
		Code:     "menu_level:1b:2b:3b",
		Path:     stringPtr("/menu_level/1b/2b/3b"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/menu-level/menu-level-1b/menu-level-2b/menu-level-3b"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建菜单3-b失败: %w", err)
	}

	// 错误页面目录
	errorCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &pagesGroup.ID,
		Name:     "错误页面",
		Code:     "error",
		Path:     stringPtr("/error"),
		Type:     entity.MenuTypeCatalogue,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon: stringPtr("bxs:error-alt"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建错误页面目录失败: %w", err)
	}

	// 403页面
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &errorCatalogue.ID,
		Name:     "403无权限",
		Code:     "error:403",
		Path:     stringPtr("/error/403"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/sys/error/Page403"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建403页面失败: %w", err)
	}

	// 404页面
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &errorCatalogue.ID,
		Name:     "404未找到",
		Code:     "error:404",
		Path:     stringPtr("/error/404"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/sys/error/Page404"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建404页面失败: %w", err)
	}

	// 500页面
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &errorCatalogue.ID,
		Name:     "500服务器错误",
		Code:     "error:500",
		Path:     stringPtr("/error/500"),
		Type:     entity.MenuTypeMenu,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/sys/error/Page500"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建500页面失败: %w", err)
	}

	// ========== UI组件 ==========

	// 组件目录
	componentsCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &uiGroup.ID,
		Name:     "组件",
		Code:     "components",
		Path:     stringPtr("/components"),
		Type:     entity.MenuTypeCatalogue,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon:    stringPtr("solar:widget-5-bold-duotone"),
			Caption: stringPtr("自定义UI组件"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建组件目录失败: %w", err)
	}

	// 图标组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "图标",
		Code:     "components:icon",
		Path:     stringPtr("/components/icon"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/icon"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建图标组件菜单失败: %w", err)
	}

	// 动画组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "动画",
		Code:     "components:animate",
		Path:     stringPtr("/components/animate"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/animate"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建动画组件菜单失败: %w", err)
	}

	// 滚动组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "滚动",
		Code:     "components:scroll",
		Path:     stringPtr("/components/scroll"),
		Type:     entity.MenuTypeMenu,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/scroll"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建滚动组件菜单失败: %w", err)
	}

	// 上传组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "上传",
		Code:     "components:upload",
		Path:     stringPtr("/components/upload"),
		Type:     entity.MenuTypeMenu,
		Order:    4,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/upload"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建上传组件菜单失败: %w", err)
	}

	// 图表组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "图表",
		Code:     "components:chart",
		Path:     stringPtr("/components/chart"),
		Type:     entity.MenuTypeMenu,
		Order:    5,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/chart"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建图表组件菜单失败: %w", err)
	}

	// 消息提示组件
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &componentsCatalogue.ID,
		Name:     "消息提示",
		Code:     "components:toast",
		Path:     stringPtr("/components/toast"),
		Type:     entity.MenuTypeMenu,
		Order:    6,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/components/toast"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建消息提示组件菜单失败: %w", err)
	}

	// ========== 其他 ==========

	// 禁用菜单（演示用）
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &othersGroup.ID,
		Name:     "禁用菜单",
		Code:     "disabled",
		Path:     stringPtr("/disabled"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon:     stringPtr("local:ic-disabled"),
			Disabled: boolPtr(true),
		},
	})
	if err != nil {
		return fmt.Errorf("创建禁用菜单失败: %w", err)
	}

	// 标签菜单（演示用）
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &othersGroup.ID,
		Name:     "标签菜单",
		Code:     "label",
		Path:     stringPtr("#label"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon: stringPtr("local:ic-label"),
			Info: stringPtr("New"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建标签菜单失败: %w", err)
	}

	// 外部链接目录
	linkCatalogue, err := menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &othersGroup.ID,
		Name:     "外部链接",
		Code:     "link",
		Path:     stringPtr("/link"),
		Type:     entity.MenuTypeCatalogue,
		Order:    3,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon: stringPtr("local:ic-external"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建外部链接目录失败: %w", err)
	}

	// 外部链接
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &linkCatalogue.ID,
		Name:     "外部链接",
		Code:     "link:external",
		Path:     stringPtr("/link/external-link"),
		Type:     entity.MenuTypeMenu,
		Order:    1,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component:    stringPtr("/pages/sys/others/link/external-link"),
			ExternalLink: stringPtr("https://ant.design/index-cn"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建外部链接菜单失败: %w", err)
	}

	// 内嵌页面
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &linkCatalogue.ID,
		Name:     "内嵌页面",
		Code:     "link:iframe",
		Path:     stringPtr("/link/iframe"),
		Type:     entity.MenuTypeMenu,
		Order:    2,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Component: stringPtr("/pages/sys/others/link/iframe"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建内嵌页面菜单失败: %w", err)
	}

	// 空白页
	_, err = menuService.CreateMenu(ctx, &rbacService.CreateMenuRequest{
		ParentID: &othersGroup.ID,
		Name:     "空白页",
		Code:     "blank",
		Path:     stringPtr("/blank"),
		Type:     entity.MenuTypeMenu,
		Order:    4,
		Status:   entity.MenuStatusActive,
		Meta: entity.MenuMeta{
			Icon:      stringPtr("local:ic-blank"),
			Component: stringPtr("/pages/sys/others/blank"),
		},
	})
	if err != nil {
		return fmt.Errorf("创建空白页菜单失败: %w", err)
	}

	fmt.Println("初始化菜单完成")
	return nil
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// boolPtr 返回布尔值指针
func boolPtr(b bool) *bool {
	return &b
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
