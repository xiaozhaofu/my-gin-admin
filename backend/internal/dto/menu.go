package dto

type MenuUpsertRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  *int64 `json:"parent_id"`
	SortOrder uint   `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
	PagePath  string `json:"page_path"`
	Icon      string `json:"icon"`
}

type MenuBatchStatusRequest struct {
	IDs      []int64 `json:"ids" binding:"required,min=1"`
	IsActive bool    `json:"is_active"`
}

type RoleUpsertRequest struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description string  `json:"description"`
	DataScope   int8    `json:"data_scope"`
	MenuIDs     []int64 `json:"menu_ids"`
}

type AdminUpsertRequest struct {
	Username string  `json:"username" binding:"required"`
	Nickname string  `json:"nickname" binding:"required"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Status   int     `json:"status"`
	DeptID   *int64  `json:"dept_id"`
	PostID   *int64  `json:"post_id"`
	RoleIDs  []int64 `json:"role_ids"`
}
