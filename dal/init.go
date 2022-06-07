// Package dal 数据库层
// db 层放置数据库的初始化方法和各个数据库对象及其数据库操作
// sql 层存放表创建语句
package dal

import (
	"github.com/RaymondCode/simple-demo/dal/cache"
	"github.com/RaymondCode/simple-demo/dal/db"
)

func init() {
	db.Init()
	cache.Init()
}
