export type FormWidgetType =
  | 'input'
  | 'textarea'
  | 'number'
  | 'select'
  | 'treeSelect'
  | 'switch';

export interface SelectOption {
  label: string;
  value: string | number | boolean | null;
  disabled?: boolean;
}

export interface TreeOption {
  label: string;
  value: string | number | null;
  disabled?: boolean;
  children?: TreeOption[];
}

export interface ProColumn<T = any> {
  title: string;
  dataIndex: string; // 支持 a.b.c 嵌套路径
  width?: number | string;
  ellipsis?: boolean;
  align?: 'left' | 'center' | 'right';
  // table 显示控制
  hideInTable?: boolean;
  // search 表单控制
  search?: boolean;
  searchPlaceholder?: string;
  // form 控制
  form?: boolean | {
    widget?: FormWidgetType;
    placeholder?: string;
    required?: boolean;
    options?: SelectOption[]; // select
    treeOptions?: TreeOption[]; // treeSelect
    props?: Record<string, any>;
  };
  // 渲染函数（可选）：由使用方通过插槽覆盖
  render?: (record: T) => any;
}

export interface ProTableRequestParams {
  // 通用的查询参数
  [key: string]: any;
}

export interface ProTableProps<T = any> {
  title?: string;
  rowKey?: string;
  columns: ProColumn<T>[];
  request: (params: ProTableRequestParams) => Promise<T[]>;
  // CRUD（存在即自动渲染按钮）
  onCreate?: (values: Partial<T>) => Promise<any>;
  onUpdate?: (id: string | number, values: Partial<T>) => Promise<any>;
  onDelete?: (id: string | number) => Promise<any>;
  // 是否树形数据
  isTree?: boolean;
}
