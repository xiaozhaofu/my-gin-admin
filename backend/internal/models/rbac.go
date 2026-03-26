package models

type Role struct {
	BaseID
	Name        string `gorm:"column:name;size:50;uniqueIndex;not null" json:"name"`
	Code        string `gorm:"column:code;size:50;uniqueIndex;not null" json:"code"`
	Description string `gorm:"column:description;size:255;not null;default:''" json:"description"`
	IsBuiltIn   bool   `gorm:"column:is_built_in;not null;default:false" json:"is_built_in"`
	DataScope   int8   `gorm:"column:data_scope;not null;default:1" json:"data_scope"`
	BaseTimeField
}

func (Role) TableName() string { return "roles" }

type AdminRole struct {
	AdminID int64 `gorm:"column:admin_id;primaryKey" json:"admin_id"`
	RoleID  int64 `gorm:"column:role_id;primaryKey" json:"role_id"`
}

func (AdminRole) TableName() string { return "admin_roles" }

type AdminMenu struct {
	BaseID
	ParentID   *int64 `gorm:"column:parent_id;index" json:"parent_id"`
	Title      string `gorm:"column:title;size:100;not null" json:"title"`
	Name       string `gorm:"column:name;size:100;not null" json:"name"`
	Path       string `gorm:"column:path;size:255;not null" json:"path"`
	Component  string `gorm:"column:component;size:255;not null;default:''" json:"component"`
	Icon       string `gorm:"column:icon;size:100;not null;default:''" json:"icon"`
	Permission string `gorm:"column:permission;size:100;not null;default:''" json:"permission"`
	Type       int8   `gorm:"column:type;not null;default:1" json:"type"`
	Sort       int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Hidden     bool   `gorm:"column:hidden;not null;default:false" json:"hidden"`
	KeepAlive  bool   `gorm:"column:keep_alive;not null" json:"keep_alive"`
	Method     string `gorm:"column:method;size:20;not null;default:''" json:"method"`
	APIPath    string `gorm:"column:api_path;size:255;not null;default:''" json:"api_path"`
	BaseTimeField
}

func (AdminMenu) TableName() string { return "admin_menus" }

type RoleMenu struct {
	RoleID int64 `gorm:"column:role_id;primaryKey" json:"role_id"`
	MenuID int64 `gorm:"column:menu_id;primaryKey" json:"menu_id"`
}

func (RoleMenu) TableName() string { return "role_menus" }
