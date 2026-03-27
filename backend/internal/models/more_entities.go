package models

import "time"

type Membership struct {
	BaseID
	Name          string `gorm:"column:name;size:255;not null" json:"name"`
	VipType       int8   `gorm:"column:vip_type;not null;default:0" json:"vip_type"`
	VipDays       int64  `gorm:"column:vip_days;not null;default:0" json:"vip_days"`
	ChannelID     int64  `gorm:"column:channel_id;not null;index" json:"channel_id"`
	OriginalPrice int64  `gorm:"column:original_price;not null;default:0" json:"original_price"`
	Price         int64  `gorm:"column:price;not null;default:0" json:"price"`
	Status        int8   `gorm:"column:status;not null;default:0" json:"status"`
	Description   string `gorm:"column:description;size:255;not null;default:''" json:"description"`
	BaseTimeField
}

func (Membership) TableName() string { return "memberships" }

type User struct {
	BaseID
	Nickname    string     `gorm:"column:nickname;size:50;not null;default:''" json:"nickname"`
	Gender      int8       `gorm:"column:gender;not null;default:0" json:"gender"`
	Avatar      string     `gorm:"column:avatar;size:255;not null;default:''" json:"avatar"`
	Phone       string     `gorm:"column:phone;size:20;index" json:"phone"`
	Email       string     `gorm:"column:email;size:100;index" json:"email"`
	Password    string     `gorm:"column:password;size:255;not null;default:''" json:"-"`
	Status      int8       `gorm:"column:status;not null;default:1" json:"status"`
	VipStatus   int8       `gorm:"column:vip_status;not null;default:0" json:"vip_status"`
	VipLevel    int8       `gorm:"column:vip_level;not null;default:0" json:"vip_level"`
	VipExpireAt *time.Time `gorm:"column:vip_expire_at" json:"vip_expire_at"`
	Follow      int8       `gorm:"column:follow;not null;default:0" json:"follow"`
	BaseTimeField
}

func (User) TableName() string { return "users" }

type UserProfile struct {
	BaseID
	UserID        int64      `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	RegisterIP    string     `gorm:"column:register_ip;size:45;not null;default:''" json:"register_ip"`
	LastLoginIP   string     `gorm:"column:last_login_ip;size:45;not null;default:''" json:"last_login_ip"`
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"last_login_time"`
	DeviceSystem  string     `gorm:"column:device_system;size:50;not null;default:''" json:"device_system"`
	DeviceModel   string     `gorm:"column:device_model;size:100;not null;default:''" json:"device_model"`
	DeviceID      string     `gorm:"column:device_id;size:100;not null;default:''" json:"device_id"`
	Source        string     `gorm:"column:source;size:50;not null;default:'direct'" json:"source"`
	ChannelID     int64      `gorm:"column:channel_id;not null;index" json:"channel_id"`
	InviteCode    string     `gorm:"column:invite_code;size:20;not null;default:''" json:"invite_code"`
	InvitedBy     uint64     `gorm:"column:invited_by;not null;default:0" json:"invited_by"`
	Remark        string     `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	ExtraData     string     `gorm:"column:extra_data;type:json" json:"extra_data"`
	BaseTimeField
}

func (UserProfile) TableName() string { return "user_profiles" }

type UserArticleSave struct {
	BaseID
	UserID    int64 `gorm:"column:user_id;not null;index" json:"user_id"`
	ArticleID int64 `gorm:"column:article_id;not null;index" json:"article_id"`
	BaseTimeField
}

func (UserArticleSave) TableName() string { return "user_article_saves" }

type Order struct {
	BaseID
	OrderNo        string     `gorm:"column:order_no;size:64;uniqueIndex;not null" json:"order_no"`
	OrderToken     string     `gorm:"column:order_token;size:64;not null" json:"order_token"`
	UserID         int64      `gorm:"column:user_id;not null;index" json:"user_id"`
	Status         int8       `gorm:"column:status;not null;default:0;index" json:"status"`
	ProductID      uint64     `gorm:"column:product_id;not null" json:"product_id"`
	ProductType    int8       `gorm:"column:product_type;not null" json:"product_type"`
	ProductTitle   string     `gorm:"column:product_title;size:255;not null" json:"product_title"`
	OriginalPrice  int64      `gorm:"column:original_price;not null" json:"original_price"`
	DiscountPrice  int64      `gorm:"column:discount_price;not null;default:0" json:"discount_price"`
	PayAmount      int64      `gorm:"column:pay_amount;not null" json:"pay_amount"`
	CouponID       *uint64    `gorm:"column:coupon_id" json:"coupon_id"`
	CouponAmount   int64      `gorm:"column:coupon_amount;not null;default:0" json:"coupon_amount"`
	PayMethod      int8       `gorm:"column:pay_method;not null;default:0" json:"pay_method"`
	PayChannel     string     `gorm:"column:pay_channel;size:32;not null;default:''" json:"pay_channel"`
	TradeNo        string     `gorm:"column:trade_no;size:128;not null;default:''" json:"trade_no"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at"`
	ChannelID      int64      `gorm:"column:channel_id;not null;index" json:"channel_id"`
	ExpireAt       *time.Time `gorm:"column:expire_at" json:"expire_at"`
	AccessExpireAt *time.Time `gorm:"column:access_expire_at" json:"access_expire_at"`
	DeliveredAt    *time.Time `gorm:"column:delivered_at" json:"delivered_at"`
	RefundStatus   int8       `gorm:"column:refund_status;not null;default:0" json:"refund_status"`
	RefundAmount   int64      `gorm:"column:refund_amount;not null;default:0" json:"refund_amount"`
	RefundNo       string     `gorm:"column:refund_no;size:128;not null;default:''" json:"refund_no"`
	RefundAt       *time.Time `gorm:"column:refund_at" json:"refund_at"`
	RefundReason   string     `gorm:"column:refund_reason;size:512;not null;default:''" json:"refund_reason"`
	ClientIP       []byte     `gorm:"column:client_ip;type:varbinary(16)" json:"-"`
	ClientIPStr    string     `gorm:"column:client_ip_str;size:45;not null;default:''" json:"client_ip"`
	Remark         string     `gorm:"column:remark;size:512;not null;default:''" json:"remark"`
	AdminRemark    string     `gorm:"column:admin_remark;size:512;not null;default:''" json:"admin_remark"`
	BaseTimeField
}

func (Order) TableName() string { return "orders" }

type OrderBill struct {
	BaseID
	UserID         int64     `gorm:"column:user_id;not null;index" json:"user_id"`
	OrderNo        string    `gorm:"column:order_no;size:64;uniqueIndex;not null" json:"order_no"`
	TradeNo        string    `gorm:"column:trade_no;size:128;not null;default:''" json:"trade_no"`
	ProductID      uint64    `gorm:"column:product_id;not null" json:"product_id"`
	ProductType    int8      `gorm:"column:product_type;not null" json:"product_type"`
	ProductTitle   string    `gorm:"column:product_title;size:255;not null" json:"product_title"`
	OriginalPrice  int64     `gorm:"column:original_price;not null" json:"original_price"`
	DiscountAmount int64     `gorm:"column:discount_amount;not null;default:0" json:"discount_amount"`
	PayAmount      int64     `gorm:"column:pay_amount;not null" json:"pay_amount"`
	PayMethod      int8      `gorm:"column:pay_method;not null" json:"pay_method"`
	PayChannel     string    `gorm:"column:pay_channel;size:32;not null;default:''" json:"pay_channel"`
	ChannelID      int64     `gorm:"column:channel_id;not null;index" json:"channel_id"`
	RefundAmount   int64     `gorm:"column:refund_amount;not null;default:0" json:"refund_amount"`
	RefundStatus   int8      `gorm:"column:refund_status;not null;default:0" json:"refund_status"`
	PaidAt         time.Time `gorm:"column:paid_at;not null" json:"paid_at"`
	BaseTimeField
}

func (OrderBill) TableName() string { return "order_bills" }
