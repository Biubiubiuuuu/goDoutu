package models

import (
	"strings"

	"github.com/Biubiubiuuuu/goDoutu/db/mysql"
)

// 表情包
type Emoticons struct {
	Model
	Url             string `gorm:"size:200" json:"url"`              // 表情包链接地址
	WordDescription string `gorm:"size:255" json:"word_description"` // 文字描述
	EmoticonsTypeID int64  `json:"emoticons_type_id"`                // 表情包类型ID
	ContributorID   int64  `json:"contributor_id"`                   // 贡献者ID
	GroupingID      int64  `json:"crouping_id"`                      // 表情包图组ID
}

// 表情包类型
type EmoticonsType struct {
	ID   int64  `gorm:"primary_key" json:"id"`
	Name string `gorm:"size:255" json:"name"` // 表情包图组标题
}

// 表情包图组
type EmoticonsGrouping struct {
	ID          int64  `gorm:"primary_key" json:"id"`
	Title       string `gorm:"size:255" json:"name"`        // 表情包类型名称
	Description string `gorm:"size:255" json:"description"` // 表情包类型描述
}

// 表情包
type EmoticonsInfo struct {
	Emoticons
	ContributorName     int64  `json:"contributor_name"`     // 贡献者ID
	EmoticonsTypeName   string `json:"emoticons_type_name"`  // 表情包类型名称
	GroupingTitle       string `json:"crouping_title"`       // 表情包图组标题
	GroupingDescription string `json:"crouping_description"` // 表情包图组描述
}

// 添加新表情包
func (e *Emoticons) AddEmoticons() error {
	db := mysql.GetMysqlDB()
	return db.Create(&e).Error
}

// 查询表情包详情 By ID
func (e *Emoticons) QueryEmoticonsByID() error {
	db := mysql.GetMysqlDB()
	return db.First(&e).Error
}

// 查询表情包
//  文字描述 word_description
//  表情包类型ID emoticons_type_id
//  表情包图组ID crouping_id
//  贡献者ID contributor_id
func QueryEmoticons(pageSize int, page int, args map[string]interface{}) (count int, emoticons []Emoticons) {
	db := mysql.GetMysqlDB()
	db.Table("emoticons")
	query := db.Select("DISTINCT emoticons.*,user.nickname AS ContributorName,emoticons_type.name AS emoticons_type_name,emoticons_grouping.title AS crouping_title,emoticons_grouping.description AS crouping_description")
	query = query.Joins("LEFT JOIN user ON user.id = emoticons.contributor_id")
	query = query.Joins("LEFT JOIN emoticons_type ON emoticons_type.id = emoticons.emoticons_type_ID")
	query = query.Joins("LEFT JOIN emoticons_grouping ON emoticons_grouping.id = emoticons.crouping_id")
	if v, ok := args["word_description"]; ok && v.(string) != "" {
		var buf strings.Builder
		buf.WriteString("%")
		buf.WriteString(v.(string))
		buf.WriteString("%")
		query = query.Where("emoticons.word_description like ?", buf.String())
	}
	if v, ok := args["emoticons_type_id"]; ok && v.(int64) > 0 {
		query = query.Where("emoticons.emoticons_type_id = ?", v.(int64))
	}
	if v, ok := args["contributor_id"]; ok && v.(int64) > 0 {
		query = query.Where("emoticons.contributor_id = ?", v.(int64))
	}
	if v, ok := args["crouping_id"]; ok && v.(int64) > 0 {
		query = query.Where("emoticons.crouping_id = ?", v.(int64))
	}
	query.Count(&count)
	query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&emoticons)
	return
}
