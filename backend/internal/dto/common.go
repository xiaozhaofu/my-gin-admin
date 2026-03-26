package dto

import "time"

type PageQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (q PageQuery) Normalize() (int, int) {
	page := max(q.Page, 1)
	pageSize := q.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}

type DateRangeQuery struct {
	CreatedFrom *time.Time `form:"created_from" time_format:"2006-01-02 15:04:05"`
	CreatedTo   *time.Time `form:"created_to" time_format:"2006-01-02 15:04:05"`
}

type IDPayload struct {
	ID int64 `json:"id" binding:"required"`
}

type BatchIDs struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}
