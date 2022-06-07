package db

import (
	"fmt"
	"time"
)

// Comment DO对象，存放评论信息
type Comment struct {
	Id             int64  `gorm:"primaryKey"`
	CommentText    string `gorm:"size:255"`
	CommentUserId  int64
	CommentVideoId int64
	CreateAt       int64 `gorm:"autoCreateTime:milli"`
	Status         int8  `gorm:"default:1"`
}

type CommentDao struct {
}

// NewCommentDao 创建Dao
func NewCommentDao() *CommentDao {
	return &CommentDao{}
}

// CreateComment 创建评论
func (c *CommentDao) CreateComment(videoId int64, userId int64, commentText string) (comment *Comment, err error) {
	comment = &Comment{
		CommentVideoId: videoId,
		CommentUserId:  userId,
		CommentText:    commentText,
		CreateAt:       time.Now().UnixMilli(),
	}
	err = DB.Create(comment).Error
	if err != nil {
		fmt.Println("create comment failed...")
		return nil, err
	}

	fmt.Println(*comment)
	return comment, err
}

// DeleteComment 根据 id 删除评论
func (c *CommentDao) DeleteComment(commentId int64) (comment *Comment) {
	comment = &Comment{
		Id: commentId,
	}
	DB.Delete(comment)
	return comment
}

// FindCommentListByVideoId 根据 videoId 查询评论列表，按评论时间倒序排列
func (c *CommentDao) FindCommentListByVideoId(videoId int64) *[]Comment {
	var commentList []Comment
	DB.Raw("SELECT * from comment WHERE comment_video_id = ? order by create_at desc", videoId).Scan(&commentList)
	return &commentList
}
