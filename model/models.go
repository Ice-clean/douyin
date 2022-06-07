// Package model 这一层存放响应对象
package model

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type Video struct {
	Id            uint   `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count"`   //零值也不希望被忽略
	FollowerCount   int64  `json:"follower_count"` //零值也不希望被忽略
	IsFollow        bool   `json:"is_follow"`      //零值也不希望被忽略
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	FavoriteCount   int64  `json:"favorite_count"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
}
