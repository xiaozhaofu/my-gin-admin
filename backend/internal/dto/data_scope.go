package dto

type DeptTree struct {
	ID       int64       `json:"id"`
	ParentID *int64      `json:"parent_id"`
	Name     string      `json:"name"`
	Code     string      `json:"code"`
	Status   int8        `json:"status"`
	Children []*DeptTree `json:"children,omitempty"`
}

type DeptUpsertRequest struct {
	ParentID *int64 `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

type PostUpsertRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}
