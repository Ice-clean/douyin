package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 注入服务
var videoService = service.NewVideoService()

// Publish check token then save upload file to public directory
// 获取参数，调用 service 层处理
func Publish(c *gin.Context) {
	// 获取标题和 token
	token := c.PostForm("token")
	title := c.PostForm("title")
	if token == "" {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "没有 token"})
		return
	}
	if title == "" {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "没有 title"})
		return
	}

	// 解析用户
	userId := Redis.Get("douyin:" + token)
	fmt.Printf("用户 Id：%v", userId)
	//if userId == redis.Nil {
	//	c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "用户不存在"})
	//	return
	//}
	user, err := userService.UserInfo(userId.Val(), token)

	fmt.Printf("%v", user)
	// 读取视频
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	// 调用投稿服务
	finalName, err := videoService.PublishVideo(user, title, data)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	// 返回上传成功
	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "上传成功",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 获取用户 ID
	userId := c.Query("user_id")
	// 调用获取投稿服务
	id, _ := strconv.Atoi(userId)
	videoList := videoService.GetPublishList(int64(id))

	fmt.Println("进来到列表了")
	fmt.Printf("%+v", videoList)
	c.JSON(http.StatusOK, model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
