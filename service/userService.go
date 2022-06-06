package service

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"strconv"
	"time"
)

const UserPrefix = "douyin:user:"

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

//注入dao对象
var userDao = db.NewUserDao()

// Redis 注入redis客户端
var Redis = db.Redis

// Register 具体注册逻辑
func (u *UserService) Register(username string, password string) (userId int64, token string, err error) {
	var user *db.User
	//创建用户。并返回封装数据的user
	user, err = userDao.CreateUser(username, password)
	if err != nil {
		return -1, "", err
	}
	fmt.Println("registerUser: ", *user) //打印日志

	//生成token
	token, err = utils.GenerateToken(user.Id, username)

	//序列化 user 并存入redis
	data, _ := json.Marshal(user)
	fmt.Println(string(data))                                 //打印序列化结果
	Redis.Set(UserPrefix+token, string(data), 7*24*time.Hour) //存入redis，有效期7天
	return user.Id, token, err
}

// Login 具体登录逻辑
func (u *UserService) Login(username string, password string) (userId int64, token string, err error) {
	var user *db.User
	user, err = userDao.FindUserByUsernameAndPassword(username, password)
	if err != nil {
		return -1, "", err
	}

	err = userDao.UpdateLastLoginTime(user.Id) //更新最后一次登录时间
	fmt.Println("loginUser: ", *user)          //打印日志

	//生成token
	token, err = utils.GenerateToken(user.Id, username)

	//序列化 user 并存入redis
	data, _ := json.Marshal(user)
	fmt.Println(string(data))                                 //打印序列化结果
	Redis.Set(UserPrefix+token, string(data), 7*24*time.Hour) //存入redis，有效期7天
	return user.Id, token, err                                //登录成功
}

// UserInfo 查看用户信息
func (u *UserService) UserInfo(userId string, token string) (user *model.User, err error) {
	var userIdInt int64
	userIdInt, err = strconv.ParseInt(userId, 10, 64)
	//解析token
	//var claims *utils.Claims
	//claims, err = utils.ParseToken(token)
	//if err != nil { //解析失败，直接返回
	//	return nil, err
	//}
	ownUserDb := u.FindUserByToken(token) //根据token查找登录用户信息
	//user = new(model.User)

	//根据id查找用户的关注量与粉丝量
	var userFollowCount *db.UserFollowCount
	userFollowCount, err = userDao.FindCountByID(userIdInt)

	if userIdInt == ownUserDb.Id { //查找用户为自己
		user = &model.User{
			Id:             userIdInt,
			Name:           ownUserDb.Username, //这里先暂且使用 Username
			FollowCount:    userFollowCount.FollowCount,
			FollowerCount:  userFollowCount.FollowerCount,
			Signature:      ownUserDb.PersonalSignature,
			TotalFavorited: 66,
			FavoriteCount:  88,
		}
	} else { //查看其它用户
		isFollow := userDao.JudgeFollow(ownUserDb.Id, userIdInt) //判断是否关注了该用户,true关注，false未关注
		otherUserDb := u.FindUserById(userIdInt)                 //根据id查找 查看的用户信息
		user = &model.User{
			Id:             userIdInt,
			Name:           otherUserDb.Username, //这里先暂且使用 Username
			FollowCount:    userFollowCount.FollowCount,
			FollowerCount:  userFollowCount.FollowerCount,
			IsFollow:       isFollow,
			Signature:      otherUserDb.PersonalSignature,
			TotalFavorited: 666,
			FavoriteCount:  888,
		}
	}
	return
}

// FindUserById 根据id从数据库从查询用户信息
func (u *UserService) FindUserById(id int64) *db.User {
	user, err := userDao.FindUserByID(id)
	if err != nil {
		return nil
	}
	return user
}

// FindUserByToken 根据token从redis中查询用户信息
func (u *UserService) FindUserByToken(token string) *db.User {
	value, err := Redis.Get(UserPrefix + token).Result()

	if err != nil {
		fmt.Println("token doesn't exist") //打印日志
		return nil
	}
	fmt.Println(value) //打印日志（redis中的 json 字符串）

	//反序列化并把结果保存到user
	user := &db.User{}
	err = json.Unmarshal([]byte(value), user)
	if err != nil {
		fmt.Println("json unmarshal failed!") //打印日志
	}
	return user
}

// FindTotalFavoritedByUserId 根据id查找用户的总获赞数
func (u *UserService) FindTotalFavoritedByUserId(userId int64) int64 {
	return userId //未实现！！！！！！
}

// FindFavoriteCountByUserId  根据id查找用户的点赞作品总数
func (u *UserService) FindFavoriteCountByUserId(userId int64) int64 {
	return userId //未实现！！！！！！
}

// Logout 退出登录（客户端无对应接口）
func (u *UserService) Logout(token string) {
	//删除redis中的token信息，下次登录校验就不会通过
	Redis.Del(token)
}

// IsUsernameCanUse 判断用户名可用，返回true表示用户名可用，false表示用户名不可用
func (u *UserService) IsUsernameCanUse(username string) bool {
	return userDao.FindUserByUsername(username)
}

// ToUserVO 将 user DO 对象转化为 VO 对象
func (u *UserService) ToUserVO(userDO db.User) *model.User {
	//根据id查找用户的关注量与粉丝量
	var userFollowCount *db.UserFollowCount
	userFollowCount, _ = userDao.FindCountByID(userDO.Id)

	return &model.User{
		Id:             userDO.Id,
		Name:           userDO.Username, //这里先暂且使用 Username
		FollowCount:    userFollowCount.FollowCount,
		FollowerCount:  userFollowCount.FollowerCount,
		Signature:      userDO.PersonalSignature,
		TotalFavorited: 66,
		FavoriteCount:  88,
	}
}
