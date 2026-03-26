import { client, type ApiResponse } from "@/api/client";

export interface UploadRecord {
  id: number;
  origin_name: string;
  path: string;
  type: number;
  md5: string;
  scene: string;
  provider: string;
  created_at: string;
}

export const uploadFileAPI = (data: FormData) =>
  client.post<any, ApiResponse<{ id: string; url: string; name: string; ftype: number; md5: string; provider: string }>>("/uploads", data, {
    headers: { "Content-Type": "multipart/form-data" }
  });

export const uploadListAPI = (params?: Record<string, unknown>) => client.get<any, ApiResponse<{ list: UploadRecord[]; total: number; page: number; page_size: number }>>("/uploads", { params });
export const uploadBatchDeleteAPI = (ids: number[]) => client.delete<any, ApiResponse<{ deleted: number }>>("/uploads", { data: { ids } });
export const uploadDeleteAPI = (id: number) => client.delete<any, ApiResponse<{ id: number }>>(`/uploads/${id}`);
