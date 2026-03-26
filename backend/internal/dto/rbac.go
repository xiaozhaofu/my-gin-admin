package dto

type RoleDetail struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	IsBuiltIn   bool    `json:"is_built_in"`
	DataScope   int8    `json:"data_scope"`
	MenuIDs     []int64 `json:"menu_ids"`
}

type AdminDetail struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	Status   int      `json:"status"`
	RoleIDs  []int64  `json:"role_ids"`
	Roles    []string `json:"roles"`
	DeptID   *int64   `json:"dept_id"`
	PostID   *int64   `json:"post_id"`
}

type AdminMenuTree struct {
	ID         int64            `json:"id"`
	ParentID   *int64           `json:"parent_id"`
	Title      string           `json:"title"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	Component  string           `json:"component"`
	Icon       string           `json:"icon"`
	Permission string           `json:"permission"`
	Type       int8             `json:"type"`
	Sort       int              `json:"sort"`
	Hidden     bool             `json:"hidden"`
	KeepAlive  bool             `json:"keep_alive"`
	Method     string           `json:"method"`
	APIPath    string           `json:"api_path"`
	Children   []*AdminMenuTree `json:"children,omitempty"`
}

type AdminMenuUpsertRequest struct {
	ParentID   *int64 `json:"parent_id"`
	Title      string `json:"title" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Permission string `json:"permission"`
	Type       int8   `json:"type"`
	Sort       int    `json:"sort"`
	Hidden     bool   `json:"hidden"`
	KeepAlive  bool   `json:"keep_alive"`
	Method     string `json:"method"`
	APIPath    string `json:"api_path"`
}
