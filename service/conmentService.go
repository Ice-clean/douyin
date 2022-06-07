package service

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"time"
)

type CommentQuery struct { //评论请求参数
	Token       string `form:"token" binding:"required"`
	VideoId     int64  `form:"video_id" binding:"required"`
	ActionType  int    `form:"action_type" binding:"required"`
	CommentText string `form:"comment_text"`
	CommentId   int64  `form:"comment_id"`
}

type CommentService struct {
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

//注入dao对象
var commentDao = db.NewCommentDao()

//注入userService对象
var userService = NewUserService()

// CommentAction 评论操作
func (c *CommentService) CommentAction(cq *CommentQuery) (commentRes *model.Comment, err error) {
	//解析token
	var claims *utils.Claims
	claims, err = utils.ParseToken(cq.Token)
	if err != nil { //解析失败，直接返回
		fmt.Println("token parse failed...")
		return nil, err
	}

	var comment *db.Comment //接受数据库返回对象

	if cq.ActionType == 1 { //发布评论
		comment, err = commentDao.CreateComment(cq.VideoId, claims.Id, cq.CommentText)
	} else if cq.ActionType == 2 { //删除评论
		comment = commentDao.DeleteComment(cq.CommentId)
	} else { //参数有误
		return nil, errors.New("参数有误")
	}

	if err != nil {
		return nil, err
	}

	var user *model.User
	user, err = userService.FindUserModelById(claims.Id, cq.Token)

	commentRes = &model.Comment{
		Id:         comment.Id,
		User:       *user,
		Content:    comment.CommentText,
		CreateDate: time.Unix(comment.CreateAt/1000, 0).Format("2006/01/02 15:04"),
	}

	return commentRes, err
}

// CommentList 评论列表
func (c *CommentService) CommentList(token string, videoId int64) (list *[]model.Comment, err error) {

	var commentListDB = commentDao.FindCommentListByVideoId(videoId)
	var commentList = make([]model.Comment, len(*commentListDB))

	for index, commentDB := range *commentListDB {
		var comment *model.Comment
		var user *model.User
		user, err = userService.FindUserModelById(commentDB.CommentUserId, token)
		if err != nil {
			return nil, err
		}

		comment = &model.Comment{
			Id:         commentDB.Id,
			User:       *user,
			Content:    commentDB.CommentText,
			CreateDate: time.Unix(commentDB.CreateAt/1000, 0).Format("2006/01/02 15:04"),
		}
		commentList[index] = *comment
	}
	return &commentList, err
}
