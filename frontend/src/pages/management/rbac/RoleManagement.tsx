import React, { useState, useEffect } from 'react';
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
  Shield,
  Search,
  UserCheck,
} from 'lucide-react';
import { toast } from 'sonner';

import roleService, { 
  type CreateRoleReq, 
  type UpdateRoleReq,
  type AssignRolePermissionsReq,
  type RoleListReq 
} from '@/api/services/roleService';
import permissionService from '@/api/services/permissionService';
import type { Role, Permission } from '@/types/entity';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

// 角色表单组件
interface RoleFormProps {
  role?: Role;
  onSubmit: (data: CreateRoleReq | UpdateRoleReq) => void;
  isLoading?: boolean;
}

const RoleForm: React.FC<RoleFormProps> = ({ role, onSubmit, isLoading }) => {
  const [formData, setFormData] = useState<CreateRoleReq>({
    name: role?.name || '',
    code: role?.code || '',
    description: role?.description || '',
    status: role?.status || 'active',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">角色名称</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="请输入角色名称"
          required
        />
      </div>
      
      <div className="space-y-2">
        <Label htmlFor="code">角色代码</Label>
        <Input
          id="code"
          value={formData.code}
          onChange={(e) => setFormData({ ...formData, code: e.target.value })}
          placeholder="请输入角色代码（如：admin）"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">描述</Label>
        <Textarea
          id="description"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="请输入角色描述"
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

// 权限分配组件
interface PermissionAssignProps {
  role: Role;
  onClose: () => void;
}

const PermissionAssign: React.FC<PermissionAssignProps> = ({ role, onClose }) => {
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const queryClient = useQueryClient();

  // 获取所有权限
  const { data: allPermissions } = useQuery({
    queryKey: ['permissions', 'all'],
    queryFn: () => permissionService.getPermissions({ page: 1, page_size: 1000 }),
  });

  // 获取角色当前权限
  const { data: rolePermissions } = useQuery({
    queryKey: ['role-permissions', role.id],
    queryFn: () => roleService.getRolePermissions(role.id),
  });

  // 分配权限
  const assignMutation = useMutation({
    mutationFn: (data: AssignRolePermissionsReq) => 
      roleService.assignRolePermissions(role.id, data),
    onSuccess: () => {
      toast.success('权限分配成功');
      queryClient.invalidateQueries({ queryKey: ['role-permissions', role.id] });
      onClose();
    },
    onError: (error: any) => {
      toast.error(error.message || '权限分配失败');
    },
  });

  useEffect(() => {
    if (rolePermissions?.data) {
      setSelectedPermissions(rolePermissions.data.map(p => p.id));
    }
  }, [rolePermissions]);

  const handleTogglePermission = (permissionId: string) => {
    setSelectedPermissions(prev => 
      prev.includes(permissionId)
        ? prev.filter(id => id !== permissionId)
        : [...prev, permissionId]
    );
  };

  const handleSubmit = () => {
    assignMutation.mutate({ permission_ids: selectedPermissions });
  };

  const permissions = allPermissions?.data?.list || [];

  return (
    <div className="space-y-4">
      <div className="text-sm text-gray-600">
        为角色 "{role.name}" 分配权限
      </div>
      
      <div className="max-h-96 overflow-y-auto space-y-2">
        {permissions.map((permission) => (
          <div key={permission.id} className="flex items-center space-x-2">
            <input
              type="checkbox"
              id={permission.id}
              checked={selectedPermissions.includes(permission.id)}
              onChange={() => handleTogglePermission(permission.id)}
              className="rounded"
            />
            <label htmlFor={permission.id} className="flex-1 cursor-pointer">
              <div className="font-medium">{permission.name}</div>
              <div className="text-sm text-gray-500">
                {permission.code} - {permission.description}
              </div>
            </label>
          </div>
        ))}
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={onClose}>
          取消
        </Button>
        <Button 
          onClick={handleSubmit} 
          disabled={assignMutation.isPending}
        >
          {assignMutation.isPending ? '分配中...' : '确定'}
        </Button>
      </DialogFooter>
    </div>
  );
};

// 角色管理主组件
export const RoleManagement: React.FC = () => {
  const [searchParams, setSearchParams] = useState<RoleListReq>({
    page: 1,
    page_size: 10,
  });
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isPermissionDialogOpen, setIsPermissionDialogOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);

  const queryClient = useQueryClient();

  // 获取角色列表
  const { data: rolesData, isLoading } = useQuery({
    queryKey: ['roles', searchParams],
    queryFn: () => roleService.getRoles(searchParams),
  });

  // 创建角色
  const createMutation = useMutation({
    mutationFn: roleService.createRole,
    onSuccess: () => {
      toast.success('角色创建成功');
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      setIsCreateDialogOpen(false);
    },
    onError: (error: any) => {
      toast.error(error.message || '角色创建失败');
    },
  });

  // 更新角色
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateRoleReq }) =>
      roleService.updateRole(id, data),
    onSuccess: () => {
      toast.success('角色更新成功');
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      setIsEditDialogOpen(false);
      setSelectedRole(null);
    },
    onError: (error: any) => {
      toast.error(error.message || '角色更新失败');
    },
  });

  // 删除角色
  const deleteMutation = useMutation({
    mutationFn: roleService.deleteRole,
    onSuccess: () => {
      toast.success('角色删除成功');
      queryClient.invalidateQueries({ queryKey: ['roles'] });
    },
    onError: (error: any) => {
      toast.error(error.message || '角色删除失败');
    },
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ ...searchParams, page: 1 });
  };

  const handleEdit = (role: Role) => {
    setSelectedRole(role);
    setIsEditDialogOpen(true);
  };

  const handleDelete = (role: Role) => {
    if (window.confirm(`确定要删除角色 "${role.name}" 吗？`)) {
      deleteMutation.mutate(role.id);
    }
  };

  const handleAssignPermissions = (role: Role) => {
    setSelectedRole(role);
    setIsPermissionDialogOpen(true);
  };

  const roles = rolesData?.data?.list || [];
  const total = rolesData?.data?.total || 0;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">角色管理</h1>
        <PermissionWrapper resource="role" action="create">
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                创建角色
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>创建角色</DialogTitle>
                <DialogDescription>
                  创建新的系统角色
                </DialogDescription>
              </DialogHeader>
              <RoleForm
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
              <Label htmlFor="search-name">角色名称</Label>
              <Input
                id="search-name"
                placeholder="请输入角色名称"
                value={searchParams.name || ''}
                onChange={(e) => setSearchParams({ ...searchParams, name: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-code">角色代码</Label>
              <Input
                id="search-code"
                placeholder="请输入角色代码"
                value={searchParams.code || ''}
                onChange={(e) => setSearchParams({ ...searchParams, code: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-status">状态</Label>
              <Select
                value={searchParams.status || ''}
                onValueChange={(value) => setSearchParams({ ...searchParams, status: value as any })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="选择状态" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">全部</SelectItem>
                  <SelectItem value="active">激活</SelectItem>
                  <SelectItem value="inactive">禁用</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <Button type="submit">
              <Search className="w-4 h-4 mr-2" />
              搜索
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* 角色列表 */}
      <Card>
        <CardHeader>
          <CardTitle>角色列表 (共 {total} 条)</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>角色名称</TableHead>
                <TableHead>角色代码</TableHead>
                <TableHead>描述</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>创建时间</TableHead>
                <TableHead className="text-right">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8">
                    加载中...
                  </TableCell>
                </TableRow>
              ) : roles.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8">
                    暂无数据
                  </TableCell>
                </TableRow>
              ) : (
                roles.map((role) => (
                  <TableRow key={role.id}>
                    <TableCell className="font-medium">{role.name}</TableCell>
                    <TableCell>
                      <code className="bg-gray-100 px-2 py-1 rounded text-sm">
                        {role.code}
                      </code>
                    </TableCell>
                    <TableCell>{role.description || '-'}</TableCell>
                    <TableCell>
                      <Badge variant={role.status === 'active' ? 'default' : 'secondary'}>
                        {role.status === 'active' ? '激活' : '禁用'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {role.created_at ? new Date(role.created_at).toLocaleDateString() : '-'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex gap-2 justify-end">
                        <PermissionWrapper resource="role" action="update">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEdit(role)}
                          >
                            <Edit className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                        
                        <PermissionWrapper resource="role" action="assign_permission">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleAssignPermissions(role)}
                          >
                            <Shield className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                        
                        <PermissionWrapper resource="role" action="delete">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDelete(role)}
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

      {/* 编辑角色对话框 */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑角色</DialogTitle>
            <DialogDescription>
              修改角色信息
            </DialogDescription>
          </DialogHeader>
          {selectedRole && (
            <RoleForm
              role={selectedRole}
              onSubmit={(data) => updateMutation.mutate({ id: selectedRole.id, data })}
              isLoading={updateMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* 权限分配对话框 */}
      <Dialog open={isPermissionDialogOpen} onOpenChange={setIsPermissionDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>分配权限</DialogTitle>
            <DialogDescription>
              为角色分配系统权限
            </DialogDescription>
          </DialogHeader>
          {selectedRole && (
            <PermissionAssign
              role={selectedRole}
              onClose={() => {
                setIsPermissionDialogOpen(false);
                setSelectedRole(null);
              }}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default RoleManagement;
