import { client, type ApiResponse } from "@/api/client";

export interface ArticleItem {
  id: number;
  title: string;
  summary: string;
  type: number;
  cover: string;
  cover_type: string;
  menu_id: number;
  channel_id: number;
  sort_order: number;
  is_paid: number;
  admin_id: number;
  is_top: number;
  is_hot: number;
  is_recommend: number;
  comment_num: number;
  share_num: number;
  view_num: number;
  collect_num: number;
  status: number;
  created_at: string;
  content?: { content: string };
}

export interface ArticleListResult {
  list: ArticleItem[];
  total: number;
  page: number;
  page_size: number;
}

export interface ArticleBatchCreateItem {
  title: string;
  summary: string;
  cover?: string;
  cover_type?: string;
  content: string;
}

export interface ArticleBatchCreatePayload {
  type: number;
  cover: string;
  cover_type: string;
  menu_id: number;
  channel_id: number;
  sort_order?: number;
  is_paid?: number;
  is_top?: number;
  is_hot?: number;
  is_recommend?: number;
  status: number;
  items: ArticleBatchCreateItem[];
}

export const articleListAPI = (params: Record<string, unknown>) => client.get<any, ApiResponse<ArticleListResult>>("/articles", { params });
export const articleDetailAPI = (id: string | number) => client.get<any, ApiResponse<ArticleItem>>(`/articles/${id}`);
export const articleCreateAPI = (data: Record<string, unknown>) => client.post<any, ApiResponse<{ id: number }>>("/articles", data);
export const articleBatchCreateAPI = (data: ArticleBatchCreatePayload) => client.post<any, ApiResponse<{ created: number; ids: number[] }>>("/articles/batch", data);
export const articleUpdateAPI = (id: string | number, data: Record<string, unknown>) => client.put<any, ApiResponse<{ id: number }>>(`/articles/${id}`, data);
export const articleStatusAPI = (data: { ids: number[]; status: number }) => client.put<any, ApiResponse<{ updated: number }>>("/articles/status", data);
export const articleDeleteAPI = (ids: number[]) => client.delete<any, ApiResponse<{ deleted: number }>>("/articles", { data: { ids } });
