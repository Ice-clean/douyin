package result

import "C"
import (
	"errors"
	"github.com/RaymondCode/simple-demo/constant"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) AbortWithStatusJSON() {
	g.C.Abort()
	g.OperateResult(errors.New("服务器繁忙，请稍后再试"))
}

func (g *Gin) OperateResult(err error) {
	if err == nil {
		g.C.JSON(constant.SUCCESS, gin.H{
			"status_code": constant.SUCCESS,
			"status_msg":  constant.GetMsg(constant.SUCCESS),
		})
	} else {
		g.C.JSON(constant.ERROR, gin.H{
			"status_code": constant.ERROR,
			"status_msg":  err.Error(),
		})
	}
	return
}

func (g *Gin) AutoResult(err error, data interface{}) {
	if err == nil {
		g.C.JSON(constant.SUCCESS, gin.H{
			"status_code": constant.SUCCESS,
			"status_msg":  constant.GetMsg(constant.SUCCESS),
			"data":        data,
		})
	} else {
		g.C.JSON(constant.ERROR, gin.H{
			"status_code": constant.ERROR,
			"status_msg":  err.Error(),
			"data":        data,
		})
	}
	return
}

func (g *Gin) CommonResult(httpCode, code int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"status_code": httpCode,
		"status_msg":  constant.GetMsg(code),
		"data":        data,
	})
	return
}
