import { client, type ApiResponse } from "@/api/client";

export interface DashboardShortcutItem {
  path: string;
  title: string;
  permission: string;
  icon: string;
  tip: string;
}

export interface DashboardShortcutGroup {
  key: string;
  title: string;
  icon: string;
  items: DashboardShortcutItem[];
}

export interface DashboardMetric {
  key: string;
  title: string;
  value: string;
  tip: string;
  color: string;
  trend: "up" | "down";
}

export interface DashboardFinanceItem {
  title: string;
  value: string;
  color: string;
}

export interface DashboardTrendItem {
  month: string;
  total: number;
}

export interface DashboardPieItem {
  type: string;
  value: number;
}

export interface DashboardPayload {
  shortcut_groups: DashboardShortcutGroup[];
  metrics: DashboardMetric[];
  finance: DashboardFinanceItem[];
  order_trend: DashboardTrendItem[];
  order_pie: DashboardPieItem[];
}

export const dashboardAPI = () => client.get<any, ApiResponse<DashboardPayload>>("/dashboard");
