package routes

import (
	"go-crontab/bootstrap"
	"go-crontab/web/controller"
)

func ApiConfigure(b *bootstrap.Bootstrapper)  {
	d := b.Group("/api")

	d.GET("/index", new(controller.Index).Welcome)
	d.POST("/job/save", new(controller.Job).Save)
}
