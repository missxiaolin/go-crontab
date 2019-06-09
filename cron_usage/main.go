package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main()  {
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
