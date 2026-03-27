package dto

import "time"

type OrderQuery struct {
	PageQuery
	TradeNo      string `form:"trade_no"`
	UserPhone    string `form:"user_phone"`
	OrderNo      string `form:"order_no"`
	PayChannel   string `form:"pay_channel"`
	PayMethod    *int   `form:"pay_method"`
	ChannelID    *int64 `form:"channel_id"`
	Status       *int   `form:"status"`
	RefundStatus *int   `form:"refund_status"`
	CreatedFrom  string `form:"created_from"`
	CreatedTo    string `form:"created_to"`
}

type OrderListItem struct {
	ID            int64      `json:"id"`
	OrderNo       string     `json:"order_no"`
	OrderToken    string     `json:"order_token"`
	UserID        int64      `json:"user_id"`
	UserNickname  string     `json:"user_nickname"`
	UserPhone     string     `json:"user_phone"`
	Status        int8       `json:"status"`
	ProductID     uint64     `json:"product_id"`
	ProductType   int8       `json:"product_type"`
	ProductTitle  string     `json:"product_title"`
	OriginalPrice int64      `json:"original_price"`
	DiscountPrice int64      `json:"discount_price"`
	PayAmount     int64      `json:"pay_amount"`
	PayMethod     int8       `json:"pay_method"`
	PayChannel    string     `json:"pay_channel"`
	TradeNo       string     `json:"trade_no"`
	PaidAt        *time.Time `json:"paid_at"`
	ChannelID     int64      `json:"channel_id"`
	ChannelName   string     `json:"channel_name"`
	ChannelCode   string     `json:"channel_code"`
	RefundStatus  int8       `json:"refund_status"`
	RefundAmount  int64      `json:"refund_amount"`
	CreatedAt     time.Time  `json:"created_at"`
}

type OrderDetailOrder struct {
	ID             int64      `json:"id"`
	OrderNo        string     `json:"order_no"`
	OrderToken     string     `json:"order_token"`
	UserID         int64      `json:"user_id"`
	Status         int8       `json:"status"`
	ProductID      uint64     `json:"product_id"`
	ProductType    int8       `json:"product_type"`
	ProductTitle   string     `json:"product_title"`
	OriginalPrice  int64      `json:"original_price"`
	DiscountPrice  int64      `json:"discount_price"`
	PayAmount      int64      `json:"pay_amount"`
	CouponID       *uint64    `json:"coupon_id"`
	CouponAmount   int64      `json:"coupon_amount"`
	PayMethod      int8       `json:"pay_method"`
	PayChannel     string     `json:"pay_channel"`
	TradeNo        string     `json:"trade_no"`
	PaidAt         *time.Time `json:"paid_at"`
	ChannelID      int64      `json:"channel_id"`
	ExpireAt       *time.Time `json:"expire_at"`
	AccessExpireAt *time.Time `json:"access_expire_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	RefundStatus   int8       `json:"refund_status"`
	RefundAmount   int64      `json:"refund_amount"`
	RefundNo       string     `json:"refund_no"`
	RefundAt       *time.Time `json:"refund_at"`
	RefundReason   string     `json:"refund_reason"`
	ClientIP       string     `json:"client_ip"`
	ClientIPRaw    string     `json:"client_ip_raw"`
	Remark         string     `json:"remark"`
	AdminRemark    string     `json:"admin_remark"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type OrderDetailUser struct {
	ID          int64      `json:"id"`
	Nickname    string     `json:"nickname"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Status      int8       `json:"status"`
	VipStatus   int8       `json:"vip_status"`
	VipLevel    int8       `json:"vip_level"`
	VipExpireAt *time.Time `json:"vip_expire_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type OrderDetailChannel struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Status    int8       `json:"status"`
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrderDetailBill struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	OrderNo        string     `json:"order_no"`
	TradeNo        string     `json:"trade_no"`
	ProductID      uint64     `json:"product_id"`
	ProductType    int8       `json:"product_type"`
	ProductTitle   string     `json:"product_title"`
	OriginalPrice  int64      `json:"original_price"`
	DiscountAmount int64      `json:"discount_amount"`
	PayAmount      int64      `json:"pay_amount"`
	PayMethod      int8       `json:"pay_method"`
	PayChannel     string     `json:"pay_channel"`
	ChannelID      int64      `json:"channel_id"`
	RefundAmount   int64      `json:"refund_amount"`
	RefundStatus   int8       `json:"refund_status"`
	PaidAt         time.Time  `json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type OrderDetail struct {
	Order   OrderDetailOrder    `json:"order"`
	User    *OrderDetailUser    `json:"user,omitempty"`
	Channel *OrderDetailChannel `json:"channel,omitempty"`
	Bill    *OrderDetailBill    `json:"bill,omitempty"`
}
