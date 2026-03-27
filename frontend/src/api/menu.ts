import { client, type ApiResponse } from "@/api/client";

export interface ContentMenu {
  id: number;
  name: string;
  parent_id?: number;
  level: number;
  path: string;
  page_path: string;
  icon: string;
  sort_order: number;
  is_active: boolean;
  children?: ContentMenu[];
}

export const menuTreeAPI = () => client.get<any, ApiResponse<ContentMenu[]>>("/menus/tree");
export const menuCascaderAPI = () => client.get<any, ApiResponse<ContentMenu[]>>("/menus/cascader");
export const menuCreateAPI = (data: Record<string, unknown>) => client.post<any, ApiResponse<{ id: number }>>("/menus", data);
export const menuUpdateAPI = (id: number, data: Record<string, unknown>) => client.put<any, ApiResponse<{ id: number }>>(`/menus/${id}`, data);
export const menuStatusAPI = (data: { ids: number[]; is_active: boolean }) => client.put<any, ApiResponse<{ updated: number }>>("/menus/status", data);
export const menuDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/menus/${id}`);
