package controller

import (
	"github.com/RaymondCode/simple-demo/constant"
	"github.com/RaymondCode/simple-demo/model"
)

var DemoVideos = []model.Video{
	{
		Id:            2,
		Author:        DemoUser,
		PlayUrl:       constant.Host + "/static/1_1_bear.mp4",
		CoverUrl:      constant.Host + "/static/1_1_bear.mp4.jpeg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "http://39.98.41.126:31101/static/1_多屏协同.mp4",
		CoverUrl:      constant.Host + "/static/1_多屏协同.mp4.jpeg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []model.Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = model.User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
