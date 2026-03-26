package models

type DictType struct {
	BaseID
	Name     string `gorm:"column:name;size:100;not null" json:"name"`
	TypeCode string `gorm:"column:type_code;size:100;uniqueIndex;not null" json:"type_code"`
	Status   int8   `gorm:"column:status;not null;default:1" json:"status"`
	Remark   string `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
	Items []DictItem `gorm:"foreignKey:TypeID" json:"items,omitempty"`
}

func (DictType) TableName() string { return "dict_types" }

type DictItem struct {
	BaseID
	TypeID    int64  `gorm:"column:type_id;index;not null" json:"type_id"`
	Label     string `gorm:"column:label;size:100;not null" json:"label"`
	Value     string `gorm:"column:value;size:100;not null" json:"value"`
	Sort      int    `gorm:"column:sort;not null;default:0" json:"sort"`
	Status    int8   `gorm:"column:status;not null;default:1" json:"status"`
	CSSClass  string `gorm:"column:css_class;size:100;not null;default:''" json:"css_class"`
	ListClass string `gorm:"column:list_class;size:100;not null;default:''" json:"list_class"`
	IsDefault bool   `gorm:"column:is_default;not null;default:false" json:"is_default"`
	Remark    string `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
}

func (DictItem) TableName() string { return "dict_items" }

type SysConfig struct {
	BaseID
	ConfigName  string `gorm:"column:config_name;size:100;not null" json:"config_name"`
	ConfigKey   string `gorm:"column:config_key;size:100;uniqueIndex;not null" json:"config_key"`
	ConfigValue string `gorm:"column:config_value;size:1000;not null;default:''" json:"config_value"`
	ConfigType  int8   `gorm:"column:config_type;not null;default:0" json:"config_type"`
	Remark      string `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
}

func (SysConfig) TableName() string { return "sys_configs" }

type OperationLog struct {
	BaseID
	AdminID      int64  `gorm:"column:admin_id;index;not null;default:0" json:"admin_id"`
	Username     string `gorm:"column:username;size:100;not null;default:''" json:"username"`
	Method       string `gorm:"column:method;size:16;index;not null" json:"method"`
	Path         string `gorm:"column:path;size:255;index;not null" json:"path"`
	StatusCode   int    `gorm:"column:status_code;index;not null;default:200" json:"status_code"`
	Success      bool   `gorm:"column:success;index;not null" json:"success"`
	ClientIP     string `gorm:"column:client_ip;size:64;not null;default:''" json:"client_ip"`
	UserAgent    string `gorm:"column:user_agent;size:500;not null;default:''" json:"user_agent"`
	RequestBody  string `gorm:"column:request_body;type:text" json:"request_body"`
	DurationMS   int64  `gorm:"column:duration_ms;not null;default:0" json:"duration_ms"`
	ErrorMessage string `gorm:"column:error_message;size:1000;not null;default:''" json:"error_message"`
	BaseTimeField
}

func (OperationLog) TableName() string { return "operation_logs" }
