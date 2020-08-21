package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Biubiubiuuuu/goDoutu/models"
	"github.com/beevik/etree"
	"github.com/google/uuid"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

// 获取图片Channel
var ImageDataChannel = make(chan []string, 100)

const (
	AccessKey    = "VKYaoHZ9no66HILp2XmBMl4RwkvZNLX6F67ek3Qd"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         // 七牛云 access_key
	SecretKey    = "OhvUGMTl45_DykHfhjBMIpd3IKl9g_Qqae2PaXWI"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         // 七牛云 secret_key
	Bucketname   = "godoutu"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          // 七牛云空间名称
	NQYBasicUrl  = "http://qezg1i20h.hn-bkt.clouddn.com/"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             // 七牛云外链url
	MacPath      = "/Users/gaochao/LTWorks/goDoutu/webcrawler/zhihu/basic/image/"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     // 本机下载地址
	Url          = "https://www.zhihu.com/api/v4/questions/%v/answers?include=data[*].is_normal,admin_closed_comment,reward_info,is_collapsed,annotation_action,annotation_detail,collapse_reason,is_sticky,collapsed_by,suggest_edit,comment_count,can_comment,content,editable_content,voteup_count,reshipment_settings,comment_permission,created_time,updated_time,review_info,relevant_info,question,excerpt,relationship.is_authorized,is_author,voting,is_thanked,is_nothelp,is_labeled,is_recognized,paid_info,paid_info_content;data[*].mark_infos[*].url;data[*].author.follower_count,badge[*].topics&offset=%d&limit=%d&sort_by=default&platform=desktop" // 知乎问题api地址
	Limit        = 5                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  // 知乎问题返回问题数 最大20条
	Question     = "310564833"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        // 知乎问题ID
	ClientID     = "iHCdpZwXdG7TXb5jld7MzvTa"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         // 百度API client_id
	ClientSecret = "xK6R7zHLLhaTonRzNTTBkwhoLFjVkETr"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 // 百度api client_secret
)

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

// 百度AL api 授权返回结构体
type BaiDuALTokenResponse struct {
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	SessionSecret    string `json:"session_secret"`
	ErrorDescription string `json:"error_description"`
}

// 百度AL 图片文字识别返回结构体
type BaiDuALPictureIdentifyResponse struct {
	LogID          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
	ErrorMsg string `json:"error_msg"`
}

func main() {
	go DownloadImgToNiuqiyun()
	// 问题ID 跳转页 页大小
	url := fmt.Sprintf(Url, Question, 0, Limit)
	response := HttpGet(url)
	var data ResponseData
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		panic(err)
	}
	//总回答数
	count := data.Paging.Totals
	fmt.Println(count)
	//总页大小 向上取整
	pages := int(math.Ceil(float64(count) / float64(Limit)))
	var arr []string
	fmt.Println(pages)
	for i := 0; i < pages; i++ {
		url := fmt.Sprintf(Url, Question, i, Limit)
		response := GetResponseData(url)
		for _, dataV := range response.Data {
			doc := etree.NewDocument()
			if err := doc.ReadFromString(dataV.Content); err != nil {
				panic(err)
			}
			// 遍历是否存在重复答案
			if len(arr) == 0 {
				arr = append(arr, dataV.Author.ID)
			}
			if !ExistStr(arr, dataV.Author.ID) || dataV.Author.ID == "0" {
				arr = append(arr, dataV.Author.ID)
				var urls []string
				for _, e := range doc.FindElements("./figure/noscript/img") {
					arrtValue := e.SelectAttrValue("src", "")
					// 遍历是否存在重复图片
					if len(urls) == 0 {
						urls = append(urls, arrtValue)
					}
					if !ExistStr(urls, arrtValue) {
						urls = append(urls, arrtValue)
					}
				}
				ImageDataChannel <- urls
			}
		}
	}
}

// 下载图片到七牛云
func DownloadImgToNiuqiyun() {
	for url := range ImageDataChannel {
		for _, v := range url {
			if v != "" {
				filePath := Download(v)
				if filePath != "" {
					_, nqyUrl := Upload(filePath)
					emoticons := models.Emoticons{
						WordDescription: BaiDuALPictureIdentify(nqyUrl),
						Url:             nqyUrl,
					}
					emoticons.AddEmoticons()
				}
			}
		}
	}
	defer close(ImageDataChannel)
}

// 循环请求获取图片
func GetResponseData(url string) (data ResponseData) {
	response := HttpGet(url)
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		panic(err)
	}
	return
}

// 七牛云获取token
func Token() string {
	putPolicy := storage.PutPolicy{
		Scope: Bucketname,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	return putPolicy.UploadToken(mac)
}

// 七牛云上传图片
func Upload(imgUrl string) (err error, url string) {
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	upToken := Token()
	uuidKey := uuid.New()
	str := strings.Split(imgUrl, ".")
	key := uuidKey.String() + "." + str[len(str)-1]
	if err := formUploader.PutFile(context.Background(), &ret, upToken, key, imgUrl, nil); err != nil {
		return err, ""
	}
	return nil, NQYBasicUrl + ret.Key
}

// 百度API获取授权
func BaiDuALToken() string {
	tokenResStr := HttpGet(fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%v&client_secret=%v", ClientID, ClientSecret))
	var res BaiDuALTokenResponse
	if err := json.Unmarshal([]byte(tokenResStr), &res); err != nil {
		panic(err)
	}
	if res.AccessToken != "" {
		return res.AccessToken
	}
	panic(res.ErrorDescription)
}

// 百度API文字识别
func BaiDuALPictureIdentify(imageUrl string) string {
	access_token := BaiDuALToken()
	reqUrl := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%v", access_token)
	postData := fmt.Sprintf("url=%v", imageUrl)
	resJson := HttpPost(reqUrl, postData)
	var res BaiDuALPictureIdentifyResponse
	if err := json.Unmarshal([]byte(resJson), &res); err != nil {
		panic(err)
	}
	if len(res.WordsResult) > 0 {
		var str []string
		for _, v := range res.WordsResult {
			str = append(str, v.Words)
		}
		return strings.Join(str, ",")
	}
	return ""
}

// 指定的字符串是否存在数组中
func ExistStr(array []string, str string) bool {
	exist := false
	for _, v := range array {
		if v == str {
			exist = true
		}
	}
	return exist
}

// 下载图片到本地
func Download(url string) string {
	str := strings.Split(url, ".")
	fileName := uuid.New().String() + "." + strings.Replace(str[len(str)-1], "?source=1940ef5c", "", -1)
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	time := time.Now().Format("20060102")
	path := MacPath + time
	if !IsExist(path) {
		CreateDir(path)
	}
	path = path + "/" + fileName
	file, err := os.Create(path)
	if err != nil {
		return ""
	}
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	if GetFileSize(path) == 0 {
		return ""
	}
	return path
}

// 文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// 判断文件夹/文件是否存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// 判断文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func HttpGet(url string) string {
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

// 发送POST请求
// url：         请求地址
// response：    请求返回的内容
func HttpPost(url string, postData string) string {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
