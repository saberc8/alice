import React, { useState, useMemo } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/ui/table';
import { Button } from '@/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card';
import { Badge } from '@/ui/badge';
import { Input } from '@/ui/input';
import { Label } from '@/ui/label';
import { Textarea } from '@/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/ui/select';
import { 
  Plus, 
  Edit, 
  Trash2, 
  Search,
  FolderTree,
  File,
  Folder,
  Square,
} from 'lucide-react';
import { toast } from 'sonner';

import menuService, { 
  type CreateMenuReq, 
  type UpdateMenuReq,
  type MenuListReq 
} from '@/api/services/menuService';
import type { Menu, MenuTree } from '@/types/entity';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

// 菜单类型图标映射
const getMenuTypeIcon = (type: number) => {
  switch (type) {
    case 0: return <FolderTree className="w-4 h-4" />; // 分组
    case 1: return <Folder className="w-4 h-4" />; // 目录
    case 2: return <File className="w-4 h-4" />; // 菜单
    case 3: return <Square className="w-4 h-4" />; // 按钮
    default: return <File className="w-4 h-4" />;
  }
};

// 菜单类型名称映射
const getMenuTypeName = (type: number) => {
  switch (type) {
    case 0: return '分组';
    case 1: return '目录';
    case 2: return '菜单';
    case 3: return '按钮';
    default: return '未知';
  }
};

// 菜单表单组件
interface MenuFormProps {
  menu?: Menu;
  parentMenus: MenuTree[];
  onSubmit: (data: CreateMenuReq | UpdateMenuReq) => void;
  isLoading?: boolean;
}

const MenuForm: React.FC<MenuFormProps> = ({ menu, parentMenus, onSubmit, isLoading }) => {
  const [formData, setFormData] = useState<CreateMenuReq>({
    parent_id: menu?.parent_id || '',
    name: menu?.name || '',
    code: menu?.code || '',
    path: menu?.path || '',
    type: menu?.type ?? 2, // 默认菜单类型
    order: menu?.order || 1,
    status: menu?.status || 'active',
    meta: menu?.meta || {},
  });

  const [metaJson, setMetaJson] = useState(
    JSON.stringify(formData.meta, null, 2)
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const meta = metaJson.trim() ? JSON.parse(metaJson) : {};
      onSubmit({ ...formData, meta });
    } catch (error) {
      toast.error('Meta数据格式错误，请检查JSON格式');
    }
  };

  // 构建菜单树选项
  const renderMenuOptions = (menus: MenuTree[], level = 0): React.ReactNode[] => {
    return menus.reduce((options: React.ReactNode[], menuItem) => {
      // 不能选择自己作为父级
      if (menu && menuItem.id === menu.id) {
        return options;
      }
      
      const indent = '　'.repeat(level);
      const typeIcon = getMenuTypeIcon(menuItem.type);
      
      options.push(
        <SelectItem key={menuItem.id} value={menuItem.id}>
          <div className="flex items-center gap-2">
            {typeIcon}
            {indent}{menuItem.name}
          </div>
        </SelectItem>
      );
      
      if (menuItem.children) {
        options.push(...renderMenuOptions(menuItem.children, level + 1));
      }
      
      return options;
    }, []);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-h-[600px] overflow-y-auto">
      <div className="space-y-2">
        <Label htmlFor="parent_id">父级菜单</Label>
        <Select
          value={formData.parent_id}
          onValueChange={(value) => setFormData({ ...formData, parent_id: value })}
        >
          <SelectTrigger>
            <SelectValue placeholder="选择父级菜单（可为空）" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">根级菜单</SelectItem>
            {renderMenuOptions(parentMenus)}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">菜单名称</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="请输入菜单名称"
          required
        />
      </div>
      
      <div className="space-y-2">
        <Label htmlFor="code">菜单代码</Label>
        <Input
          id="code"
          value={formData.code}
          onChange={(e) => setFormData({ ...formData, code: e.target.value })}
          placeholder="请输入菜单代码"
          required
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="type">菜单类型</Label>
          <Select
            value={formData.type.toString()}
            onValueChange={(value) => setFormData({ ...formData, type: parseInt(value) })}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="0">
                <div className="flex items-center gap-2">
                  <FolderTree className="w-4 h-4" />
                  分组
                </div>
              </SelectItem>
              <SelectItem value="1">
                <div className="flex items-center gap-2">
                  <Folder className="w-4 h-4" />
                  目录
                </div>
              </SelectItem>
              <SelectItem value="2">
                <div className="flex items-center gap-2">
                  <File className="w-4 h-4" />
                  菜单
                </div>
              </SelectItem>
              <SelectItem value="3">
                <div className="flex items-center gap-2">
                  <Square className="w-4 h-4" />
                  按钮
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="order">排序</Label>
          <Input
            id="order"
            type="number"
            value={formData.order}
            onChange={(e) => setFormData({ ...formData, order: parseInt(e.target.value) || 1 })}
            placeholder="排序值"
            min="1"
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="path">路径</Label>
        <Input
          id="path"
          value={formData.path}
          onChange={(e) => setFormData({ ...formData, path: e.target.value })}
          placeholder="请输入路径（如：/user/list）"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="status">状态</Label>
        <Select
          value={formData.status}
          onValueChange={(value: 'active' | 'inactive') => setFormData({ ...formData, status: value })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="active">激活</SelectItem>
            <SelectItem value="inactive">禁用</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="meta">Meta数据（JSON格式）</Label>
        <Textarea
          id="meta"
          value={metaJson}
          onChange={(e) => setMetaJson(e.target.value)}
          placeholder='{"icon": "user", "component": "/pages/user/list"}'
          rows={6}
          className="font-mono text-sm"
        />
        <div className="text-sm text-gray-500">
          用于存储图标、组件路径等额外信息，格式为JSON
        </div>
      </div>

      <DialogFooter>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? '提交中...' : '确定'}
        </Button>
      </DialogFooter>
    </form>
  );
};

// 菜单树表格行组件
interface MenuTreeRowProps {
  menu: MenuTree;
  level: number;
  onEdit: (menu: Menu) => void;
  onDelete: (menu: Menu) => void;
}

const MenuTreeRow: React.FC<MenuTreeRowProps> = ({ menu, level, onEdit, onDelete }) => {
  const indent = '　'.repeat(level);
  const typeIcon = getMenuTypeIcon(menu.type);
  
  return (
    <>
      <TableRow>
        <TableCell>
          <div className="flex items-center gap-2">
            {typeIcon}
            {indent}{menu.name}
          </div>
        </TableCell>
        <TableCell>
          <code className="bg-gray-100 px-2 py-1 rounded text-sm">
            {menu.code}
          </code>
        </TableCell>
        <TableCell>{menu.path || '-'}</TableCell>
        <TableCell>
          <Badge variant="outline">{getMenuTypeName(menu.type)}</Badge>
        </TableCell>
        <TableCell>{menu.order}</TableCell>
        <TableCell>
          <Badge variant={menu.status === 'active' ? 'default' : 'secondary'}>
            {menu.status === 'active' ? '激活' : '禁用'}
          </Badge>
        </TableCell>
        <TableCell>
          {menu.created_at ? new Date(menu.created_at).toLocaleDateString() : '-'}
        </TableCell>
        <TableCell className="text-right">
          <div className="flex gap-2 justify-end">
            <PermissionWrapper resource="menu" action="update">
              <Button
                variant="outline"
                size="sm"
                onClick={() => onEdit(menu)}
              >
                <Edit className="w-4 h-4" />
              </Button>
            </PermissionWrapper>
            
            <PermissionWrapper resource="menu" action="delete">
              <Button
                variant="destructive"
                size="sm"
                onClick={() => onDelete(menu)}
              >
                <Trash2 className="w-4 h-4" />
              </Button>
            </PermissionWrapper>
          </div>
        </TableCell>
      </TableRow>
      {menu.children?.map((child) => (
        <MenuTreeRow
          key={child.id}
          menu={child}
          level={level + 1}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </>
  );
};

// 菜单管理主组件
export const MenuManagement: React.FC = () => {
  const [searchParams, setSearchParams] = useState<MenuListReq>({
    page: 1,
    page_size: 100, // 菜单一般不多，可以一次性加载
  });
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [selectedMenu, setSelectedMenu] = useState<Menu | null>(null);

  const queryClient = useQueryClient();

  // 获取菜单树
  const { data: menuTreeData, isLoading } = useQuery({
    queryKey: ['menu-tree'],
    queryFn: () => menuService.getMenuTree(),
  });

  // 获取菜单列表（用于父级选择）
  const { data: menusData } = useQuery({
    queryKey: ['menus', searchParams],
    queryFn: () => menuService.getMenus(searchParams),
  });

  // 创建菜单
  const createMutation = useMutation({
    mutationFn: menuService.createMenu,
    onSuccess: () => {
      toast.success('菜单创建成功');
      queryClient.invalidateQueries({ queryKey: ['menu-tree'] });
      queryClient.invalidateQueries({ queryKey: ['menus'] });
      setIsCreateDialogOpen(false);
    },
    onError: (error: any) => {
      toast.error(error.message || '菜单创建失败');
    },
  });

  // 更新菜单
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMenuReq }) =>
      menuService.updateMenu(id, data),
    onSuccess: () => {
      toast.success('菜单更新成功');
      queryClient.invalidateQueries({ queryKey: ['menu-tree'] });
      queryClient.invalidateQueries({ queryKey: ['menus'] });
      setIsEditDialogOpen(false);
      setSelectedMenu(null);
    },
    onError: (error: any) => {
      toast.error(error.message || '菜单更新失败');
    },
  });

  // 删除菜单
  const deleteMutation = useMutation({
    mutationFn: menuService.deleteMenu,
    onSuccess: () => {
      toast.success('菜单删除成功');
      queryClient.invalidateQueries({ queryKey: ['menu-tree'] });
      queryClient.invalidateQueries({ queryKey: ['menus'] });
    },
    onError: (error: any) => {
      toast.error(error.message || '菜单删除失败');
    },
  });

  const handleEdit = (menu: Menu) => {
    setSelectedMenu(menu);
    setIsEditDialogOpen(true);
  };

  const handleDelete = (menu: Menu) => {
    if (window.confirm(`确定要删除菜单 "${menu.name}" 吗？`)) {
      deleteMutation.mutate(menu.id);
    }
  };

  const menuTree = menuTreeData?.data || [];
  const parentMenus = menuTreeData?.data || [];

  // 计算菜单总数
  const countMenus = (menus: MenuTree[]): number => {
    return menus.reduce((count, menu) => {
      return count + 1 + (menu.children ? countMenus(menu.children) : 0);
    }, 0);
  };

  const totalMenus = countMenus(menuTree);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">菜单管理</h1>
        <PermissionWrapper resource="menu" action="create">
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                创建菜单
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-3xl">
              <DialogHeader>
                <DialogTitle>创建菜单</DialogTitle>
                <DialogDescription>
                  创建新的系统菜单
                </DialogDescription>
              </DialogHeader>
              <MenuForm
                parentMenus={parentMenus}
                onSubmit={(data) => createMutation.mutate(data)}
                isLoading={createMutation.isPending}
              />
            </DialogContent>
          </Dialog>
        </PermissionWrapper>
      </div>

      {/* 菜单树表格 */}
      <Card>
        <CardHeader>
          <CardTitle>菜单树结构 (共 {totalMenus} 个菜单)</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>菜单名称</TableHead>
                <TableHead>菜单代码</TableHead>
                <TableHead>路径</TableHead>
                <TableHead>类型</TableHead>
                <TableHead>排序</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>创建时间</TableHead>
                <TableHead className="text-right">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8">
                    加载中...
                  </TableCell>
                </TableRow>
              ) : menuTree.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8">
                    暂无数据
                  </TableCell>
                </TableRow>
              ) : (
                menuTree.map((menu) => (
                  <MenuTreeRow
                    key={menu.id}
                    menu={menu}
                    level={0}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* 编辑菜单对话框 */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>编辑菜单</DialogTitle>
            <DialogDescription>
              修改菜单信息
            </DialogDescription>
          </DialogHeader>
          {selectedMenu && (
            <MenuForm
              menu={selectedMenu}
              parentMenus={parentMenus}
              onSubmit={(data) => updateMutation.mutate({ id: selectedMenu.id, data })}
              isLoading={updateMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default MenuManagement;
