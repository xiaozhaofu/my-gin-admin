package models

import "time"

type BaseID struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
}

type BaseTimeField struct {
	CreatedAt time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"-"`
}

type BaseTimeFieldNoDelete struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}
