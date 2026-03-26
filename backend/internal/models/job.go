package models

type SysJob struct {
	BaseID
	Name           string `gorm:"column:name;size:100;not null" json:"name"`
	Executor       string `gorm:"column:executor;size:100;not null;default:''" json:"executor"`
	JobKey         string `gorm:"column:job_key;size:100;uniqueIndex;not null" json:"job_key"`
	CronExpression string `gorm:"column:cron_expression;size:100;not null;default:''" json:"-"`
	CronExpr       string `gorm:"column:cron_expr;size:100;not null" json:"cron_expr"`
	Target         string `gorm:"column:target;size:255;not null" json:"target"`
	Status         int8   `gorm:"column:status;not null;default:1" json:"status"`
	Concurrent     bool   `gorm:"column:concurrent;not null;default:false" json:"concurrent"`
	Remark         string `gorm:"column:remark;size:255;not null;default:''" json:"remark"`
	BaseTimeField
}

func (SysJob) TableName() string { return "sys_jobs" }
