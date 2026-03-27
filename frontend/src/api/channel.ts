import { client, type ApiResponse } from "@/api/client";

export interface ChannelItem {
  id: number;
  name: string;
  code: string;
  status: number;
  remark: string;
}

export const channelListAPI = () => client.get<any, ApiResponse<ChannelItem[]>>("/channels");
export const channelSaveAPI = (id: number | null, data: Record<string, unknown>) =>
  id ? client.put<any, ApiResponse<{ id: number }>>(`/channels/${id}`, data) : client.post<any, ApiResponse<{ id: number }>>("/channels", data);
export const channelStatusAPI = (data: { ids: number[]; status: number }) =>
  client.put<any, ApiResponse<{ updated: number }>>("/channels/status", data);
