package models

import "time"

type LoginLog struct {
	BaseID
	Username     string `gorm:"column:username;size:100;index;not null;default:''" json:"username"`
	Success      bool   `gorm:"column:success;index;not null;default:false" json:"success"`
	IP           string `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	UserAgent    string `gorm:"column:user_agent;size:500;not null;default:''" json:"user_agent"`
	FailureCause string `gorm:"column:failure_cause;size:500;not null;default:''" json:"failure_cause"`
	BaseTimeField
}

func (LoginLog) TableName() string { return "login_logs" }

type OnlineSession struct {
	BaseID
	AdminID        int64      `gorm:"column:admin_id;index;not null" json:"admin_id"`
	Username       string     `gorm:"column:username;size:100;index;not null" json:"username"`
	AccessToken    string     `gorm:"column:access_token;size:512;uniqueIndex;not null" json:"access_token"`
	RefreshToken   string     `gorm:"column:refresh_token;size:512;not null;default:''" json:"refresh_token"`
	IP             string     `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	UserAgent      string     `gorm:"column:user_agent;size:500;not null;default:''" json:"user_agent"`
	LastActiveAt   time.Time  `gorm:"column:last_active_at;index;not null" json:"last_active_at"`
	ExpiredAt      time.Time  `gorm:"column:expired_at;index;not null" json:"expired_at"`
	ForceOfflineAt *time.Time `gorm:"column:force_offline_at" json:"force_offline_at"`
	BaseTimeField
}

func (OnlineSession) TableName() string { return "online_sessions" }
