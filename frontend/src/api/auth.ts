import { client, type ApiResponse } from "@/api/client";

export interface LoginPayload {
  username: string;
  password: string;
}

export interface UserProfile {
  id: number;
  username: string;
  roles: string[];
}

export interface LoginResult {
  access_token: string;
  access_token_expires: number;
  refresh_token: string;
  refresh_token_expires: number;
  user: {
    id: number;
    username: string;
    nickname: string;
    avatar: string;
    roles: string[];
    permissions: string[];
  };
}

export const loginAPI = (data: LoginPayload) => client.post<any, ApiResponse<LoginResult>>("/auth/login", data);
export const meAPI = () => client.get<any, ApiResponse<UserProfile>>("/auth/me");
