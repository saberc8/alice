import React, { useMemo } from 'react';
import { useUserMenuTree } from '@/store/userStore';
import type { MenuTree } from '@/types/entity';
import { Link, useLocation } from 'react-router-dom';
import { cn } from '@/utils';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/ui/collapsible';
import { ChevronDown, ChevronRight } from 'lucide-react';

interface DynamicMenuProps {
  className?: string;
}

// 菜单项组件
interface MenuItemProps {
  menu: MenuTree;
  level?: number;
}

const MenuItem: React.FC<MenuItemProps> = ({ menu, level = 0 }) => {
  const location = useLocation();
  const [isOpen, setIsOpen] = React.useState(false);
  
  const hasChildren = menu.children && menu.children.length > 0;
  const isActive = menu.path && location.pathname === menu.path;
  const isParentActive = menu.path && location.pathname.startsWith(menu.path);

  // 菜单项样式
  const itemClasses = cn(
    'flex items-center gap-3 px-3 py-2 rounded-lg transition-colors',
    'hover:bg-gray-100 dark:hover:bg-gray-800',
    {
      'bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400': isActive,
      'text-gray-700 dark:text-gray-300': !isActive,
      'pl-6': level === 1,
      'pl-9': level === 2,
      'pl-12': level >= 3,
    }
  );

  // 根据菜单类型决定是否可点击
  const isClickable = menu.type === 2 && menu.path; // 只有菜单类型且有路径的才可点击

  // 渲染菜单图标
  const renderIcon = () => {
    if (menu.meta?.icon) {
      // 这里可以根据图标名称渲染对应的图标
      return <div className="w-4 h-4 text-gray-500" />;
    }
    return null;
  };

  // 如果有子菜单，渲染折叠菜单
  if (hasChildren) {
    return (
      <Collapsible open={isOpen || isParentActive} onOpenChange={setIsOpen}>
        <CollapsibleTrigger className={cn(itemClasses, 'w-full justify-between')}>
          <div className="flex items-center gap-3">
            {renderIcon()}
            <span className="font-medium">{menu.name}</span>
          </div>
          {isOpen || isParentActive ? (
            <ChevronDown className="w-4 h-4" />
          ) : (
            <ChevronRight className="w-4 h-4" />
          )}
        </CollapsibleTrigger>
        <CollapsibleContent className="space-y-1 mt-1">
          {menu.children?.map((child) => (
            <MenuItem key={child.id} menu={child} level={level + 1} />
          ))}
        </CollapsibleContent>
      </Collapsible>
    );
  }

  // 如果是可点击的菜单项
  if (isClickable) {
    return (
      <Link to={menu.path!} className={itemClasses}>
        {renderIcon()}
        <span>{menu.name}</span>
      </Link>
    );
  }

  // 如果是按钮类型（不显示在导航中）
  if (menu.type === 3) {
    return null;
  }

  // 其他类型（分组、目录等）
  return (
    <div className={cn(itemClasses, 'cursor-default')}>
      {renderIcon()}
      <span className="font-medium">{menu.name}</span>
    </div>
  );
};

/**
 * 动态菜单组件
 * 根据用户权限显示菜单
 */
export const DynamicMenu: React.FC<DynamicMenuProps> = ({ className }) => {
  const userMenuTree = useUserMenuTree();

  // 过滤掉按钮类型的菜单项
  const filteredMenuTree = useMemo(() => {
    const filterMenus = (menus: MenuTree[]): MenuTree[] => {
      return menus
        .filter(menu => menu.type !== 3) // 过滤掉按钮类型
        .map(menu => ({
          ...menu,
          children: menu.children ? filterMenus(menu.children) : undefined,
        }));
    };
    return filterMenus(userMenuTree);
  }, [userMenuTree]);

  if (!filteredMenuTree.length) {
    return (
      <div className={cn('p-4 text-center text-gray-500', className)}>
        暂无可用菜单
      </div>
    );
  }

  return (
    <nav className={cn('space-y-1', className)}>
      {filteredMenuTree.map((menu) => (
        <MenuItem key={menu.id} menu={menu} />
      ))}
    </nav>
  );
};

/**
 * 菜单面包屑组件
 * 根据当前路径显示面包屑导航
 */
export const MenuBreadcrumb: React.FC = () => {
  const location = useLocation();
  const userMenuTree = useUserMenuTree();

  // 查找当前路径对应的菜单项
  const findMenuByPath = (menus: MenuTree[], path: string): MenuTree[] => {
    for (const menu of menus) {
      if (menu.path === path) {
        return [menu];
      }
      if (menu.children) {
        const found = findMenuByPath(menu.children, path);
        if (found.length > 0) {
          return [menu, ...found];
        }
      }
    }
    return [];
  };

  const breadcrumbPath = findMenuByPath(userMenuTree, location.pathname);

  if (breadcrumbPath.length === 0) {
    return null;
  }

  return (
    <nav className="flex items-center space-x-2 text-sm text-gray-600">
      {breadcrumbPath.map((menu, index) => (
        <React.Fragment key={menu.id}>
          {index > 0 && <span>/</span>}
          <span className={index === breadcrumbPath.length - 1 ? 'font-medium text-gray-900' : ''}>
            {menu.name}
          </span>
        </React.Fragment>
      ))}
    </nav>
  );
};

export default DynamicMenu;
