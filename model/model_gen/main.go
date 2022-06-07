package main

import (
	_ "github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/dal/db"
)

func main() {
	_ = db.DB.AutoMigrate(&db.Favorite{})
}
