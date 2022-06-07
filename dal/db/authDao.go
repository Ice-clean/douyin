package db

type Role struct { //角色
	Id   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:64"`
}

type UserRole struct { //用户角色关联
	Id     int64 `gorm:"primaryKey"`
	UserId int64
	RoleId int64
}

type Permission struct { //许可
	Id   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:64"`
	Uri  string `gorm:"size:256"`
}

type RolePermission struct { //角色许可关联
	Id           int64 `gorm:"primaryKey"`
	RoleId       int64
	PermissionId int64
}

type AuthDao struct {
}

// NewAuthDao 创建Dao
func NewAuthDao() *AuthDao {
	return &AuthDao{}
}

//FindPermissionByUserIdAndUri 根据用户 id 和请求 uri 进行关联查询
//err不为nil表示查询不到对应权限，返回false，反之返回true
func (a *AuthDao) FindPermissionByUserIdAndUri(userId int64, uri string) bool {
	var permission = new(Permission)
	err := DB.Raw("SELECT p.id AS id, p.name AS name, p.uri AS uri "+
		"FROM user_role ur, role_permission rp, permission p "+
		"where ur.role_id = rp.role_id and rp.permission_id = p.id and ur.user_id = ? and p.uri = ?",
		userId, uri).
		Take(permission).Error
	if err != nil {
		return false
	}
	return true
}

// CreateUserRole 创建对应的用户角色关系
func (a *AuthDao) CreateUserRole(userId int64, roleId int64) error {
	userRole := UserRole{
		UserId: userId,
		RoleId: roleId,
	}
	err := DB.Create(&userRole).Error
	return err
}
