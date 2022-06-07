package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

type FollowListResponse struct {
	model.Response
	FollowList []model.User `json:"user_list"`
}

var followService = service.NewFollowService()

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	var userId, _ = strconv.ParseInt(c.Query("user_id"), 10, 64)
	var toUserId, _ = strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	var actionType, _ = strconv.ParseInt(c.Query("action_type"), 10, 64)
	if actionType == 1 {
		err := followService.DoFollow(userId, toUserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{StatusCode: 400, StatusMsg: "未知错误"})
		} else {
			c.JSON(http.StatusOK, model.Response{StatusCode: 200, StatusMsg: "关注成功"})
		}
	} else if actionType == 2 {
		err := followService.CancelFollow(userId, toUserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{StatusCode: 400, StatusMsg: "未知错误"})
		} else {
			c.JSON(http.StatusOK, model.Response{StatusCode: 200, StatusMsg: "取消关注成功"})
		}
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {

	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	followList := followService.GetFollowList(userId)
	c.JSON(http.StatusOK, FollowListResponse{
		Response: model.Response{
			StatusCode: 200,
		},
		FollowList: followList,
	})

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {

	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fmt.Println(c.Query("user_id"))
	followList := followService.GetFollowerList(userId)
	c.JSON(http.StatusOK, FollowListResponse{
		Response: model.Response{
			StatusCode: 200,
		},
		FollowList: followList,
	})
}
