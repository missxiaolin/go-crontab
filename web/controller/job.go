package controller

import (
	"github.com/gin-gonic/gin"
	"go-crontab/common/etcd"
)

type Job struct {
	Base
}

/**
 * 保存任务接口 etcd
 * POST job={name: xxx, command: xxxx, cronExpr: xxx}
 */
func (t *Job) Save(c *gin.Context) {
	var (
		job etcd.Job
		oldJob *etcd.Job
		err error
	)
	if err = c.BindJSON(&job); err != nil {
		goto ERR
	}
	// 保存到etcd
	if oldJob, err = etcd.G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	t.Succ(c, oldJob)
	return
ERR:
	t.Err(c, "job保存失败", 500)
}
