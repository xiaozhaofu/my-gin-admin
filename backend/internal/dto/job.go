package dto

type SysJobUpsertRequest struct {
	Name       string `json:"name" binding:"required"`
	Executor   string `json:"executor"`
	JobKey     string `json:"job_key" binding:"required"`
	CronExpr   string `json:"cron_expr" binding:"required"`
	Target     string `json:"target" binding:"required"`
	Status     int8   `json:"status"`
	Concurrent bool   `json:"concurrent"`
	Remark     string `json:"remark"`
}
