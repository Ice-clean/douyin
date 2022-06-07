package controller

import (
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteQuery struct {
	Token      string `form:"token" binding:"required"`
	VideoId    int64  `form:"video_id" binding:"required"`
	ActionType int    `form:"actionType" binding:"required"`
}

type FavoriteListQuery struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type FavoriteListResponse struct {
	model.Response
	FavoriteList []model.Video `json:"video_list"`
}

var favoriteService = service.NewFavoriteService()

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, _ := strconv.Atoi(c.Query("action_type"))
	if user := db.Redis.Get(token); user == nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "用户未登录！"})
	}
	loginUser := service.NewUserService().FindUserByToken(token)
	if actionType == 1 {
		favoriteService.DoLike(loginUser.Id, videoId)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "点赞成功"})
	} else if actionType == 2 {
		favoriteService.CancelLike(loginUser.Id, videoId)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "取消点赞成功"})
	} else {
		c.JSON(http.StatusBadRequest, model.Response{StatusCode: 400, StatusMsg: "未知错误"})
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	loginUser := service.NewUserService().FindUserByToken(token)
	if user := db.Redis.Get(token); user == nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "用户未登录！"})
		return
	}
	likeVideoList := favoriteService.GetLikeList(loginUser.Id)
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: model.Response{
			StatusCode: 200,
		},
		FavoriteList: likeVideoList,
	})
}
