package etcd

// 定时任务
type Job struct {
	Name string `json:"name"` // 任务名
	Command string `json:"command"` // shell 命令
	CronExpr string `json:"CronExpr"` // cron 表达式
}

