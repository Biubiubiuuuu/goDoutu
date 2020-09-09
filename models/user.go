package models

import (
	"time"

	"github.com/Biubiubiuuuu/goDoutu/db/mysql"
)

// 用户
type User struct {
	Model
	Telephone string    `gorm:"unique;size:50;" json:"telephone"` // 手机号
	Password  string    `json:"-"`                                // 密码
	Birthday  time.Time `json:"birthday"`                         // 生日 格式：yyyy-MM-dd HH:mm:ss
	Sex       int64     `json:"sex"`                              // 性别 0:未知 1:男 2:女
	Avatar    string    `gorm:"size:200" json:"avatar"`           // 头像
	Email     string    `gorm:"unique;size:50;" json:"email"`     // 邮箱
	QQ        string    `gorm:"size:200" json:"QQ"`               // QQ
	Nickname  string    `gorm:"size:150" json:"nickname"`         // 微信用户昵称
	OpenID    string    `gorm:"size:150" json:"open_id"`          // 微信用户唯一标识
	Country   string    `gorm:"size:30" json:"country"`           // 所在国家
	Province  string    `gorm:"size:30" json:"province"`          // 所在省份
	City      string    `gorm:"size:30" json:"city"`              // 所在城市
	Token     string    `gorm:"size:255" json:"token"`            // token
}

// 粉丝
type UserFans struct {
	ID       int64 `gorm:"index"`
	BeUserID int64 // 被关注用户ID
	UserID   int64 // 关注用户ID
}

// 我的关注
type UserFollows struct {
	ID       int64 `gorm:"index"`
	BeUserID int64 // 被关注用户ID
	UserID   int64 // 关注用户ID
}

//var buf strings.Builder
//buf.WriteString("%")
//buf.WriteString(v.(string))
//buf.WriteString("%")

// 用户登录 by
//  telephone
//  password
func (u *User) Login() error {
	db := mysql.GetMysqlDB()
	return db.Where("telephone = ? ", u.Telephone).Where("password = ?", u.Password).First(&u).Error
}

// 获取用户信息 by token
func (u *User) QueryUserByToken() error {
	db := mysql.GetMysqlDB()
	return db.Where("token = ? AND ISNULL(token)=0 AND LENGTH(trim(token))>0", u.Token).First(&u).Error
}

// 获取用户信息 by open_id
func (u *User) QueryUserByOpenID() error {
	db := mysql.GetMysqlDB()
	return db.Where("open_id = ? AND ISNULL(open_id)=0 AND LENGTH(trim(open_id))>0", u.OpenID).First(&u).Error
}

// 修改用户信息
func (u *User) UpdatesUserByID(args interface{}) error {
	db := mysql.GetMysqlDB()
	return db.Model(&u).Updates(args).Error
}

// 用户信息详情
func (u *User) QueryUserByID() error {
	db := mysql.GetMysqlDB()
	return db.First(&u, u.ID).Error
}

// 我的粉丝
func (u *UserFans) QueryMyFansByBeUserID() (fans []UserFans) {
	db := mysql.GetMysqlDB()
	db.Where("be_user_id = ?", u.BeUserID).Find(&fans)
	return
}

// 我的关注
func (u *UserFollows) QueryMyFollowsByBeUserID() (follows []UserFollows) {
	db := mysql.GetMysqlDB()
	db.Where("be_user_id = ?", u.BeUserID).Find(&follows)
	return
}

// 新增用户
func (u *User) NewUser() error {
	db := mysql.GetMysqlDB()
	return db.Create(&u).Error
}
