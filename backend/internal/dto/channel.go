package dto

type ChannelUpsertRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

type ChannelBatchStatusRequest struct {
	IDs    []int64 `json:"ids" binding:"required,min=1"`
	Status int     `json:"status" binding:"required"`
}
