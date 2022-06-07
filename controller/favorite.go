package controller

import (
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FavoriteQuery struct {
	UserId     int64  `form:"user_id" binding:"required"`
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
	var favoriteQuery FavoriteQuery
	err := c.ShouldBindQuery(&favoriteQuery)
	if err != nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "参数不合法"})
		return
	}
	if user := db.Redis.Get(favoriteQuery.Token); user == nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "用户未登录！"})
		return
	}
	if favoriteQuery.ActionType == 1 {
		favoriteService.DoLike(favoriteQuery.UserId, favoriteQuery.VideoId)
	} else if favoriteQuery.ActionType == 2 {
		favoriteService.CancelLike(favoriteQuery.UserId, favoriteQuery.VideoId)
	} else {
		c.JSON(http.StatusBadRequest, model.Response{StatusCode: 400, StatusMsg: "未知错误"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var favoriteListQuery FavoriteListQuery
	err := c.ShouldBindQuery(&favoriteListQuery)
	if err != nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "参数不合法"})
		return
	}
	if user := db.Redis.Get(favoriteListQuery.Token); user == nil {
		c.JSON(http.StatusForbidden, model.Response{StatusCode: 403, StatusMsg: "用户未登录！"})
		return
	}
	likeVideoList := favoriteService.GetLikeList(favoriteListQuery.UserId)
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: model.Response{
			StatusCode: 200,
		},
		FavoriteList: likeVideoList,
	})
}
