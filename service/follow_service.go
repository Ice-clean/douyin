package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dal/db"
	"github.com/RaymondCode/simple-demo/model"
)

type FollowService struct {
}

// NewFollowService 创建服务
func NewFollowService() *FollowService {
	return &FollowService{}
}

var userFollowDao = db.NewUserFollowDao()

// DoFollow 关注
func (f *FollowService) DoFollow(userId int64, toUserId int64) (err error) {
	key := fmt.Sprintf("douyin:follow:toUser%d:user%d", toUserId, userId)
	fmt.Println(key)
	var userFollow = userFollowDao.FindFollow(userId, toUserId)
	if userFollow.Id > 0 {
		return err
	}
	err = userFollowDao.UserFollow(userId, toUserId)
	if err != nil {
		return err
	}
	var userFollowCount *db.UserFollowCount
	userFollowCount, err = userDao.FindCountByID(userId)
	if err != nil {
		return err
	}
	fmt.Sprintf("asdasd")
	userFollowCount.FollowCount = userFollowCount.FollowCount + 1
	userFollowDao.UserFollowCount(userFollowCount)
	userFollowCount, err = userDao.FindCountByID(toUserId)
	if err != nil {
		return err
	}
	fmt.Sprintf("as123123123dasd")
	userFollowCount.FollowerCount = userFollowCount.FollowerCount + 1
	userFollowDao.UserFollowCount(userFollowCount)
	return err
}

// CancelLike 取消关注
func (f *FollowService) CancelFollow(userId int64, toUserId int64) (err error) {
	key := fmt.Sprintf("douyin:follow:toUser%d:user%d", toUserId, userId)
	fmt.Println(key)
	err = userFollowDao.UserCancelFollow(userId, toUserId)
	if err != nil {
		return err
	}
	var userFollowCount *db.UserFollowCount
	userFollowCount, err = userDao.FindCountByID(userId)
	if err != nil {
		return err
	}
	userFollowCount.FollowCount = userFollowCount.FollowCount - 1
	userFollowDao.UserFollowCount(userFollowCount)
	userFollowCount, err = userDao.FindCountByID(toUserId)
	if err != nil {
		return err
	}
	userFollowCount.FollowerCount = userFollowCount.FollowerCount - 1
	userFollowDao.UserFollowCount(userFollowCount)
	return err
}

// GetFollowList 获取关注列表
func (f *FollowService) GetFollowList(userId int64) []model.User {
	key := fmt.Sprintf("douyin:follow:user%d", userId)
	fmt.Println(key)
	var userFollows = userFollowDao.FindUserFollow(userId)
	var userList = make([]model.User, len(userFollows))
	for count := range userFollows {
		var userFollow = userFollows[count]
		var user = new(model.User)
		var userInfo *db.User
		var _ error
		userInfo, _ = userDao.FindUserByID(userFollow.FollowedUserId)
		user.Id = userInfo.Id
		user.Name = userInfo.Name
		user.FollowCount = int64(len(userFollowDao.FindUserFollow(userInfo.Id)))
		user.FollowerCount = int64(len(userFollowDao.FindUserFollower(userInfo.Id)))
		user.IsFollow = true
		userList[count] = *user
	}
	return userList
}

// GetFollowerList 获取粉丝列表
func (f *FollowService) GetFollowerList(userId int64) []model.User {
	key := fmt.Sprintf("douyin:follow:user%d", userId)
	fmt.Println(key)
	var userFollowers = userFollowDao.FindUserFollower(userId)
	var userList = make([]model.User, len(userFollowers))
	for count := range userFollowers {
		var userFollower = userFollowers[count]
		var user = new(model.User)
		var userInfo *db.User
		var _ error
		userInfo, _ = userDao.FindUserByID(userFollower.FollowUserId)
		user.Id = userInfo.Id
		user.Name = userInfo.Name
		user.FollowCount = int64(len(userFollowDao.FindUserFollow(userInfo.Id)))
		user.FollowerCount = int64(len(userFollowDao.FindUserFollower(userInfo.Id)))
		user.IsFollow = userDao.JudgeFollow(userId, userInfo.Id)
		userList[count] = *user
	}
	return userList
}
