# 菜单初始化指南

## 概述

本文档介绍如何初始化 Alice 项目的菜单数据。菜单数据根据前端路由配置生成，包含完整的导航结构。

## 菜单结构

菜单系统采用四层结构：

### 1. 菜单类型
- **分组 (Group)**: 顶级分类，如"仪表板"、"页面管理"等
- **目录 (Catalogue)**: 中间层级，有子菜单但本身不是页面
- **菜单 (Menu)**: 具体的页面菜单项
- **按钮 (Button)**: 页面内的操作按钮权限

### 2. 分组组织
```
├── 仪表板 (Group)
│   ├── 工作台
│   └── 分析页
├── 页面管理 (Group)  
│   ├── 系统管理 (目录)
│   ├── 多级菜单 (目录)
│   └── 错误页面 (目录)
├── UI组件 (Group)
│   └── 组件 (目录)
└── 其他 (Group)
    ├── 各种演示菜单
    └── ...
```

## 初始化步骤

### 1. 确保数据库配置正确
检查 `config.yaml` 中的数据库连接配置：
```yaml
database:
  host: localhost
  port: 3306
  username: your_username
  password: your_password
  database: alice
```

### 2. 运行初始化命令
```bash
# 进入后端目录
cd backend

# 编译初始化程序
go build -o bin/init ./cmd/init

# 运行初始化（这将清空并重新创建所有数据）
./bin/init
```

### 3. 验证初始化结果
初始化完成后，检查数据库中的 `menus` 表，应该包含约 40+ 条记录。

## 菜单配置说明

### Meta 配置字段
```json
{
  "icon": "图标名称，如 local:ic-workbench",
  "caption": "菜单说明文字",
  "info": "标签信息，如 New、Hot 等", 
  "disabled": "是否禁用",
  "auth": "是否需要认证",
  "hidden": "是否隐藏",
  "external_link": "外部链接地址",
  "component": "前端组件路径"
}
```

### 路径规则
- 分组菜单：无路径
- 目录菜单：设置父级路径，如 `/management`
- 页面菜单：完整路径，如 `/management/user/profile`
- 特殊路径：`#label` 表示标签菜单，不对应实际页面

## 自定义菜单

### 1. 修改现有菜单
编辑 `cmd/init/main.go` 中的 `initMenus` 函数，修改对应的菜单配置。

### 2. 添加新菜单
在相应的分组或目录下添加新的菜单项：

```go
_, err = menuService.CreateMenu(ctx, &service.CreateMenuRequest{
    ParentID: &parentCatalogue.ID,
    Name:     "新菜单",
    Code:     "new:menu",
    Path:     stringPtr("/new/menu"),
    Type:     entity.MenuTypeMenu,
    Order:    1,
    Status:   entity.MenuStatusActive,
    Meta: entity.MenuMeta{
        Icon:      stringPtr("local:ic-new"),
        Component: stringPtr("/pages/new/menu"),
    },
})
```

### 3. 重新初始化
修改后重新编译和运行初始化命令。

## 注意事项

1. **数据清空**: 初始化会清空现有的用户、角色、权限、菜单数据
2. **权限关联**: 菜单会自动关联到管理员角色
3. **组件路径**: 确保前端组件路径与实际文件路径一致
4. **图标资源**: 确保使用的图标在前端项目中可用
5. **国际化**: 当前使用中文标题，已移除 i18n 支持

## 故障排除

### 常见错误
1. **数据库连接失败**: 检查配置文件和数据库服务状态
2. **外键约束错误**: 确保数据库表结构正确
3. **重复键错误**: 检查菜单代码是否唯一

### 调试建议
1. 查看初始化日志输出
2. 检查数据库错误日志
3. 验证菜单层级关系是否正确

## 相关文件

- `cmd/init/main.go` - 初始化逻辑
- `domain/rbac/entity/menu.go` - 菜单实体定义
- `frontend/src/routes/sections/dashboard/frontend.tsx` - 前端路由配置
- `frontend/src/layouts/dashboard/nav/nav-data/nav-data-frontend.tsx` - 导航数据配置
- `docs/menu_init_data.sql` - 菜单结构说明
