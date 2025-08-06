import React, { useState } from 'react';
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
} from 'lucide-react';
import { toast } from 'sonner';

import permissionService, { 
  type CreatePermissionReq, 
  type UpdatePermissionReq,
  type PermissionListReq 
} from '@/api/services/permissionService';
import type { Permission } from '@/types/entity';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

// 权限表单组件
interface PermissionFormProps {
  permission?: Permission;
  onSubmit: (data: CreatePermissionReq | UpdatePermissionReq) => void;
  isLoading?: boolean;
}

const PermissionForm: React.FC<PermissionFormProps> = ({ permission, onSubmit, isLoading }) => {
  const [formData, setFormData] = useState<CreatePermissionReq>({
    name: permission?.name || '',
    code: permission?.code || '',
    resource: permission?.resource || '',
    action: permission?.action || '',
    description: permission?.description || '',
    status: permission?.status || 'active',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  // 自动生成权限代码
  const generateCode = () => {
    if (formData.resource && formData.action) {
      setFormData(prev => ({ ...prev, code: `${prev.resource}:${prev.action}` }));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">权限名称</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="请输入权限名称"
          required
        />
      </div>
      
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="resource">资源</Label>
          <Input
            id="resource"
            value={formData.resource}
            onChange={(e) => setFormData({ ...formData, resource: e.target.value })}
            onBlur={generateCode}
            placeholder="如：user"
            required
          />
        </div>
        
        <div className="space-y-2">
          <Label htmlFor="action">操作</Label>
          <Select
            value={formData.action}
            onValueChange={(value) => {
              setFormData({ ...formData, action: value });
              setTimeout(generateCode, 0);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="选择操作" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="create">create (创建)</SelectItem>
              <SelectItem value="read">read (查看)</SelectItem>
              <SelectItem value="update">update (更新)</SelectItem>
              <SelectItem value="delete">delete (删除)</SelectItem>
              <SelectItem value="manage">manage (管理)</SelectItem>
              <SelectItem value="assign">assign (分配)</SelectItem>
              <SelectItem value="export">export (导出)</SelectItem>
              <SelectItem value="import">import (导入)</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="code">权限代码</Label>
        <Input
          id="code"
          value={formData.code}
          onChange={(e) => setFormData({ ...formData, code: e.target.value })}
          placeholder="如：user:read"
          required
        />
        <div className="text-sm text-gray-500">
          格式：resource:action，如 user:read
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">描述</Label>
        <Textarea
          id="description"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="请输入权限描述"
          rows={3}
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

      <DialogFooter>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? '提交中...' : '确定'}
        </Button>
      </DialogFooter>
    </form>
  );
};

// 权限管理主组件
export const PermissionManagement: React.FC = () => {
  const [searchParams, setSearchParams] = useState<PermissionListReq>({
    page: 1,
    page_size: 10,
  });
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [selectedPermission, setSelectedPermission] = useState<Permission | null>(null);

  const queryClient = useQueryClient();

  // 获取权限列表
  const { data: permissionsData, isLoading } = useQuery({
    queryKey: ['permissions', searchParams],
    queryFn: () => permissionService.getPermissions(searchParams),
  });

  // 创建权限
  const createMutation = useMutation({
    mutationFn: permissionService.createPermission,
    onSuccess: () => {
      toast.success('权限创建成功');
      queryClient.invalidateQueries({ queryKey: ['permissions'] });
      setIsCreateDialogOpen(false);
    },
    onError: (error: any) => {
      toast.error(error.message || '权限创建失败');
    },
  });

  // 更新权限
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePermissionReq }) =>
      permissionService.updatePermission(id, data),
    onSuccess: () => {
      toast.success('权限更新成功');
      queryClient.invalidateQueries({ queryKey: ['permissions'] });
      setIsEditDialogOpen(false);
      setSelectedPermission(null);
    },
    onError: (error: any) => {
      toast.error(error.message || '权限更新失败');
    },
  });

  // 删除权限
  const deleteMutation = useMutation({
    mutationFn: permissionService.deletePermission,
    onSuccess: () => {
      toast.success('权限删除成功');
      queryClient.invalidateQueries({ queryKey: ['permissions'] });
    },
    onError: (error: any) => {
      toast.error(error.message || '权限删除失败');
    },
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ ...searchParams, page: 1 });
  };

  const handleEdit = (permission: Permission) => {
    setSelectedPermission(permission);
    setIsEditDialogOpen(true);
  };

  const handleDelete = (permission: Permission) => {
    if (window.confirm(`确定要删除权限 "${permission.name}" 吗？`)) {
      deleteMutation.mutate(permission.id);
    }
  };

  const permissions = permissionsData?.data?.list || [];
  const total = permissionsData?.data?.total || 0;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">权限管理</h1>
        <PermissionWrapper resource="permission" action="create">
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                创建权限
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>创建权限</DialogTitle>
                <DialogDescription>
                  创建新的系统权限
                </DialogDescription>
              </DialogHeader>
              <PermissionForm
                onSubmit={(data) => createMutation.mutate(data)}
                isLoading={createMutation.isPending}
              />
            </DialogContent>
          </Dialog>
        </PermissionWrapper>
      </div>

      {/* 搜索表单 */}
      <Card>
        <CardHeader>
          <CardTitle>搜索条件</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSearch} className="flex gap-4 items-end">
            <div className="flex-1">
              <Label htmlFor="search-name">权限名称</Label>
              <Input
                id="search-name"
                placeholder="请输入权限名称"
                value={searchParams.name || ''}
                onChange={(e) => setSearchParams({ ...searchParams, name: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-code">权限代码</Label>
              <Input
                id="search-code"
                placeholder="请输入权限代码"
                value={searchParams.code || ''}
                onChange={(e) => setSearchParams({ ...searchParams, code: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-resource">资源</Label>
              <Input
                id="search-resource"
                placeholder="请输入资源"
                value={searchParams.resource || ''}
                onChange={(e) => setSearchParams({ ...searchParams, resource: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-action">操作</Label>
              <Input
                id="search-action"
                placeholder="请输入操作"
                value={searchParams.action || ''}
                onChange={(e) => setSearchParams({ ...searchParams, action: e.target.value })}
              />
            </div>
            <Button type="submit">
              <Search className="w-4 h-4 mr-2" />
              搜索
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* 权限列表 */}
      <Card>
        <CardHeader>
          <CardTitle>权限列表 (共 {total} 条)</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>权限名称</TableHead>
                <TableHead>权限代码</TableHead>
                <TableHead>资源</TableHead>
                <TableHead>操作</TableHead>
                <TableHead>描述</TableHead>
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
              ) : permissions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8">
                    暂无数据
                  </TableCell>
                </TableRow>
              ) : (
                permissions.map((permission) => (
                  <TableRow key={permission.id}>
                    <TableCell className="font-medium">{permission.name}</TableCell>
                    <TableCell>
                      <code className="bg-gray-100 px-2 py-1 rounded text-sm">
                        {permission.code}
                      </code>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">{permission.resource}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant="secondary">{permission.action}</Badge>
                    </TableCell>
                    <TableCell className="max-w-xs truncate">
                      {permission.description || '-'}
                    </TableCell>
                    <TableCell>
                      <Badge variant={permission.status === 'active' ? 'default' : 'secondary'}>
                        {permission.status === 'active' ? '激活' : '禁用'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {permission.created_at ? new Date(permission.created_at).toLocaleDateString() : '-'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex gap-2 justify-end">
                        <PermissionWrapper resource="permission" action="update">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEdit(permission)}
                          >
                            <Edit className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                        
                        <PermissionWrapper resource="permission" action="delete">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDelete(permission)}
                          >
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* 编辑权限对话框 */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>编辑权限</DialogTitle>
            <DialogDescription>
              修改权限信息
            </DialogDescription>
          </DialogHeader>
          {selectedPermission && (
            <PermissionForm
              permission={selectedPermission}
              onSubmit={(data) => updateMutation.mutate({ id: selectedPermission.id, data })}
              isLoading={updateMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default PermissionManagement;
