package main

import (
	"fmt"
	"go-crontab/bootstrap"
	"go-crontab/common"
	"go-crontab/config"
	"go-crontab/web/middleware"
	"go-crontab/web/routes"
	"net/http"
	"runtime"
	"time"
)

func initEnv()  {
	// 配置线程数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New("任务调度", "xiaolin")
	app.Bootstrap()
	app.Use(middleware.Cors())
	// 接入prometheus
	p := middleware.NewPrometheus("gin")
	//p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
	//	url := c.Request.URL.Path
	//	for _, p := range c.Params {
	//		if p.Key == "name" {
	//			url = strings.Replace(url, p.Value, ":name", 1)
	//			break
	//		}
	//	}
	//	return url
	//}
	p.Use(app.Engine);
	app.Configure(routes.ApiConfigure)

	return app
}

func main ()  {
	initEnv();
	app := newApp()

	startServer(app)
}

func startServer (b *bootstrap.Bootstrapper)  {
	server := &http.Server{
		Addr:           ":" + config.Cfg.Produce.Port,
		Handler:        b,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}




	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// 平滑退出，先结束所有在执行的任务
	common.GracefulExitWeb(server)
}
