package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base struct{
	Ctx *gin.Context
}

// 抛错
func (t *Base) Err(c *gin.Context, msg string, ErrorCode uint) {

	fmt.Println("[Error]", msg)
	t.Api(c, http.StatusOK, gin.H{
		"success": false,
		"ErrorCode": ErrorCode,
		"msg": msg,
	})
}

// 成功
func (t *Base) Succ(c *gin.Context, data ...interface{}) {

	var d interface{}
	if data != nil && len(data) > 0 {
		d = data[0]
	}

	t.Api(c, http.StatusOK, gin.H{
		"success": true,
		"msg": "ok",
		"data": d,
	})
}


// 返回
func (t *Base) Api(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}