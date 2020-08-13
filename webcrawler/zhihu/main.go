package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

// 获取图片Channel
var ImageDataChannel = make(chan ImageData, 100)

const (
	Url      = "https://www.zhihu.com/api/v4/questions/%v/answers?include=data[*].is_normal,admin_closed_comment,reward_info,is_collapsed,annotation_action,annotation_detail,collapse_reason,is_sticky,collapsed_by,suggest_edit,comment_count,can_comment,content,editable_content,voteup_count,reshipment_settings,comment_permission,created_time,updated_time,review_info,relevant_info,question,excerpt,relationship.is_authorized,is_author,voting,is_thanked,is_nothelp,is_labeled,is_recognized,paid_info,paid_info_content;data[*].mark_infos[*].url;data[*].author.follower_count,badge[*].topics&offset=%d&limit=%d&sort_by=default&platform=desktop"
	Female   = 1
	Male     = 2
	Limit    = 20
	Question = "310564833"
)

type ImageData struct {
	Urls     []string `json:"urls"`      // 图片地址
	AuthorID string   `json:"author_id"` // 作者ID
}

// 知乎问题回答返回结构体
type ResponseData struct {
	Data []struct {
		ID         int    `json:"id"`
		Type       string `json:"type"`
		AnswerType string `json:"answer_type"`
		Question   struct {
			Type         string `json:"type"`
			ID           int    `json:"id"`
			Title        string `json:"title"`
			QuestionType string `json:"question_type"`
			Created      int    `json:"created"`
			UpdatedTime  int    `json:"updated_time"`
			URL          string `json:"url"`
			Relationship struct {
			} `json:"relationship"`
		} `json:"question"`
		Author struct {
			ID                string        `json:"id"`
			URLToken          string        `json:"url_token"`
			Name              string        `json:"name"`
			AvatarURL         string        `json:"avatar_url"`
			AvatarURLTemplate string        `json:"avatar_url_template"`
			IsOrg             bool          `json:"is_org"`
			Type              string        `json:"type"`
			URL               string        `json:"url"`
			UserType          string        `json:"user_type"`
			Headline          string        `json:"headline"`
			Badge             []interface{} `json:"badge"`
			BadgeV2           struct {
				Title        string        `json:"title"`
				MergedBadges []interface{} `json:"merged_badges"`
				DetailBadges []interface{} `json:"detail_badges"`
				Icon         string        `json:"icon"`
				NightIcon    string        `json:"night_icon"`
			} `json:"badge_v2"`
			Gender        int  `json:"gender"`
			IsAdvertiser  bool `json:"is_advertiser"`
			FollowerCount int  `json:"follower_count"`
			IsFollowed    bool `json:"is_followed"`
			IsPrivacy     bool `json:"is_privacy"`
		} `json:"author"`
		URL                string `json:"url"`
		IsCollapsed        bool   `json:"is_collapsed"`
		CreatedTime        int    `json:"created_time"`
		UpdatedTime        int    `json:"updated_time"`
		Extras             string `json:"extras"`
		IsCopyable         bool   `json:"is_copyable"`
		IsNormal           bool   `json:"is_normal"`
		VoteupCount        int    `json:"voteup_count"`
		CommentCount       int    `json:"comment_count"`
		IsSticky           bool   `json:"is_sticky"`
		AdminClosedComment bool   `json:"admin_closed_comment"`
		CommentPermission  string `json:"comment_permission"`
		CanComment         struct {
			Reason string `json:"reason"`
			Status bool   `json:"status"`
		} `json:"can_comment"`
		ReshipmentSettings string        `json:"reshipment_settings"`
		Content            string        `json:"content"`
		EditableContent    string        `json:"editable_content"`
		Excerpt            string        `json:"excerpt"`
		CollapsedBy        string        `json:"collapsed_by"`
		CollapseReason     string        `json:"collapse_reason"`
		AnnotationAction   interface{}   `json:"annotation_action"`
		MarkInfos          []interface{} `json:"mark_infos"`
		RelevantInfo       struct {
			IsRelevant   bool   `json:"is_relevant"`
			RelevantType string `json:"relevant_type"`
			RelevantText string `json:"relevant_text"`
		} `json:"relevant_info"`
		SuggestEdit struct {
			Reason          string `json:"reason"`
			Status          bool   `json:"status"`
			Tip             string `json:"tip"`
			Title           string `json:"title"`
			UnnormalDetails struct {
				Status      string `json:"status"`
				Description string `json:"description"`
				Reason      string `json:"reason"`
				ReasonID    int    `json:"reason_id"`
				Note        string `json:"note"`
			} `json:"unnormal_details"`
			URL string `json:"url"`
		} `json:"suggest_edit"`
		IsLabeled  bool `json:"is_labeled"`
		RewardInfo struct {
			CanOpenReward     bool   `json:"can_open_reward"`
			IsRewardable      bool   `json:"is_rewardable"`
			RewardMemberCount int    `json:"reward_member_count"`
			RewardTotalMoney  int    `json:"reward_total_money"`
			Tagline           string `json:"tagline"`
		} `json:"reward_info"`
		Relationship struct {
			IsAuthor         bool          `json:"is_author"`
			IsAuthorized     bool          `json:"is_authorized"`
			IsNothelp        bool          `json:"is_nothelp"`
			IsThanked        bool          `json:"is_thanked"`
			IsRecognized     bool          `json:"is_recognized"`
			Voting           int           `json:"voting"`
			UpvotedFollowees []interface{} `json:"upvoted_followees"`
		} `json:"relationship"`
		AdAnswer interface{} `json:"ad_answer"`
	} `json:"data"`
	Paging struct {
		IsEnd    bool   `json:"is_end"`
		IsStart  bool   `json:"is_start"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Totals   int    `json:"totals"`
	} `json:"paging"`
}

func main() {
	go DownloadImgToNiuqiyun()
	// 问题ID 跳转页 页大小
	url := fmt.Sprintf(Url, Question, 0, Limit)
	response := Get(url)
	var data ResponseData
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		panic(err)
	}
	//总回答数
	count := data.Paging.Totals
	//总页大小 向上取整
	pages := int(math.Ceil(float64(count) / float64(Limit)))
	var arr []string
	go func() {
		for i := 0; i < pages; i++ {
			if i < 1 {
				url := fmt.Sprintf(Url, Question, i, Limit)
				response := GetResponseData(url)
				for _, v := range response.Data {
					doc := etree.NewDocument()
					if err := doc.ReadFromString(v.Content); err != nil {
						panic(err)
					}
					imaData := ImageData{
						AuthorID: v.Author.ID,
					}
					// 遍历是否存在重复答案
					if len(arr) == 0 {
						arr = append(arr, v.Author.ID)
					}
					for _, v2 := range arr {
						if v2 != imaData.AuthorID {
							arr = append(arr, v.Author.ID)
							for _, e := range doc.FindElements("./figure/noscript/img") {
								imaData.Urls = append(imaData.Urls, e.SelectAttrValue("src", ""))
							}
							ImageDataChannel <- imaData
						}
					}
				}
			}
		}
	}()
	time.Sleep(time.Second * 3)
}

// 下载图片到七牛云
func DownloadImgToNiuqiyun() {
	for data := range ImageDataChannel {
		for _, url := range data.Urls {
			if url != "" {
				fmt.Println(url)
			}
		}
	}
	defer close(ImageDataChannel)
}

// 循环请求获取图片
func GetResponseData(url string) (data ResponseData) {
	response := Get(url)
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		panic(err)
	}
	return
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}
