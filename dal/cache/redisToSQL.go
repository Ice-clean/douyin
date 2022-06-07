package cache

import (
	"github.com/RaymondCode/simple-demo/dal/db"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

type Fn func()

type MyTicker struct {
	MyTick *time.Ticker
	Runner Fn
}

func NewMyTick(interval int, f Fn) *MyTicker {
	return &MyTicker{
		MyTick: time.NewTicker(time.Duration(interval) * time.Second),
		Runner: f,
	}
}

func (t *MyTicker) Start() {
	for {
		select {
		case <-t.MyTick.C:
			t.Runner()
		}
	}
}

func redisToSQL() {
	datas := db.Redis.Keys("douyin:favorite*").Val()
	var favorites []db.Favorite
	for _, data := range datas {
		SplitData := strings.Split(data, ":")
		video := SplitData[2][5:]
		user := SplitData[3][4:]
		videoId, _ := strconv.ParseInt(video, 10, 64)
		userId, _ := strconv.ParseInt(user, 10, 64)
		result, _ := db.Redis.Get(data).Result()
		tag, _ := strconv.Atoi(result)
		favorites = append(favorites, db.Favorite{
			UserId:  userId,
			VideoID: videoId,
			Tag:     int32(tag),
		})
	}
	db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "video_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tag"}),
	}).Create(&favorites)
}
