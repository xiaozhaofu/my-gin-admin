package repository

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) List(query dto.OrderQuery) ([]dto.OrderListItem, int64, error) {
	var (
		items []dto.OrderListItem
		total int64
	)

	page, pageSize := query.Normalize()
	base := r.buildListQuery(query)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base.
		Select(`
			orders.id,
			orders.order_no,
			orders.order_token,
			orders.user_id,
			orders.status,
			orders.product_id,
			orders.product_type,
			orders.product_title,
			orders.original_price,
			orders.discount_price,
			orders.pay_amount,
			orders.pay_method,
			orders.pay_channel,
			orders.trade_no,
			orders.paid_at,
			orders.channel_id,
			orders.refund_status,
			orders.refund_amount,
			orders.created_at,
			users.nickname as user_nickname,
			users.phone as user_phone,
			channels.name as channel_name,
			channels.code as channel_code
		`).
		Order("orders.id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *OrderRepository) ByID(id int64) (*dto.OrderDetail, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}

	result := &dto.OrderDetail{
		Order: dto.OrderDetailOrder{
			ID:             order.ID,
			OrderNo:        order.OrderNo,
			OrderToken:     order.OrderToken,
			UserID:         order.UserID,
			Status:         order.Status,
			ProductID:      order.ProductID,
			ProductType:    order.ProductType,
			ProductTitle:   order.ProductTitle,
			OriginalPrice:  order.OriginalPrice,
			DiscountPrice:  order.DiscountPrice,
			PayAmount:      order.PayAmount,
			CouponID:       order.CouponID,
			CouponAmount:   order.CouponAmount,
			PayMethod:      order.PayMethod,
			PayChannel:     order.PayChannel,
			TradeNo:        order.TradeNo,
			PaidAt:         order.PaidAt,
			ChannelID:      order.ChannelID,
			ExpireAt:       order.ExpireAt,
			AccessExpireAt: order.AccessExpireAt,
			DeliveredAt:    order.DeliveredAt,
			RefundStatus:   order.RefundStatus,
			RefundAmount:   order.RefundAmount,
			RefundNo:       order.RefundNo,
			RefundAt:       order.RefundAt,
			RefundReason:   order.RefundReason,
			ClientIP:       order.ClientIPStr,
			ClientIPRaw:    hex.EncodeToString(order.ClientIP),
			Remark:         order.Remark,
			AdminRemark:    order.AdminRemark,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
			DeletedAt:      order.DeletedAt,
		},
	}

	if order.UserID > 0 {
		var user models.User
		if err := r.db.First(&user, order.UserID).Error; err == nil {
			result.User = &dto.OrderDetailUser{
				ID:          user.ID,
				Nickname:    user.Nickname,
				Phone:       user.Phone,
				Email:       user.Email,
				Status:      user.Status,
				VipStatus:   user.VipStatus,
				VipLevel:    user.VipLevel,
				VipExpireAt: user.VipExpireAt,
				CreatedAt:   user.CreatedAt,
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if order.ChannelID > 0 {
		var channel models.Channel
		if err := r.db.First(&channel, order.ChannelID).Error; err == nil {
			result.Channel = &dto.OrderDetailChannel{
				ID:        channel.ID,
				Name:      channel.Name,
				Code:      channel.Code,
				Status:    int8(channel.Status),
				Remark:    channel.Remark,
				CreatedAt: channel.CreatedAt,
				UpdatedAt: channel.UpdatedAt,
				DeletedAt: channel.DeletedAt,
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	var bill models.OrderBill
	if err := r.db.Where("order_no = ?", order.OrderNo).First(&bill).Error; err == nil {
		result.Bill = &dto.OrderDetailBill{
			ID:             bill.ID,
			UserID:         bill.UserID,
			OrderNo:        bill.OrderNo,
			TradeNo:        bill.TradeNo,
			ProductID:      bill.ProductID,
			ProductType:    bill.ProductType,
			ProductTitle:   bill.ProductTitle,
			OriginalPrice:  bill.OriginalPrice,
			DiscountAmount: bill.DiscountAmount,
			PayAmount:      bill.PayAmount,
			PayMethod:      bill.PayMethod,
			PayChannel:     bill.PayChannel,
			ChannelID:      bill.ChannelID,
			RefundAmount:   bill.RefundAmount,
			RefundStatus:   bill.RefundStatus,
			PaidAt:         bill.PaidAt,
			CreatedAt:      bill.CreatedAt,
			UpdatedAt:      bill.UpdatedAt,
			DeletedAt:      bill.DeletedAt,
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return result, nil
}

func IsOrderNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *OrderRepository) buildListQuery(query dto.OrderQuery) *gorm.DB {
	db := r.db.Model(&models.Order{}).
		Joins("LEFT JOIN users ON users.id = orders.user_id").
		Joins("LEFT JOIN channels ON channels.id = orders.channel_id")

	orderNo := strings.TrimSpace(query.OrderNo)
	if orderNo != "" {
		db = db.Where("orders.order_no LIKE ?", "%"+orderNo+"%")
	}
	tradeNo := strings.TrimSpace(query.TradeNo)
	if tradeNo != "" {
		db = db.Where("orders.trade_no LIKE ?", "%"+tradeNo+"%")
	}
	userPhone := strings.TrimSpace(query.UserPhone)
	if userPhone != "" {
		db = db.Where("users.phone LIKE ?", "%"+userPhone+"%")
	}
	payChannel := strings.TrimSpace(query.PayChannel)
	if payChannel != "" {
		db = db.Where("orders.pay_channel = ?", payChannel)
	}
	if query.PayMethod != nil {
		db = db.Where("orders.pay_method = ?", *query.PayMethod)
	}
	if query.ChannelID != nil {
		db = db.Where("orders.channel_id = ?", *query.ChannelID)
	}
	if query.Status != nil {
		db = db.Where("orders.status = ?", *query.Status)
	}
	if query.RefundStatus != nil {
		db = db.Where("orders.refund_status = ?", *query.RefundStatus)
	}
	if t := parseOrderTime(query.CreatedFrom, false); !t.IsZero() {
		db = db.Where("orders.created_at >= ?", t)
	}
	if t := parseOrderTime(query.CreatedTo, true); !t.IsZero() {
		db = db.Where("orders.created_at <= ?", t)
	}
	return db
}

func parseOrderTime(value string, endOfDay bool) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}

	if t, err := time.ParseInLocation(time.DateTime, value, time.Local); err == nil {
		return t
	}
	if t, err := time.ParseInLocation(time.DateOnly, value, time.Local); err == nil {
		if endOfDay {
			return t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		return t
	}
	return time.Time{}
}
