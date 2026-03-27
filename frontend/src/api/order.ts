import { client, type ApiResponse } from "@/api/client";

export interface OrderListItem {
  id: number;
  order_no: string;
  order_token: string;
  user_id: number;
  user_nickname: string;
  user_phone: string;
  status: number;
  product_id: number;
  product_type: number;
  product_title: string;
  original_price: number;
  discount_price: number;
  pay_amount: number;
  pay_method: number;
  pay_channel: string;
  trade_no: string;
  paid_at?: string | null;
  channel_id: number;
  channel_name: string;
  channel_code: string;
  refund_status: number;
  refund_amount: number;
  created_at: string;
}

export interface OrderListResult {
  list: OrderListItem[];
  total: number;
  page: number;
  page_size: number;
}

export interface OrderDetailOrder {
  id: number;
  order_no: string;
  order_token: string;
  user_id: number;
  status: number;
  product_id: number;
  product_type: number;
  product_title: string;
  original_price: number;
  discount_price: number;
  pay_amount: number;
  coupon_id?: number | null;
  coupon_amount: number;
  pay_method: number;
  pay_channel: string;
  trade_no: string;
  paid_at?: string | null;
  channel_id: number;
  expire_at?: string | null;
  access_expire_at?: string | null;
  delivered_at?: string | null;
  refund_status: number;
  refund_amount: number;
  refund_no: string;
  refund_at?: string | null;
  refund_reason: string;
  client_ip: string;
  client_ip_raw: string;
  remark: string;
  admin_remark: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface OrderDetailUser {
  id: number;
  nickname: string;
  phone: string;
  email: string;
  status: number;
  vip_status: number;
  vip_level: number;
  vip_expire_at?: string | null;
  created_at: string;
}

export interface OrderDetailChannel {
  id: number;
  name: string;
  code: string;
  status: number;
  remark: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface OrderDetailBill {
  id: number;
  user_id: number;
  order_no: string;
  trade_no: string;
  product_id: number;
  product_type: number;
  product_title: string;
  original_price: number;
  discount_amount: number;
  pay_amount: number;
  pay_method: number;
  pay_channel: string;
  channel_id: number;
  refund_amount: number;
  refund_status: number;
  paid_at: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface OrderDetail {
  order: OrderDetailOrder;
  user?: OrderDetailUser;
  channel?: OrderDetailChannel;
  bill?: OrderDetailBill;
}

export const orderListAPI = (params: Record<string, unknown>) =>
  client.get<any, ApiResponse<OrderListResult>>("/orders", { params });

export const orderDetailAPI = (id: number | string) =>
  client.get<any, ApiResponse<OrderDetail>>(`/orders/${id}`);

export const orderExportAPI = (params: Record<string, unknown>) =>
  client.get<any, Blob>("/orders/export", { params, responseType: "blob" });
