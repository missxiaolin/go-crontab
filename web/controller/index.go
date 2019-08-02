package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Index struct {
	Base
}

func (t *Index) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"ErrorCode": 0,
		"msg": "welcome",
		"data": nil,
	})
}