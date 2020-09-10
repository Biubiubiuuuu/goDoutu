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

// 关注 or 被关注
type UserFans struct {
	ID       int64 `gorm:"index"`
	UserID   int64 // 关注用户ID
	BeUserID int64 // 被关注用户ID
}

// 我的收藏
type UserEmoticons struct {
	ID          int64 `gorm:"index"`
	UserID      int64 // 用户ID
	EmoticonsID int64 // 表情包ID
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

// 添加关注
func (u *UserFans) AddFans() error {
	db := mysql.GetMysqlDB()
	return db.Create(&u).Error
}

// 取消关注
func (u *UserFans) CancelFans() error {
	db := mysql.GetMysqlDB()
	return db.Unscoped().Delete(&u).Error
}

// 我的粉丝
func (u *User) QueryMyFansByBeUserID() (count int, user []User) {
	db := mysql.GetMysqlDB()
	query := db.Table("user")
	query = query.Select("id in (SELECT user_id FROM user_fans WHERE be_user_id = ?)", u.ID)
	query.Count(&count)
	query.Find(&user)
	return
}

// 查看是否关注用户
func (u *User) QueryFansIsBe(be_user_id int64) bool {
	db := mysql.GetMysqlDB()
	query := db.Table("user_fans")
	var fans []UserFans
	query = query.Select("be_user_id = ? AND user_id = ?", be_user_id, u.ID).Find(&fans)
	if len(fans) > 0 {
		return true
	}
	return false
}

// 查看是否互相关注
func (u *User) QueryFansMutualed(be_user_id int64) bool {
	db := mysql.GetMysqlDB()
	query := db.Table("user_fans")
	var fans []UserFans
	query = query.Select("be_user_id = ? AND user_id = ? AND be_user_id = ? AND user_id = ?", be_user_id, u.ID, u.ID, be_user_id).Find(&fans)
	if len(fans) > 0 {
		return true
	}
	return false
}

// 新增用户
func (u *User) NewUser() error {
	db := mysql.GetMysqlDB()
	return db.Create(&u).Error
}

// 收藏表情包
func (u *UserEmoticons) AddUserEmoticons() error {
	db := mysql.GetMysqlDB()
	return db.Create(&u).Error
}

// 是否已收藏
func (u *UserEmoticons) BoolUserEmoticons() error {
	db := mysql.GetMysqlDB()
	return db.Where("user_id = ? AND emoticons_id = ?", u.UserID, u.EmoticonsID).First(&u).Error
}
