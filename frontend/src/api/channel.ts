import { client, type ApiResponse } from "@/api/client";

export interface ChannelItem {
  id: number;
  name: string;
  code: string;
  status: number;
  remark: string;
}

export const channelListAPI = () => client.get<any, ApiResponse<ChannelItem[]>>("/channels");
