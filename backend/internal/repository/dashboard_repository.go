package repository

import (
	"time"

	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) CountArticles() (int64, error) {
	var total int64
	err := r.db.Model(&models.Article{}).Count(&total).Error
	return total, err
}

func (r *DashboardRepository) CountMenus() (int64, error) {
	var total int64
	err := r.db.Model(&models.Menu{}).Count(&total).Error
	return total, err
}

func (r *DashboardRepository) CountUploads() (int64, error) {
	var total int64
	err := r.db.Model(&models.UploadFile{}).Count(&total).Error
	return total, err
}

func (r *DashboardRepository) CountAdmins() (int64, error) {
	var total int64
	err := r.db.Model(&models.Admin{}).Count(&total).Error
	return total, err
}

func (r *DashboardRepository) CountRoles() (int64, error) {
	var total int64
	err := r.db.Model(&models.Role{}).Count(&total).Error
	return total, err
}

func (r *DashboardRepository) CountOnlineSessions(now time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&models.OnlineSession{}).
		Where("expired_at > ?", now).
		Where("force_offline_at IS NULL").
		Count(&total).Error
	return total, err
}

type OrderTrendItem struct {
	Month string
	Total int64
}

func (r *DashboardRepository) OrderTrend(lastMonths int) ([]OrderTrendItem, error) {
	if lastMonths <= 0 {
		lastMonths = 12
	}

	start := time.Now().AddDate(0, -(lastMonths - 1), 0)
	type row struct {
		Month string
		Total int64
	}
	var rows []row
	err := r.db.Model(&models.OrderBill{}).
		Select("DATE_FORMAT(paid_at, '%Y-%m') as month, COALESCE(SUM(pay_amount), 0) as total").
		Where("paid_at >= ?", time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())).
		Group("DATE_FORMAT(paid_at, '%Y-%m')").
		Order("month asc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]int64, len(rows))
	for _, row := range rows {
		lookup[row.Month] = row.Total
	}

	result := make([]OrderTrendItem, 0, lastMonths)
	for i := lastMonths - 1; i >= 0; i-- {
		monthTime := time.Now().AddDate(0, -i, 0)
		monthKey := monthTime.Format("2006-01")
		result = append(result, OrderTrendItem{
			Month: monthTime.Format("1月"),
			Total: lookup[monthKey],
		})
	}

	return result, nil
}

type OrderChannelItem struct {
	Channel string
	Total   int64
}

func (r *DashboardRepository) OrderChannelBreakdown(limit int) ([]OrderChannelItem, error) {
	if limit <= 0 {
		limit = 8
	}

	type row struct {
		Channel string
		Total   int64
	}
	var rows []row
	err := r.db.Model(&models.OrderBill{}).
		Select("CASE WHEN pay_channel = '' THEN '未标记' ELSE pay_channel END as channel, COALESCE(SUM(pay_amount), 0) as total").
		Group("channel").
		Order("total desc").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]OrderChannelItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, OrderChannelItem{
			Channel: row.Channel,
			Total:   row.Total,
		})
	}
	return result, nil
}

type FinanceSummary struct {
	OrderAmount      int64
	RefundAmount     int64
	PaidOrderCount   int64
	RefundOrderCount int64
}

func (r *DashboardRepository) Finance() (*FinanceSummary, error) {
	type row struct {
		OrderAmount      int64
		RefundAmount     int64
		PaidOrderCount   int64
		RefundOrderCount int64
	}
	var data row
	err := r.db.Model(&models.OrderBill{}).
		Select(`
			COALESCE(SUM(pay_amount), 0) as order_amount,
			COALESCE(SUM(refund_amount), 0) as refund_amount,
			COUNT(*) as paid_order_count,
			COALESCE(SUM(CASE WHEN refund_status > 0 THEN 1 ELSE 0 END), 0) as refund_order_count
		`).
		Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return &FinanceSummary{
		OrderAmount:      data.OrderAmount,
		RefundAmount:     data.RefundAmount,
		PaidOrderCount:   data.PaidOrderCount,
		RefundOrderCount: data.RefundOrderCount,
	}, nil
}
