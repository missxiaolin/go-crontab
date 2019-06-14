package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time // expr.Next(now)
}

func main()  {
	//task()

	manyTask()
}

// 多任务
func manyTask()  {
	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*CronJob // key: 任务名字
	)

	scheduleTable = make(map[string]*CronJob)

	// 当前时间
	now = time.Now()

	// 定义2个cronJob
	expr = cronexpr.MustParse("*/5 * * * *")
	cronJob = &CronJob{
		expr: expr,
		nextTime: expr.Next(now),
	}

	scheduleTable["jobOne"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * *")
	cronJob = &CronJob{
		expr: expr,
		nextTime: expr.Next(now),
	}

	scheduleTable["jobTwo"] = cronJob

	go func() {
		var (
			jobName string
			cronJob *CronJob
			now time.Time
		)
		// 定时检查任务调度表
		for  {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，执行任务
					go func(jobName string) {
						fmt.Printf("执行:", jobName)
					}(jobName)

					// 计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Printf(jobName, "下次执行时间:", cronJob.nextTime)
				}
			}
			select {
			case <- time.NewTimer(100 * time.Millisecond).C:
				
			}
		}
	}()

	time.Sleep(100000 * time.Millisecond)

}

// 单任务
func task()  {
	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)

	// 每分钟执行一次
	if expr, err = cronexpr.Parse("*/5 * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	nextTime = expr.Next(now)

	fmt.Println(nextTime)

	expr = expr
}
