package dto

type UploadListQuery struct {
	PageQuery
	OriginName string `form:"origin_name"`
	Type       *int8  `form:"type"`
	Scene      string `form:"scene"`
	Provider   string `form:"provider"`
}
