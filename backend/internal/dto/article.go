package dto

type ArticleQuery struct {
	PageQuery
	Title       string `form:"title"`
	Type        *int   `form:"type"`
	MenuID      *int64 `form:"menu_id"`
	ChannelID   *int64 `form:"channel_id"`
	IsPaid      *int   `form:"is_paid"`
	Status      *int   `form:"status"`
	CreatedFrom string `form:"created_from"`
	CreatedTo   string `form:"created_to"`
}

type ArticleUpsertRequest struct {
	Title       string  `json:"title" binding:"required"`
	Summary     string  `json:"summary"`
	Type        int     `json:"type" binding:"required"`
	Cover       string  `json:"cover"`
	CoverLarge  string  `json:"cover_large"`
	CoverMedium string  `json:"cover_medium"`
	CoverSmall  string  `json:"cover_small"`
	CoverType   string  `json:"cover_type"`
	MenuID      int64   `json:"menu_id"`
	MenuIDs     []int64 `json:"menu_ids"`
	ChannelID   int64   `json:"channel_id" binding:"required"`
	SortOrder   int8    `json:"sort_order"`
	IsPaid      int     `json:"is_paid"`
	IsTop       int     `json:"is_top"`
	IsHot       int     `json:"is_hot"`
	IsRecommend int     `json:"is_recommend"`
	Status      int     `json:"status"`
	Content     string  `json:"content" binding:"required"`
}

type ArticleBatchCreateItem struct {
	Title       string `json:"title" binding:"required"`
	Summary     string `json:"summary"`
	Cover       string `json:"cover"`
	CoverLarge  string `json:"cover_large"`
	CoverMedium string `json:"cover_medium"`
	CoverSmall  string `json:"cover_small"`
	CoverType   string `json:"cover_type"`
	Content     string `json:"content" binding:"required"`
}

type ArticleBatchCreateRequest struct {
	Type        int                      `json:"type" binding:"required"`
	Cover       string                   `json:"cover"`
	CoverLarge  string                   `json:"cover_large"`
	CoverMedium string                   `json:"cover_medium"`
	CoverSmall  string                   `json:"cover_small"`
	CoverType   string                   `json:"cover_type"`
	MenuID      int64                    `json:"menu_id" binding:"required"`
	MenuIDs     []int64                  `json:"menu_ids"`
	ChannelID   int64                    `json:"channel_id" binding:"required"`
	SortOrder   int8                     `json:"sort_order"`
	IsPaid      int                      `json:"is_paid"`
	IsTop       int                      `json:"is_top"`
	IsHot       int                      `json:"is_hot"`
	IsRecommend int                      `json:"is_recommend"`
	Status      int                      `json:"status"`
	Items       []ArticleBatchCreateItem `json:"items" binding:"required,min=1"`
}

type ArticleBatchStatusRequest struct {
	IDs    []int64 `json:"ids" binding:"required,min=1"`
	Status int     `json:"status" binding:"required"`
}
