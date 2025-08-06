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
  UserCheck,
  Shield,
} from 'lucide-react';
import { toast } from 'sonner';

import userService, { 
  type CreateUserReq, 
  type UpdateUserReq,
  type AssignUserRolesReq,
  type UserListReq 
} from '@/api/services/userService';
import roleService from '@/api/services/roleService';
import type { User, Role } from '@/types/entity';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

// 用户表单组件
interface UserFormProps {
  user?: User;
  onSubmit: (data: CreateUserReq | UpdateUserReq) => void;
  isLoading?: boolean;
  isEdit?: boolean;
}

const UserForm: React.FC<UserFormProps> = ({ user, onSubmit, isLoading, isEdit = false }) => {
  const [formData, setFormData] = useState<CreateUserReq>({
    username: user?.username || '',
    email: user?.email || '',
    password: '',
    status: user?.status || 'active',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (isEdit) {
      // 编辑时，如果密码为空则不传递密码字段
      const updateData: UpdateUserReq = {
        username: formData.username,
        email: formData.email,
        status: formData.status,
      };
      if (formData.password) {
        updateData.password = formData.password;
      }
      onSubmit(updateData);
    } else {
      onSubmit(formData);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="username">用户名</Label>
        <Input
          id="username"
          value={formData.username}
          onChange={(e) => setFormData({ ...formData, username: e.target.value })}
          placeholder="请输入用户名"
          required
        />
      </div>
      
      <div className="space-y-2">
        <Label htmlFor="email">邮箱</Label>
        <Input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder="请输入邮箱"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">
          密码 {isEdit && <span className="text-sm text-gray-500">(留空则不修改)</span>}
        </Label>
        <Input
          id="password"
          type="password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          placeholder={isEdit ? "留空不修改密码" : "请输入密码"}
          required={!isEdit}
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

// 角色分配组件
interface RoleAssignProps {
  user: User;
  onClose: () => void;
}

const RoleAssign: React.FC<RoleAssignProps> = ({ user, onClose }) => {
  const [selectedRoles, setSelectedRoles] = useState<string[]>([]);
  const queryClient = useQueryClient();

  // 获取所有角色
  const { data: allRoles } = useQuery({
    queryKey: ['roles', 'all'],
    queryFn: () => roleService.getRoles({ page: 1, page_size: 1000 }),
  });

  // 获取用户当前角色
  const { data: userRoles } = useQuery({
    queryKey: ['user-roles', user.id],
    queryFn: () => userService.getUserRoles(user.id),
  });

  // 分配角色
  const assignMutation = useMutation({
    mutationFn: (data: AssignUserRolesReq) => 
      userService.assignUserRoles(user.id, data),
    onSuccess: () => {
      toast.success('角色分配成功');
      queryClient.invalidateQueries({ queryKey: ['user-roles', user.id] });
      queryClient.invalidateQueries({ queryKey: ['users'] });
      onClose();
    },
    onError: (error: any) => {
      toast.error(error.message || '角色分配失败');
    },
  });

  useEffect(() => {
    if (userRoles?.data) {
      setSelectedRoles(userRoles.data.map(r => r.id));
    }
  }, [userRoles]);

  const handleToggleRole = (roleId: string) => {
    setSelectedRoles(prev => 
      prev.includes(roleId)
        ? prev.filter(id => id !== roleId)
        : [...prev, roleId]
    );
  };

  const handleSubmit = () => {
    assignMutation.mutate({ role_ids: selectedRoles });
  };

  const roles = allRoles?.data?.list || [];

  return (
    <div className="space-y-4">
      <div className="text-sm text-gray-600">
        为用户 "{user.username}" 分配角色
      </div>
      
      <div className="max-h-96 overflow-y-auto space-y-2">
        {roles.map((role) => (
          <div key={role.id} className="flex items-center space-x-2">
            <input
              type="checkbox"
              id={role.id}
              checked={selectedRoles.includes(role.id)}
              onChange={() => handleToggleRole(role.id)}
              className="rounded"
            />
            <label htmlFor={role.id} className="flex-1 cursor-pointer">
              <div className="font-medium">{role.name}</div>
              <div className="text-sm text-gray-500">
                {role.code} - {role.description}
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

// 用户管理主组件
export const UserManagement: React.FC = () => {
  const [searchParams, setSearchParams] = useState<UserListReq>({
    page: 1,
    page_size: 10,
  });
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isRoleDialogOpen, setIsRoleDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const queryClient = useQueryClient();

  // 获取用户列表
  const { data: usersData, isLoading } = useQuery({
    queryKey: ['users', searchParams],
    queryFn: () => userService.getUsers(searchParams),
  });

  // 创建用户
  const createMutation = useMutation({
    mutationFn: userService.createUser,
    onSuccess: () => {
      toast.success('用户创建成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setIsCreateDialogOpen(false);
    },
    onError: (error: any) => {
      toast.error(error.message || '用户创建失败');
    },
  });

  // 更新用户
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserReq }) =>
      userService.updateUser(id, data),
    onSuccess: () => {
      toast.success('用户更新成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setIsEditDialogOpen(false);
      setSelectedUser(null);
    },
    onError: (error: any) => {
      toast.error(error.message || '用户更新失败');
    },
  });

  // 删除用户
  const deleteMutation = useMutation({
    mutationFn: userService.deleteUser,
    onSuccess: () => {
      toast.success('用户删除成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
    onError: (error: any) => {
      toast.error(error.message || '用户删除失败');
    },
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ ...searchParams, page: 1 });
  };

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsEditDialogOpen(true);
  };

  const handleDelete = (user: User) => {
    if (window.confirm(`确定要删除用户 "${user.username}" 吗？`)) {
      deleteMutation.mutate(user.id);
    }
  };

  const handleAssignRoles = (user: User) => {
    setSelectedUser(user);
    setIsRoleDialogOpen(true);
  };

  const users = usersData?.data?.list || [];
  const total = usersData?.data?.total || 0;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">用户管理</h1>
        <PermissionWrapper resource="user" action="create">
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                创建用户
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>创建用户</DialogTitle>
                <DialogDescription>
                  创建新的系统用户
                </DialogDescription>
              </DialogHeader>
              <UserForm
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
              <Label htmlFor="search-username">用户名</Label>
              <Input
                id="search-username"
                placeholder="请输入用户名"
                value={searchParams.username || ''}
                onChange={(e) => setSearchParams({ ...searchParams, username: e.target.value })}
              />
            </div>
            <div className="flex-1">
              <Label htmlFor="search-email">邮箱</Label>
              <Input
                id="search-email"
                placeholder="请输入邮箱"
                value={searchParams.email || ''}
                onChange={(e) => setSearchParams({ ...searchParams, email: e.target.value })}
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

      {/* 用户列表 */}
      <Card>
        <CardHeader>
          <CardTitle>用户列表 (共 {total} 条)</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>用户名</TableHead>
                <TableHead>邮箱</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>创建时间</TableHead>
                <TableHead className="text-right">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8">
                    加载中...
                  </TableCell>
                </TableRow>
              ) : users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8">
                    暂无数据
                  </TableCell>
                </TableRow>
              ) : (
                users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.username}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Badge variant={user.status === 'active' ? 'default' : 'secondary'}>
                        {user.status === 'active' ? '激活' : '禁用'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {user.created_at ? new Date(user.created_at).toLocaleDateString() : '-'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex gap-2 justify-end">
                        <PermissionWrapper resource="user" action="update">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEdit(user)}
                          >
                            <Edit className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                        
                        <PermissionWrapper resource="user" action="assign_role">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleAssignRoles(user)}
                          >
                            <UserCheck className="w-4 h-4" />
                          </Button>
                        </PermissionWrapper>
                        
                        <PermissionWrapper resource="user" action="delete">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDelete(user)}
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

      {/* 编辑用户对话框 */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑用户</DialogTitle>
            <DialogDescription>
              修改用户信息
            </DialogDescription>
          </DialogHeader>
          {selectedUser && (
            <UserForm
              user={selectedUser}
              onSubmit={(data) => updateMutation.mutate({ id: selectedUser.id, data })}
              isLoading={updateMutation.isPending}
              isEdit
            />
          )}
        </DialogContent>
      </Dialog>

      {/* 角色分配对话框 */}
      <Dialog open={isRoleDialogOpen} onOpenChange={setIsRoleDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>分配角色</DialogTitle>
            <DialogDescription>
              为用户分配系统角色
            </DialogDescription>
          </DialogHeader>
          {selectedUser && (
            <RoleAssign
              user={selectedUser}
              onClose={() => {
                setIsRoleDialogOpen(false);
                setSelectedUser(null);
              }}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default UserManagement;
