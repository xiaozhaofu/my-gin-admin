package dto

type DictTypeUpsertRequest struct {
	Name     string `json:"name" binding:"required"`
	TypeCode string `json:"type_code" binding:"required"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
}

type DictItemUpsertRequest struct {
	TypeID    int64  `json:"type_id" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Sort      int    `json:"sort"`
	Status    int8   `json:"status"`
	CSSClass  string `json:"css_class"`
	ListClass string `json:"list_class"`
	IsDefault bool   `json:"is_default"`
	Remark    string `json:"remark"`
}

type SysConfigUpsertRequest struct {
	ConfigName  string `json:"config_name" binding:"required"`
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value"`
	ConfigType  int8   `json:"config_type"`
	Remark      string `json:"remark"`
}
