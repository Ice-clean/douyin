package db

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserId  int64
	VideoID int64
	Tag     int32
}
