package middleWare

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Redis 注入redis客户端
var Redis = db.Redis

//注入权限service对象
var authService = service.NewAuthService()

// LoginCheck 登录校验
func LoginCheck(c *gin.Context) {
	fmt.Println("LoginCheck start......") //打印日志

	token := c.Query("token")

	// token不合法
	if token == "" || len(token) == 0 {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "请登录后再进行操作",
		})
		c.Abort()
		fmt.Println("token不合法 LoginCheck over......") //打印日志
		return
	}

	claims, err1 := utils.ParseToken(token)
	//token解析失败
	if err1 != nil || claims == nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "token不合法",
		})
		c.Abort()
		fmt.Println("token解析失败 LoginCheck over......") //打印日志
		return
	}

	value, err2 := Redis.Get(service.UserPrefix + token).Result()
	//redis中不存在相应的token信息
	if err2 != nil || value == "" {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "token失效,请重新登录",
		})
		c.Abort()
		fmt.Println("redis中不存在相应的token信息 token失效 LoginCheck over......") //打印日志
		return
	}

	fmt.Println("登录校验成功 LoginCheck over......") //打印日志

	//校验成功，执行后续请求
	c.Next()
}

// AuthorityCheck 权限校验（uri为对应的请求路由，不包含参数!）(token统一从 query 中获取！)
func AuthorityCheck(c *gin.Context) {
	fmt.Println("AuthorityCheck start......") //打印日志

	uri := strings.Split(c.Request.RequestURI, "?")[0] //截取请求uri
	fmt.Println("截取的uri为:", uri)                       //打印日志
	token := c.Query("token")
	claims, _ := utils.ParseToken(token) //解析token,获取当前登录用户id

	if flag := authService.IsPermit(claims.Id, uri); !flag { //没有权限执行此操作
		c.Abort()
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "权限不足",
		})
		fmt.Println("没有权限执行此操作 AuthorityCheck stop......") //打印日志
		return
	} else { //拥有权限，放行
		fmt.Println("拥有权限，放行 AuthorityCheck stop......") //打印日志
		c.Next()
	}
}
