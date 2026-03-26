package models

type Dept struct {
	BaseID
	ParentID *int64 `gorm:"column:parent_id;index" json:"parent_id"`
	Name     string `gorm:"column:name;size:100;not null" json:"name"`
	Code     string `gorm:"column:code;size:100;uniqueIndex;not null" json:"code"`
	Leader   string `gorm:"column:leader;size:100;not null;default:''" json:"leader"`
	Phone    string `gorm:"column:phone;size:30;not null;default:''" json:"phone"`
	Email    string `gorm:"column:email;size:100;not null;default:''" json:"email"`
	Sort     int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Status   int8   `gorm:"column:status;not null;default:1" json:"status"`
	BaseTimeField
}

func (Dept) TableName() string { return "depts" }

type Post struct {
	BaseID
	Name   string `gorm:"column:name;size:100;not null" json:"name"`
	Code   string `gorm:"column:code;size:100;uniqueIndex;not null" json:"code"`
	Sort   int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Status int8   `gorm:"column:status;not null;default:1" json:"status"`
	Remark string `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
}

func (Post) TableName() string { return "posts" }
