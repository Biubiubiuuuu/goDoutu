package baidu

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/Biubiubiuuuu/goDoutu/helper/request"
)

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

// 百度API获取授权
func BaiDuALToken() string {
	tokenResStr := request.HttpGet(fmt.Sprintf("%v/oauth/2.0/token?grant_type=client_credentials&client_id=%v&client_secret=%v", config.BAIDUUrl, config.BAIDUClientID, config.BAIDUClientSecret))
	var res BaiDuALTokenResponse
	if err := json.Unmarshal([]byte(tokenResStr), &res); err != nil {
		log.Println("Baidu AL auth fail:" + err.Error())
		return ""
	}
	if res.AccessToken != "" {
		return res.AccessToken
	}
	log.Println("Baidu AL auth fail:" + res.ErrorDescription)
	return ""
}

// 百度API文字识别
func BaiDuALPictureIdentify(imageUrl string) string {
	access_token := BaiDuALToken()
	reqUrl := fmt.Sprintf("%v/rest/2.0/ocr/v1/general_basic?access_token=%v", config.BAIDUUrl, access_token)
	postData := fmt.Sprintf("url=%v", imageUrl)
	resJson := request.HttpPost(reqUrl, postData)
	var res BaiDuALPictureIdentifyResponse
	if err := json.Unmarshal([]byte(resJson), &res); err != nil {
		log.Println("Baidu AL fail:" + err.Error())
		return ""
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
