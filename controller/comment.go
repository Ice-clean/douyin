package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListQuery struct {
	Token   string `form:"token"`
	VideoId int64  `form:"video_id" binding:"required"`
}

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentResponse struct {
	model.Response
	model.Comment
}

//注入service对象
var commentService = service.NewCommentService()

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	var p service.CommentQuery
	err1 := c.ShouldBindQuery(&p)
	if err1 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}

	comment, err2 := commentService.CommentAction(&p) //调用
	if err2 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "系统异常"})
		return
	} else {
		c.JSON(http.StatusOK, CommentResponse{
			model.Response{
				StatusCode: 0,
				StatusMsg:  "操作成功",
			},
			*comment,
		})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	var p CommentListQuery
	err1 := c.ShouldBindQuery(&p)
	if err1 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "参数不合法"})
		return
	}

	commentList, err2 := commentService.CommentList(p.Token, p.VideoId)
	if err2 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "系统异常"})
		return
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			model.Response{
				StatusCode: 0,
				StatusMsg:  "加载成功",
			},
			*commentList,
		})
	}
}
