package models

import "time"

// 用户
type User struct {
	Model
	Telephone      string    `gorm:"unique;size:50;" json:"telephone"` // 手机号
	Username       string    `gorm:"unique;size:200;" json:"username"` // 用户名
	Password       string    `json:"-"`                                // 密码
	Birthday       time.Time `json:"birthday"`                         // 生日
	Sex            int64     `json:"sex"`                              // 性别 0:未知 1:男 2:女
	Avatar         string    `gorm:"size:200" json:"avatar"`           // 头像
	Email          string    `gorm:"unique;size:50;" json:"email"`     // 邮箱
	QQ             string    `gorm:"size:200" json:"QQ"`               // QQ
	WeChatNickname string    `gorm:"size:150" json:"WeChat_nickname"`  // 微信昵称
}
