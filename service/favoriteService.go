package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	mapset "github.com/deckarep/golang-set"
	"strconv"
	"strings"
	"time"
)

type FavoriteService struct {
}

// NewFavoriteService 创建服务
func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (f *FavoriteService) DoLike(userId int64, videoId int64) {
	key := fmt.Sprintf("douyin:favorite:video%d:user%d", videoId, userId)
	fmt.Println(key)
	_, err := db.Redis.Set(key, 1, time.Minute*5).Result()
	if err != nil {
		panic(err)
	}
	return
}

func (f *FavoriteService) CancelLike(userId int64, videoId int64) {
	key := fmt.Sprintf("douyin:favorite:video%d:user%d", videoId, userId)
	fmt.Println(key)
	_, err := db.Redis.Set(key, 0, time.Minute*5).Result()
	if err != nil {
		panic(err)
	}
	return
}

func (f *FavoriteService) IsLike(userId int64, videoId int64) bool {
	// 先到缓存中查找
	temp := fmt.Sprintf("douyin:favorite:video%d:user%d", videoId, userId)
	tempTag, err := db.Redis.Get(temp).Result()
	// 如果在redis中存在的话，优先采用redis
	if err == nil {
		tag, _ := strconv.Atoi(tempTag)
		if tag == 1 {
			return true
		} else if tag == 0 {
			return false
		}
	}
	// 再到数据库中查找
	var favorite db.Favorite
	result := db.DB.Find(&favorite, "user_id = ? and video_id = ?", userId, videoId)
	if result.RowsAffected == 0 {
		return false
	} else if favorite.Tag == 0 {
		return false
	} else {
		return true
	}
}

func (f *FavoriteService) GetLikeList(userId int64) []model.Video {
	var videoList []model.Video
	videoIdSet := mapset.NewSet()

	// 先从mysql中获取用户已经点赞的视频
	var favorites []db.Favorite
	db.DB.Find(&favorites, "user_id = ?", userId)
	for _, favorite := range favorites {
		if favorite.Tag == 0 {
			continue
		}
		videoIdSet.Add(favorite.VideoID)
	}

	// 再从缓存中获取用户已经点赞的视频
	temp := fmt.Sprintf("douyin:favorite:*user%d", userId)
	videos, _ := db.Redis.Keys(temp).Result()
	for _, video := range videos {
		tempTag, _ := db.Redis.Get(video).Result()
		tag, _ := strconv.Atoi(tempTag)
		SplitData := strings.Split(video, ":")
		tempVideo := SplitData[2][5:]
		videoId, _ := strconv.ParseInt(tempVideo, 10, 64)
		if tag == 0 {
			if videoIdSet.Contains(videoId) {
				videoIdSet.Remove(videoId)
			}
			continue
		}
		videoIdSet.Add(videoId)
	}

	//将找到的视频id转换为视频列表
	for id := range videoIdSet.Iterator().C {
		fmt.Println(id)
		//TODO 用id查video
		// video := db.GetVideo(id.(int64))
		video := model.Video{
			Id:            uint(id.(int64)),
			Author:        model.User{},
			PlayUrl:       "",
			CoverUrl:      "",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		}
		fmt.Println(video)
		//var user model.User
		//videoVo := model.Video{}
		videoList = append(videoList, video)
	}
	return videoList
}
