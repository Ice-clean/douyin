package db

import (
	"gorm.io/gorm"
)

// Video DO 对象
type Video struct {
	gorm.Model
	UserId        int64  // 发布者 id
	PlayUrl       string // 视频 url
	CoverUrl      string // 封面 url
	Tag           string // 视频标签
	FavoriteCount int64  // 点赞数量
	CommentCount  int64  // 评论数量
	Title         string // 视频标题
}

type VideoDao struct {
}

func NewVideoDao() *VideoDao {
	return &VideoDao{}
}

// CreateVideo 创建视频
func (v *VideoDao) CreateVideo(video *Video) bool {
	return DB.Create(video).Error != nil
}

// GetPublishByUserId 通过用户 ID 获取该用户发布的视频列表
func (v *VideoDao) GetPublishByUserId(userId int64) []Video {
	var videoList []Video
	DB.Where("user_id = ?", userId).Find(&videoList)
	return videoList
}

func (v *VideoDao) GetVideoById(videoId int64) Video {
	var video Video
	DB.Where("id = ?", videoId).First(&video)
	return video
}

// GetVideoList 分页获取视频列表
func (v *VideoDao) GetVideoList(page, num int) []Video {
	var videoList []Video
	DB.Order("id DESC").Limit(num).Offset(page * num).Find(&videoList)
	return videoList
}
