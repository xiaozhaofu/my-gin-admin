package dto

type OperationLogQuery struct {
	PageQuery
	Username string `form:"username"`
	Method   string `form:"method"`
	Path     string `form:"path"`
	Success  *bool  `form:"success"`
}
