package service

import (
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"strconv"
	"time"
)

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
	//创建用户，返回插入的主键id
	userId, err = userDao.CreateUser(username, password)
	if err != nil {
		return -1, "", err
	}
	//生成token
	token, err = utils.GenerateToken(userId, username)
	Redis.Set("douyin:"+token, userId, 7*24*time.Hour) //存入redis
	return
}

// Login 具体登录逻辑
func (u *UserService) Login(username string, password string) (userId int64, token string, err error) {
	var user *db.User
	user, err = userDao.FindUserByUsernameAndPassword(username, password)
	if err != nil {
		return -1, "", err
	}

	//生成token
	token, err = utils.GenerateToken(userId, username)
	Redis.Set("douyin:"+token, userId, 7*24*time.Hour) //存入redis，有效期7天
	return user.Id, token, err                         //登录成功
}

// UserInfo 查看用户信息
func (u *UserService) UserInfo(userId string, token string) (user *model.User, err error) {
	//解析token
	var claims *utils.Claims
	claims, err = utils.ParseToken(token)
	if err != nil { //解析失败，直接返回
		return nil, err
	}

	user = new(model.User)

	var userIdInt int64
	userIdInt, err = strconv.ParseInt(userId, 10, 64)
	//根据id查找用户的关注量与粉丝量
	var userFollowCount *db.UserFollowCount
	userFollowCount, err = userDao.FindCountByID(userIdInt)

	if userIdInt == claims.Id { //查找用户为自己
		user = &model.User{
			Id:            userIdInt,
			Name:          claims.Username, //这里先暂且使用用户名
			FollowCount:   userFollowCount.FollowCount,
			FollowerCount: userFollowCount.FollowerCount,
		}
	} else { //查看其它用户
		//判断是否关注了该用户,true关注，false未关注
		isFollow := userDao.JudgeFollow(claims.Id, userIdInt)
		user = &model.User{
			Id:            userIdInt,
			Name:          claims.Username, //这里先暂且使用用户名
			FollowCount:   userFollowCount.FollowCount,
			FollowerCount: userFollowCount.FollowerCount,
			IsFollow:      isFollow,
		}
	}
	return
}

func (u *UserService) Logout(token string) {
	//删除redis中的token信息，下次登录校验就不会通过
	Redis.Del(token)
}

// IsUsernameCanUse 判断用户名可用，返回true表示用户名可用，false表示用户名不可用
func (u *UserService) IsUsernameCanUse(username string) bool {
	return userDao.FindUserByUsername(username)
}
