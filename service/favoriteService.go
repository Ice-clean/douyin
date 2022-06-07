package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	mapset "github.com/deckarep/golang-set"
	"regexp"
	"sort"
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

// DoLike 点赞
func (f *FavoriteService) DoLike(userId int64, videoId int64) {
	key := fmt.Sprintf("douyin:favorite:video%d:user%d", videoId, userId)
	fmt.Println(key)
	_, err := db.Redis.Set(key, 1, time.Minute*5).Result()
	if err != nil {
		panic(err)
	}
	return
}

// CancelLike 取消点赞
func (f *FavoriteService) CancelLike(userId int64, videoId int64) {
	key := fmt.Sprintf("douyin:favorite:video%d:user%d", videoId, userId)
	fmt.Println(key)
	_, err := db.Redis.Set(key, 0, time.Minute*5).Result()
	if err != nil {
		panic(err)
	}
	return
}

// IsLike 判断用户是否喜欢
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

// GetLikeList 获取用户喜欢的视频列表
func (f *FavoriteService) GetLikeList(userId int64) []model.Video {
	videoIdList := f.GetLikeId(userId)
	var videoList = make([]model.Video, len(videoIdList))

	//将找到的视频id转换为视频列表
	for id := range videoIdList {
		fmt.Println(id)
		video := NewVideoService().GetVideoById(int64(id))
		user := NewUserService().FindUserById(video.UserId)
		userVo := NewUserService().ToUserVO(*user)
		videoVo := NewVideoService().ToVideoVO(video, *userVo, true)
		videoList = append(videoList, videoVo)
	}
	return videoList
}

// Recommend 对用户喜欢的视频进行推荐
func (f *FavoriteService) Recommend(userId int64) []db.Video {
	var videoList []db.Video
	// 根据用户点赞过的视频进行推荐
	videoIdList := f.GetLikeId(userId)
	type APIVideo struct {
		Tag string
	}
	var apiVideos = make([]APIVideo, len(videoList))
	var TagCount = make(map[string]int)
	db.DB.Model(&db.Video{}).Find(&apiVideos, videoIdList)
	compile := regexp.MustCompile(`[\p{Han}]+`)
	// 收集每个tag的数量
	for _, apiVideo := range apiVideos {
		tags := compile.FindAllString(apiVideo.Tag, -1)
		for _, tag := range tags {
			_, ok := TagCount[tag]
			if ok {
				TagCount[tag] += 1
			} else {
				TagCount[tag] = 1
			}
		}
	}
	// 对tag进行降序排序
	type count struct {
		Tag      string
		TagCount int
	}
	var counts []count
	for k, v := range TagCount {
		counts = append(counts, count{k, v})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].TagCount > counts[j].TagCount
	})
	// 拿到top3的tag
	var TopTags []string
	for i, v := range counts {
		if i > 2 {
			break
		} else {
			TopTags = append(TopTags, v.Tag)
		}
	}
	db.DB.Where("1=1")
	for i, _ := range TopTags {
		sql := fmt.Sprintf("%%%s%%", TopTags[i])
		db.DB.Or("tag LIKE ?", sql)
	}
	db.DB.Find(&videoList).Limit(50)
	return videoList
}

// GetLikeId 获取用户喜欢的视频的id列表
func (f *FavoriteService) GetLikeId(userId int64) []interface{} {
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
	videoIdList := videoIdSet.ToSlice()
	return videoIdList
}
