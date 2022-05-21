package controller

import (
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Redis = db.Redis

func LoginCheck(c *gin.Context) {
	token := c.Query("token")

	// token不合法
	if token == "" || len(token) == 0 {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "请登录",
		})
		c.Abort()
		return
	}

	claims, err := utils.ParseToken(token)
	//token解析失败
	if err != nil || claims == nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "token不合法",
		})
		c.Abort()
		return
	}

	value := Redis.Get("token")
	//redis中不存在相应的token信息
	if value == nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "token失效",
		})
		c.Abort()
		return
	}

	//校验成功，执行后续请求
	c.Next()
}
