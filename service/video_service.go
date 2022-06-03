package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/constant"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

type VideoService struct {
}

// NewVideoService 创建服务
func NewVideoService() *VideoService {
	return &VideoService{}
}

// GetVideoById 通过视频 ID 获取视频对象
func (v *VideoService) GetVideoById(videoId int64) db.Video {
	return db.GetVideo(videoId)
}

// PublishVideo 发布视频
// 发布成功则返回 fileName，否则抛出错误 err
func (v *VideoService) PublishVideo(user *model.User, title string, file *multipart.FileHeader) (finalName string, err error) {
	// 获取文件名以及生成最终的文件名
	filename := filepath.Base(file.Filename)
	fmt.Printf("视频：%v", user)
	finalName = fmt.Sprintf("%d_%s", user.Id, filename)

	// 设置存放路径，并保存到服务器
	c := gin.Context{}
	saveFile := filepath.Join("./public/", finalName)
	if err = c.SaveUploadedFile(file, saveFile); err != nil {
		return
	}

	// 生成缩略图并保存到服务器
	snapshotName := utils.GetSnapshot(saveFile, filepath.Join("./public/", finalName), 50)

	// 封装视频对象
	video := &db.Video{}
	video.UserId = user.Id
	video.Tag = v.parseTag(title)
	video.PlayUrl = fmt.Sprintf("%s/%s/%s", constant.Host, "static", finalName)
	video.CoverUrl = fmt.Sprintf("%s/%s/%s", constant.Host, "static", snapshotName)

	// 将视频信息保存到数据库
	db.CreateVideo(video)
	return
}

// GetPublishList 获取视频列表，并将 DO 对象转成 VO 对象
func (v *VideoService) GetPublishList(userId int64) []model.Video {
	// 获取视频列表，并准备转化成投稿列表
	var videoList = db.GetPublishByUserId(userId)
	var publishList = make([]model.Video, len(videoList))

	// TODO 获取投稿用户（等用户功能完成）
	// user := userService.GetUserById(userId)
	var user model.User

	// 对象转化
	// favorite := videoService.IsFavorite(user， videoList)
	for _, v := range videoList {
		// TODO 判断用户是否点赞本视频

		var favorite bool
		publishList = append(publishList, model.Video{
			Id:            v.ID,
			Author:        user,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    favorite,
		})
	}
	return publishList
}

// getTag 解析标题中的标签，若有标签则去除标题中该标签
func (v *VideoService) parseTag(title string) string {
	// 取出第一行，作正则匹配
	tagString := strings.Split(title, "\n")[0]
	compile := regexp.MustCompile("#\\S+")
	tags := compile.FindAllString(tagString, -1)

	// 拼接标签为字符串
	tagBuilder := strings.Builder{}
	for _, v := range tags {
		tagBuilder.Write([]byte(v))
	}
	return tagBuilder.String()
}
