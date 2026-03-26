import { client, type ApiResponse } from "@/api/client";

export interface RoleItem {
  id: number;
  name: string;
  code: string;
  description: string;
  data_scope: number;
  menu_ids: number[];
}

export interface AdminItem {
  id: number;
  username: string;
  nickname: string;
  phone: string;
  email: string;
  status: number;
  role_ids: number[];
  roles: string[];
  dept_id?: number;
  post_id?: number;
}

export interface AdminMenuNode {
  id: number;
  parent_id?: number;
  title: string;
  name: string;
  path: string;
  component: string;
  icon: string;
  permission: string;
  type: number;
  sort: number;
  hidden: boolean;
  keep_alive: boolean;
  method: string;
  api_path: string;
  children?: AdminMenuNode[];
}

export interface DictTypeItem {
  id: number;
  type_id: number;
  label: string;
  value: string;
  sort: number;
  status: number;
  css_class: string;
  list_class: string;
  is_default: boolean;
  remark: string;
}

export interface DictType {
  id: number;
  name: string;
  type_code: string;
  status: number;
  remark: string;
  items: DictTypeItem[];
}

export interface SysConfigItem {
  id: number;
  config_name: string;
  config_key: string;
  config_value: string;
  config_type: number;
  remark: string;
}

export interface OperationLogItem {
  id: number;
  admin_id: number;
  username: string;
  method: string;
  path: string;
  status_code: number;
  success: boolean;
  client_ip: string;
  user_agent: string;
  request_body: string;
  duration_ms: number;
  error_message: string;
  created_at: string;
}

export interface LoginLogItem {
  id: number;
  username: string;
  success: boolean;
  ip: string;
  user_agent: string;
  failure_cause: string;
  created_at: string;
}

export interface OnlineSessionItem {
  id: number;
  admin_id: number;
  username: string;
  ip: string;
  user_agent: string;
  last_active_at: string;
  expired_at: string;
}

export interface SysJobItem {
  id: number;
  name: string;
  job_key: string;
  cron_expr: string;
  target: string;
  status: number;
  concurrent: boolean;
  remark: string;
}

export interface DeptNode {
  id: number;
  parent_id?: number;
  name: string;
  code: string;
  status: number;
  children?: DeptNode[];
}

export interface PostItem {
  id: number;
  name: string;
  code: string;
  sort: number;
  status: number;
  remark: string;
}

export const adminListAPI = () => client.get<any, ApiResponse<AdminItem[]>>("/admins");
export const adminSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/admins/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/admins", data);

export const roleListAPI = () => client.get<any, ApiResponse<RoleItem[]>>("/roles");
export const roleSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/roles/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/roles", data);

export const moduleMetadataAPI = () => client.get<any, ApiResponse<Array<{ slug: string; title: string; fields: string[] }>>>("/module-metadata");
export const adminMenuTreeAPI = () => client.get<any, ApiResponse<AdminMenuNode[]>>("/admin-menus/tree");
export const adminMenuSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/admin-menus/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/admin-menus", data);
export const adminMenuDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/admin-menus/${id}`);
export const dictTypeListAPI = () => client.get<any, ApiResponse<DictType[]>>("/dict-types");
export const dictTypeSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/dict-types/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/dict-types", data);
export const dictTypeDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/dict-types/${id}`);
export const dictItemSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/dict-items/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/dict-items", data);
export const dictItemDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/dict-items/${id}`);
export const sysConfigListAPI = () => client.get<any, ApiResponse<SysConfigItem[]>>("/sys-configs");
export const sysConfigSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/sys-configs/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/sys-configs", data);
export const sysConfigDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/sys-configs/${id}`);
export const operationLogListAPI = (params: Record<string, unknown>) =>
  client.get<any, ApiResponse<{ list: OperationLogItem[]; total: number; page: number; page_size: number }>>("/operation-logs", { params });
export const operationLogDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/operation-logs/${id}`);
export const operationLogClearAPI = () => client.delete<any, ApiResponse<{ cleared: boolean }>>("/operation-logs");
export const loginLogListAPI = () => client.get<any, ApiResponse<LoginLogItem[]>>("/login-logs");
export const onlineSessionListAPI = () => client.get<any, ApiResponse<OnlineSessionItem[]>>("/online-sessions");
export const onlineSessionForceOfflineAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/online-sessions/${id}`);
export const jobListAPI = () => client.get<any, ApiResponse<SysJobItem[]>>("/jobs");
export const jobSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/jobs/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/jobs", data);
export const jobDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/jobs/${id}`);
export const deptListAPI = () => client.get<any, ApiResponse<DeptNode[]>>("/depts/tree");
export const deptSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/depts/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/depts", data);
export const deptDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/depts/${id}`);
export const postListAPI = () => client.get<any, ApiResponse<PostItem[]>>("/posts");
export const postSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/posts/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/posts", data);
export const postDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/posts/${id}`);
