package service

import "github.com/RaymondCode/simple-demo/dal/db"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

//注入dao对象
var authDao = db.NewAuthDao()

//IsPermit 判断是否拥有该权限
func (a *AuthService) IsPermit(userId int64, uri string) bool {
	return authDao.FindPermissionByUserIdAndUri(userId, uri)
}
