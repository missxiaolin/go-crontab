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
func (t *Job) JobSave(c *gin.Context) {
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

/**
 * 删除任务接口
 */
func (t *Job) JobDel(c *gin.Context) {
	var (
		err error
		name string
		oldJob *etcd.Job
	)
	name = c.DefaultQuery("name", "")
	if name == "" {
		goto ERR
	}
	// 去删除任务
	if oldJob, err = etcd.G_jobMgr.DelJob(name); err != nil {
		goto ERR
	}
	t.Succ(c, oldJob)
ERR:
	t.Err(c, "job删除失败", 500)
}

/**
 * 任务接口列表
 */
func (t *Job) JobList(c gin.Context) {
	var (
		jobList []*etcd.Job
		err error
	)
	// 获取任务列表
	if jobList, err = etcd.G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	t.Succ(c, jobList)


ERR:
	t.Err(c, "job获取列表失败", 500)
}


