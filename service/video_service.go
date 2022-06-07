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
	"strconv"
	"strings"
)

type VideoService struct {
}

// 注入 videoDao
var videoDao = db.NewVideoDao()

// NewVideoService 创建服务
func NewVideoService() *VideoService {
	return &VideoService{}
}

// GetVideoById 通过视频 ID 获取视频对象
func (v *VideoService) GetVideoById(videoId int64) db.Video {
	return videoDao.GetVideoById(videoId)
}

// GetVideoList 分页获取视频列表（页数存储在缓存中，默认每页数量为 10）
// 其中 userIP 为请求者的 IP 地址，用来唯一标识用户，获取当前用户正在查看的视频页数
func (v VideoService) GetVideoList(userIP string) []model.Video {
	page := 0
	// 先从缓存中获取当前用户查看的视频页数（没有则默认 0）
	pageString, err := Redis.HGet(constant.UserVideoPage, userIP).Result()
	if err == nil && pageString != "" {
		page, _ = strconv.Atoi(pageString)
	}
	// 最后从数据库取出视频列表
	videoList := videoDao.GetVideoList(page, 10)
	// 判断是否还有列表
	if len(videoList) > 0 {
		// 有的话 page 继续 +1
		Redis.HSet(constant.UserVideoPage, userIP, page+1)
	} else {
		// 没有的话，page 归 0，循环播放
		Redis.HSet(constant.UserVideoPage, userIP, 0)
		videoList = videoDao.GetVideoList(0, 10)
	}
	// 转化为响应对象然后返回
	return v.ToVideoVOList(videoList, -1, "")
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
	utils.GetSnapshot(saveFile, filepath.Join("./public/", finalName), 50)

	// 封装视频对象
	video := &db.Video{}
	video.UserId = user.Id
	video.Tag = v.parseTag(title)
	video.PlayUrl = fmt.Sprintf("%s/%s/%s", constant.Host, "static", finalName)
	video.CoverUrl = fmt.Sprintf("%s/%s/%s.jpeg", constant.Host, "static", finalName)
	video.Title = title

	// 将视频信息保存到数据库
	videoDao.CreateVideo(video)
	return
}

// GetPublishList 获取视频列表，并将 DO 对象转成 VO 对象
func (v *VideoService) GetPublishList(userId int64, token string) []model.Video {
	// 获取视频列表，并准备转化成投稿列表
	var videoList = videoDao.GetPublishByUserId(userId)
	// 将视频列表转化成响应对象再返回
	return v.ToVideoVOList(videoList, userId, token)
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

// GetRecommend 根据用户情况获取推荐视频并封装成响应对象
func (v *VideoService) GetRecommend(userId int64, token string) []model.Video {
	// 获取推荐的视频列表
	videoList := NewFavoriteService().Recommend(userId)
	// 将视频列表转化为响应对象再返回
	return v.ToVideoVOList(videoList, userId, token)
}

// ToVideoVO 将视频 DO 对象转化为 VO 对象
func (v *VideoService) ToVideoVO(videoDO db.Video, author model.User, favorite bool) model.Video {
	return model.Video{
		Id:            videoDO.ID,
		Author:        author,
		PlayUrl:       videoDO.PlayUrl,
		CoverUrl:      videoDO.CoverUrl,
		FavoriteCount: videoDO.FavoriteCount,
		CommentCount:  videoDO.CommentCount,
		IsFavorite:    favorite,
		Title:         videoDO.Title,
	}
}

// ToVideoVOList 将视频 DO 列表转化成 VO 列表
func (v *VideoService) ToVideoVOList(videoDOList []db.Video, userId int64, token string) []model.Video {
	var videoVOList = make([]model.Video, len(videoDOList))
	userService := NewUserService()
	// 将视频列表转化为响应对象
	var favorite bool
	var user *model.User
	for i, videoDO := range videoDOList {
		// 判断用户是否点赞本视频
		if userId == -1 {
			favorite = false
		} else {
			favorite = NewFavoriteService().IsLike(userId, int64(videoDO.ID))
		}
		fmt.Println("视频：", videoDO.ID, "用户：", userId, "是否点赞？:", favorite)
		// 获取视频发布者
		if token == "" {
			user = userService.ToUserVO(*userService.FindUserById(videoDO.UserId))
		} else {
			user, _ = userService.UserInfo(strconv.Itoa(int(videoDO.UserId)), token)
		}

		// 将视频转化为 VO 对象，并存入列表
		videoVO := v.ToVideoVO(videoDO, *user, favorite)
		videoVOList[i] = videoVO
	}
	return videoVOList
}
