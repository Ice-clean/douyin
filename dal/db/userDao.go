package db

import (
	"fmt"
	"time"
)

// User DO对象，存放用户信息
type User struct {
	Id       int64  `gorm:"primaryKey"` //用户id
	Username string `gorm:"size:64"`    //用户名
	Password string `gorm:"size:64"`    //密码
	Name     string `gorm:"size:64"`    //昵称

	CreateTime        int64  `gorm:"autoCreateTime:milli"`        //注册时间（毫秒值）
	LastLogin         int64  `gorm:"autoCreateTime:milli"`        //上次登录时间(毫秒值)
	Freeze            int8   `gorm:"default:0"`                   //是否被冻结（0表示正常，1表示被冻结）
	Age               int    `gorm:"default:20"`                  //年龄（默认20）
	PersonalSignature string `gorm:"default:'你们喜欢的话题，就是我们采访的内容'"` //个性签名
	Site              string `gorm:"default:'杭州';size:64"`        //地点
}

// UserFollowCount DO对象 存放用户关注数量与粉丝数量
type UserFollowCount struct {
	Id            int64 `gorm:"primaryKey"` //用户id
	FollowCount   int64 `gorm:"default:0"`  //关注数量
	FollowerCount int64 `gorm:"default:0"`  //粉丝数量
}

// UserFollow DO对象 维护用户之间关注与被关注的关系
type UserFollow struct {
	Id             int64 `gorm:"primaryKey"` //主键
	FollowUserId   int64 //关注操作的用户id
	FollowedUserId int64 //被关注的用户id
}

type UserDao struct {
}

// NewUserDao 创建Dao
func NewUserDao() *UserDao {
	return &UserDao{}
}

// CreateUser  创建新用户
func (u *UserDao) CreateUser(username string, password string) (us *User, err error) {

	user := User{
		Username:   username,
		Password:   password,
		Name:       "张三",
		CreateTime: time.Now().UnixMilli(),
		LastLogin:  time.Now().UnixMilli(),
	}
	//创建对应的用户
	if err = DB.Create(&user).Error; err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", user)

	//创建对应的用户关注和粉丝数量记录
	userFollowCount := UserFollowCount{Id: user.Id}
	DB.Create(&userFollowCount)

	return &user, err

	//db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//db.AutoMigrate(&User{})
	//db.SingularTable(true)
}

// FindUserByUsernameAndPassword 根据用户名和密码查询用户
func (u *UserDao) FindUserByUsernameAndPassword(username string, password string) (user *User, err error) {
	user = new(User)
	if err = DB.Where("username = ? AND password = ?", username, password).First(user).Error; err != nil {
		//用户不存在，登录失败
		return nil, err
	}
	//登录成功
	return
}

// FindUserByUsername 判断用户名可用，返回true表示用户名可用，false表示用户名不可用
func (u *UserDao) FindUserByUsername(username string) bool {
	var user = new(User)
	if err := DB.Where("username = ?", username).First(user).Error; err != nil {
		//用户名可用（未存在)
		return true
	}
	return false //用户名已存在（不可用）
}

// FindUserByID 根据id查询用户信息
func (u *UserDao) FindUserByID(id int64) (user *User, err error) {
	user = new(User)
	if err = DB.First(user, id).Error; err != nil {
		return nil, err
	}
	return
}

// FindCountByID 根据用户id查询用户关注数和粉丝数
func (u *UserDao) FindCountByID(id int64) (userFollowCount *UserFollowCount, err error) {
	userFollowCount = new(UserFollowCount)
	if err = DB.First(userFollowCount, id).Error; err != nil {
		return nil, err
	}
	return
}

// JudgeFollow 判断是否关注该用户，true表示关注，false表示未关注
func (u *UserDao) JudgeFollow(followUserId int64, followedUserId int64) bool {
	var userFollow = new(UserFollow)
	err := DB.Where("follow_user_id = ? AND followed_user_id = ?", followUserId, followedUserId).First(userFollow).Error
	if err != nil {
		//没有关注该用户
		return false
	}
	return true //关注了该用户
}

//关注操作
func (u *UserDao) UserFollow(userId int64, toUserId int64) (err error) {
	userFollow := UserFollow{
		FollowUserId:   userId,
		FollowedUserId: toUserId,
	}
	return DB.Create(&userFollow).Error
}

//取消关注操作
func (u *UserDao) UserCancelFollow(userId int64, toUserId int64) (err error) {
	return DB.Where("follow_user_id = ? AND followed_user_id = ?", userId, toUserId).Delete(UserFollow{}).Error
}
func (u *UserDao) FindFollow(userId int64, toUserId int64) UserFollow {
	var userFollow UserFollow
	DB.Where("follow_user_id = ? AND followed_user_id = ?", userId, toUserId).Find(&userFollow)
	return userFollow
}

//获取关注列表
func (u *UserDao) FindUserFollow(userId int64) []UserFollow {
	var userFollows []UserFollow
	DB.Where("follow_user_id = ?", userId).Find(&userFollows)
	return userFollows
}

//获取粉丝列表
func (u *UserDao) FindUserFollower(userId int64) []UserFollow {
	var userFollows []UserFollow
	DB.Where("followed_user_id = ?", userId).Find(&userFollows)
	return userFollows
}

// UpdateLastLoginTime 更新用户最后一次的登录时间
func (u *UserDao) UpdateLastLoginTime(id int64) (err error) {
	err = DB.Model(&User{}).Where("id = ?", id).Update("last_login", time.Now().UnixMilli()).Error
	if err != nil { //更新失败
		return err
	}
	return //更新成功
}
