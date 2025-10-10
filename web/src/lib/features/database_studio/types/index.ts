export type ColumnType =
  | 'string'
  | 'integer'
  | 'float'
  | 'boolean'
  | 'date'
  | 'datetime'
  | 'timestamp'
  | 'json'
  | 'binary'
  | 'text'
  | 'uuid'
  | 'array';

export interface Column {
  name: string;
  type: ColumnType;
  nullable: boolean;
  default_value: string | null;
  primary_key: boolean;
  auto_increment: boolean;
  unique: boolean;
  comment: string;
}

export interface ForeignKey {
  column: string;
  referenced_table: string;
  referenced_column: string;
  on_delete: string;
  on_update: string;
}

export interface Index {
  name: string;
  columns: string[];
  unique: boolean;
  primary: boolean;
}

export interface TableSchema {
  name: string;
  columns: Column[];
  foreign_keys: ForeignKey[];
  indexes: Index[];
  row_count: number;
  comment: string;
}

export interface QueryResult {
  columns: string[];
  rows: Record<string, any>[];
  rows_affected: number;
  execution_time_ms: number;
}

export interface TableDataResult {
  columns: string[];
  rows: Record<string, any>[];
  total_rows: number;
  total_pages: number;
  page: number;
  page_size: number;
  limit: number;
  offset: number;
}

export type FilterOperator =
  | 'eq'
  | 'neq'
  | 'gt'
  | 'gte'
  | 'lt'
  | 'lte'
  | 'like'
  | 'not_like'
  | 'in'
  | 'not_in'
  | 'is_null'
  | 'is_not_null';

export interface Filter {
  column: string;
  operator: FilterOperator;
  value: any;
}

export interface Sort {
  column: string;
  direction: 'asc' | 'desc';
}

export interface TableDataOptions {
  page?: number;
  page_size?: number;
  limit?: number;
  offset?: number;
  filters?: Filter[];
  sorts?: Sort[];
}

export interface DatabaseInfo {
  version: string;
  size: number;
}

export interface ExecuteQueryRequest {
  query: string;
  limit?: number;
  offset?: number;
}

export interface InsertRowRequest {
  data: Record<string, any>;
}

export interface UpdateRowRequest {
  primary_key: Record<string, any>;
  data: Record<string, any>;
}

export interface DeleteRowRequest {
  primary_key: Record<string, any>;
}
