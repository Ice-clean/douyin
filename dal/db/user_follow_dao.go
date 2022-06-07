package db

type UserFollowDao struct {
}

// NewUserDao 创建Dao
func NewUserFollowDao() *UserFollowDao {
	return &UserFollowDao{}
}

//关注操作
func (u *UserFollowDao) UserFollow(userId int64, toUserId int64) (err error) {
	userFollow := UserFollow{
		FollowUserId:   userId,
		FollowedUserId: toUserId,
	}
	return DB.Create(&userFollow).Error
}

func (u *UserFollowDao) UserFollowCount(userFollowCount *UserFollowCount) {
	DB.Save(&userFollowCount)
}

//取消关注操作
func (u *UserFollowDao) UserCancelFollow(userId int64, toUserId int64) (err error) {
	return DB.Where("follow_user_id = ? AND followed_user_id = ?", userId, toUserId).Delete(UserFollow{}).Error
}

func (u *UserFollowDao) FindFollow(userId int64, toUserId int64) UserFollow {
	var userFollow UserFollow
	DB.Where("follow_user_id = ? AND followed_user_id = ?", userId, toUserId).Find(&userFollow)
	return userFollow
}

//获取关注列表
func (u *UserFollowDao) FindUserFollow(userId int64) []UserFollow {
	var userFollows []UserFollow
	DB.Where("follow_user_id = ?", userId).Find(&userFollows)
	return userFollows
}

//获取粉丝列表
func (u *UserFollowDao) FindUserFollower(userId int64) []UserFollow {
	var userFollows []UserFollow
	DB.Where("followed_user_id = ?", userId).Find(&userFollows)
	return userFollows
}
