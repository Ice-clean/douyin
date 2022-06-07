package controller

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"log"
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
	userString, err := db.Redis.Get(token).Result()
	fmt.Println("userString：", userString)
	if token != "" && err == nil && userString != "" {
		// 解析成 user 对象
		var user db.User
		err := json.Unmarshal([]byte(userString), &user)
		if err != nil {
			log.Fatal("json 解析失败：", err)
		}
		// 有登录的话使用推荐获得视频列表
		videoList = videoService.GetRecommend(user.Id, token)
	} else {
		// 没有登录，则从所有视频中获取视频列表
		fmt.Println("取得的IP：", c.ClientIP())
		videoList = videoService.GetVideoList(c.ClientIP())
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
