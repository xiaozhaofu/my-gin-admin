package models

type AdminStatus uint8

const (
	AdminStatusNormal AdminStatus = iota + 1
	AdminStatusDisabled
)

type Admin struct {
	BaseID
	Username string      `gorm:"column:username;size:50;uniqueIndex;not null" json:"username"`
	Nickname string      `gorm:"column:nickname;size:50;not null;default:''" json:"nickname"`
	Avatar   string      `gorm:"column:avatar;size:255;not null;default:''" json:"avatar"`
	Phone    string      `gorm:"column:phone;size:20;index" json:"phone"`
	Email    string      `gorm:"column:email;size:100;index" json:"email"`
	Password string      `gorm:"column:password;size:255;not null" json:"-"`
	Status   AdminStatus `gorm:"column:status;not null;default:1" json:"status"`
	DeptID   *int64      `gorm:"column:dept_id;index" json:"dept_id"`
	PostID   *int64      `gorm:"column:post_id;index" json:"post_id"`
	Remark   string      `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
	Roles []Role `gorm:"many2many:admin_roles;" json:"roles,omitempty"`
}

func (Admin) TableName() string { return "admins" }
