package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
//var usersLoginInfo = map[string]model.User{
//	"IceClean233333": {
//		Id:            1,
//		Name:          "IceClean",
//		FollowCount:   10,
//		FollowerCount: 5,
//		IsFollow:      true,
//	},
//}
//
//var userIdSequence = int64(1)

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

//注入service对象
var userService = service.NewUserService()

// Register 注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) <= 0 || len(username) > 32 || len(password) <= 0 || len(password) > 32 { //合法性校验
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名或密码不合法"},
		})
		return
	}

	if !userService.IsUsernameCanUse(username) { //用户名不可用（用户已存在）
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名已存在，换一个试试"},
		})
		return
	}

	//调用服务
	userId, token, err := userService.Register(username, password)
	if err != nil { //调用服务层出错
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "系统繁忙"},
		})
	} else { //成功注册并返回
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   userId,
			Token:    token,
		})
	}

	//token := username + password
	//
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := model.User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

// Login 登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//fmt.Println(strings.Split(c.Request.RequestURI, "?")[0]) //请求uri

	if len(username) <= 0 || len(username) > 32 || len(password) <= 0 || len(password) > 32 { //合法性校验
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名或密码不合法"},
		})
		return
	}

	userId, token, err := userService.Login(username, password)
	if err != nil { //用户名或密码错误
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名或密码错误，请重新输入"},
		})
	} else { //成功登录并返回
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   userId,
			Token:    token,
		})
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

// UserInfo 查看用户信息
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")

	user, err := userService.UserInfo(userId, token)
	if err != nil || user == nil { //查询不到用户
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	} else { //成功查询并返回
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     *user,
		})
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: model.Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

// Logout 退出登录
func Logout(c *gin.Context) {
	userService.Logout(c.Query("token"))
}
