package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// 最终的视频列表
	var videoList []model.Video

	// 判断是否有登录
	token := c.Query("token")
	user := userService.FindUserByToken(token)
	if token != "" && user != nil {
		// 有登录的话使用推荐获得视频列表
		videoList = videoService.GetRecommend(user.Id, token)
	} else {
		// 没有登录，则从所有视频中获取视频列表
		videoList = videoService.GetVideoList(c.ClientIP())
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
