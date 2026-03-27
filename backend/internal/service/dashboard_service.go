package service

import (
	"fmt"
	"time"

	"go_sleep_admin/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

type DashboardMetric struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Value string `json:"value"`
	Tip   string `json:"tip"`
	Color string `json:"color"`
	Trend string `json:"trend"`
}

type DashboardShortcut struct {
	Path       string `json:"path"`
	Title      string `json:"title"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Tip        string `json:"tip"`
}

type DashboardShortcutGroup struct {
	Key   string              `json:"key"`
	Title string              `json:"title"`
	Icon  string              `json:"icon"`
	Items []DashboardShortcut `json:"items"`
}

type DashboardTrendItem struct {
	Month string `json:"month"`
	Total int64  `json:"total"`
}

type DashboardPieItem struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

type DashboardFinanceItem struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Color string `json:"color"`
}

type DashboardPayload struct {
	ShortcutGroups []DashboardShortcutGroup `json:"shortcut_groups"`
	Metrics        []DashboardMetric        `json:"metrics"`
	Finance        []DashboardFinanceItem   `json:"finance"`
	OrderTrend     []DashboardTrendItem     `json:"order_trend"`
	OrderPie       []DashboardPieItem       `json:"order_pie"`
}

func (s *DashboardService) Overview() (*DashboardPayload, error) {
	articleCount, err := s.repo.CountArticles()
	if err != nil {
		return nil, err
	}
	menuCount, err := s.repo.CountMenus()
	if err != nil {
		return nil, err
	}
	uploadCount, err := s.repo.CountUploads()
	if err != nil {
		return nil, err
	}
	adminCount, err := s.repo.CountAdmins()
	if err != nil {
		return nil, err
	}
	roleCount, err := s.repo.CountRoles()
	if err != nil {
		return nil, err
	}
	onlineCount, err := s.repo.CountOnlineSessions(time.Now())
	if err != nil {
		return nil, err
	}
	finance, err := s.repo.Finance()
	if err != nil {
		return nil, err
	}
	trend, err := s.repo.OrderTrend(12)
	if err != nil {
		return nil, err
	}
	pie, err := s.repo.OrderChannelBreakdown(8)
	if err != nil {
		return nil, err
	}

	payload := &DashboardPayload{
		ShortcutGroups: []DashboardShortcutGroup{
			{
				Key:   "content",
				Title: "常用功能",
				Icon:  "folder-menu",
				Items: []DashboardShortcut{
					{Path: "/articles", Title: "文章管理", Permission: "/api/v1/articles#GET", Icon: "form", Tip: "内容列表与编辑"},
					{Path: "/articles/new", Title: "新增文章", Permission: "/api/v1/articles#POST", Icon: "add-voucher", Tip: "快速创建新内容"},
					{Path: "/channels", Title: "渠道管理", Permission: "/api/v1/channels#GET", Icon: "list", Tip: "维护内容渠道和发布来源"},
					{Path: "/orders", Title: "订单管理", Permission: "/api/v1/orders#GET", Icon: "list", Tip: "查看支付订单、导出和详情核对"},
					{Path: "/uploads", Title: "资源中心", Permission: "/api/v1/uploads#GET", Icon: "folder-open", Tip: "管理图片音视频资源"},
					{Path: "/content-menus", Title: "内容菜单", Permission: "/api/v1/menus/tree#GET", Icon: "classify", Tip: "维护三级菜单结构"},
				},
			},
			{
				Key:   "system",
				Title: "系统管理",
				Icon:  "set",
				Items: []DashboardShortcut{
					{Path: "/system/admins", Title: "管理员", Permission: "/api/v1/admins#GET", Icon: "user", Tip: "账号、部门、岗位归属"},
					{Path: "/system/roles", Title: "角色权限", Permission: "/api/v1/roles#GET", Icon: "permission", Tip: "角色与数据范围管理"},
					{Path: "/system/operation-logs", Title: "操作日志", Permission: "/api/v1/operation-logs#GET", Icon: "safety", Tip: "查看后台请求审计日志"},
					{Path: "/system/online-users", Title: "在线用户", Permission: "/api/v1/online-sessions#GET", Icon: "switch", Tip: "查看并强制下线在线会话"},
				},
			},
		},
		Metrics: []DashboardMetric{
			{Key: "articles", Title: "文章总数", Value: fmt.Sprintf("%d", articleCount), Tip: "内容资产规模", Color: "#165dff", Trend: "up"},
			{Key: "menus", Title: "内容菜单", Value: fmt.Sprintf("%d", menuCount), Tip: "三级菜单节点数量", Color: "#00b42a", Trend: "up"},
			{Key: "uploads", Title: "上传资源", Value: fmt.Sprintf("%d", uploadCount), Tip: "文件中心资源总量", Color: "#ff7d00", Trend: "up"},
			{Key: "admins", Title: "管理员", Value: fmt.Sprintf("%d", adminCount), Tip: "后台账号数量", Color: "#722ed1", Trend: "down"},
			{Key: "roles", Title: "角色数", Value: fmt.Sprintf("%d", roleCount), Tip: "权限角色总数", Color: "#eb0aa4", Trend: "up"},
			{Key: "online", Title: "在线会话", Value: fmt.Sprintf("%d", onlineCount), Tip: "当前在线用户", Color: "#0fc6c2", Trend: "up"},
		},
		Finance: []DashboardFinanceItem{
			{Title: "订单总额", Value: formatFen(finance.OrderAmount), Color: "#ff8625"},
			{Title: "退款金额", Value: formatFen(finance.RefundAmount), Color: "#165DFF"},
			{Title: "支付订单数", Value: fmt.Sprintf("%d", finance.PaidOrderCount), Color: "#39cbab"},
			{Title: "退款订单数", Value: fmt.Sprintf("%d", finance.RefundOrderCount), Color: "#6c73ff"},
		},
		OrderTrend: make([]DashboardTrendItem, 0, len(trend)),
		OrderPie:   make([]DashboardPieItem, 0, len(pie)),
	}

	for _, item := range trend {
		payload.OrderTrend = append(payload.OrderTrend, DashboardTrendItem{
			Month: item.Month,
			Total: item.Total,
		})
	}
	for _, item := range pie {
		payload.OrderPie = append(payload.OrderPie, DashboardPieItem{
			Type:  item.Channel,
			Value: item.Total,
		})
	}

	return payload, nil
}

func formatFen(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}
