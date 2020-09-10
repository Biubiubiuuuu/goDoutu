package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/Biubiubiuuuu/goDoutu/helper/file"
	fil "github.com/Biubiubiuuuu/goDoutu/helper/file"
	"github.com/Biubiubiuuuu/goDoutu/helper/request"
	"github.com/Biubiubiuuuu/goDoutu/thirdparty/baidu"
	"github.com/Biubiubiuuuu/goDoutu/thirdparty/qiniuyun"

	"github.com/Biubiubiuuuu/goDoutu/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 微信授权参数
type WechatAuthReq struct {
	Code string `json:"code"` // 临时登录凭证
}

// 微信授权返回参数
type WechatAuthResponse struct {
	OpenID string `json:"open_id"` // 微信用户唯一标识
	Errmsg string `json:"errmsg"`  // 错误信息
}

// 微信更新用户信息
type NewUserReq struct {
	Sex      int64  `json:"sex"`      // 性别 0:未知 1:男 2:女
	Avatar   string `json:"avatar"`   // 头像
	Nickname string `json:"nickname"` // 微信用户昵称
	OpenID   string `json:"open_id"`  // 微信用户唯一标识
	Country  string `json:"country"`  // 所在国家
	Province string `json:"province"` // 所在省份
	City     string `json:"city"`     // 所在城市
}

// 表情包请求参数
type EmoticonsRequest struct {
	Url                 string `json:"url"`                  // 表情包链接地址
	EmoticonsTypeID     int64  `json:"emoticons_type_id"`    // 表情包类型ID
	GroupingTitle       string `json:"crouping_title"`       // 表情包图组标题
	GroupingDescription string `json:"crouping_description"` // 表情包图组描述
	OpenID              string `json:"open_id"`              // 微信用户唯一ID
}

// 微信授权
func WechatAuth(c *gin.Context) {
	var req WechatAuthReq
	var resp DoutuResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Message = fmt.Sprintf("请求参数错误:%s！数据格式不合法，请检查，格式为JSON", err.Error())
		Response(c, http.StatusOK, resp)
		return
	}
	if req.Code == "" {
		resp.Message = "临时登录凭证code不能为空！"
		Response(c, http.StatusOK, resp)
		return
	}
	// 微信登录
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", config.AppID, config.AppSecret, req.Code)
	respData := request.HttpGet(url)
	var authResp WechatAuthResponse
	if err := json.Unmarshal([]byte(respData), &authResp); err != nil {
		resp.Message = fmt.Sprintf("序列化错误:%s", err.Error())
		Response(c, http.StatusOK, resp)
		return
	}
	if authResp.OpenID == "" {
		resp.Message = fmt.Sprintf("授权失败:%s", authResp.Errmsg)
		Response(c, http.StatusOK, resp)
		return
	}
	resp.Status = true
	resp.Message = "授权成功"
	resp.Data = authResp.OpenID
	Response(c, http.StatusOK, resp)
	return
}

// 微信授权之后创建用户
func NewUser(c *gin.Context) {
	var resp DoutuResponse
	var req NewUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Message = fmt.Sprintf("请求参数错误:%s！数据格式不合法，请检查，格式为JSON", err.Error())
		Response(c, http.StatusOK, resp)
		return
	}

	if req.OpenID == "" {
		resp.Message = "open_id不能为空！"
		Response(c, http.StatusOK, resp)
		return
	}
	// 默认性别为 0:未知
	if req.Sex < 0 || req.Sex > 2 {
		req.Sex = 0
	}
	user := models.User{
		OpenID:   req.OpenID,
		Avatar:   req.Avatar,
		City:     req.City,
		Country:  req.Country,
		Sex:      req.Sex,
		Nickname: req.Nickname,
		Province: req.Province,
	}
	// 不存在,则直接新增用户
	if err := user.QueryUserByOpenID(); err != nil {
		if err := user.NewUser(); err != nil {
			resp.Message = fmt.Sprintf("用户注册失败：%s", err.Error())
			Response(c, http.StatusOK, resp)
			return
		}
		user.QueryUserByOpenID()
		resp.Status = true
		resp.Message = "用户注册成功"
		Response(c, http.StatusOK, resp)
		return
	}
	// 存在，更新用户信息
	args := map[string]interface{}{
		"sex":      req.Sex,
		"avatar":   req.Avatar,
		"nickname": req.Nickname,
		"country":  req.Country,
		"province": req.Province,
		"city":     req.City,
	}
	// 更新用户信息，不考虑是否成功
	user.UpdatesUserByID(args)
	resp.Status = true
	resp.Message = "用户更新成功"
	Response(c, http.StatusOK, resp)
	return
}

// 获取表情包
func Emoticons(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "25"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	emoticons_type_id, _ := strconv.ParseInt(c.Query("emoticons_type_id"), 10, 64)
	crouping_id, _ := strconv.ParseInt(c.Query("crouping_id"), 10, 64)
	contributor_id, _ := strconv.ParseInt(c.Query("contributor_id"), 10, 64)
	args := map[string]interface{}{
		"word_description":  c.Query("word_description"),
		"emoticons_type_id": emoticons_type_id,
		"crouping_id":       crouping_id,
		"contributor_id":    contributor_id,
		"sort_condition":    c.DefaultQuery("sort_condition", "views"),
	}
	count, emoticons := models.QueryEmoticons(pageSize, page, args)
	var resp DoutuResponse
	resp.Message = "查询成功"
	if len(emoticons) == 0 {
		resp.Message = "没有更多了～"

	}
	resp.Status = true
	resp.Data = map[string]interface{}{
		"emoticons": emoticons,
		"count":     count,
	}
	Response(c, http.StatusOK, resp)
	return
}

// 上传表情包图片接口
func UploadImage(c *gin.Context) {
	var resp DoutuResponse
	// 获取主机头
	r := c.Request
	host := r.Host
	if strings.HasPrefix(host, "http://") == false {
		host = "http://" + host
	}
	var image string
	if file, err := c.FormFile("image"); err == nil {
		// 文件名 避免重复取uuid
		var filename string
		uuid, _ := uuid.NewUUID()
		arr := strings.Split(file.Filename, ".")
		if strings.EqualFold(arr[len(arr)-1], "png") {
			filename = uuid.String() + ".png"
		} else if strings.EqualFold(arr[len(arr)-1], "jpg") {
			filename = uuid.String() + ".jpg"
		} else if strings.EqualFold(arr[len(arr)-1], "jpeg") {
			filename = uuid.String() + ".jpeg"
		} else if strings.EqualFold(arr[len(arr)-1], "gif") {
			filename = uuid.String() + ".gif"
		} else {
			resp.Message = "图片格式只支持png、jpg、jpeg、gif"

		}
		pathFile := config.ImageDir
		if !fil.IsExist(pathFile) {
			fil.CreateDir(pathFile)
		}
		pathFile = pathFile + filename
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			image = host + "/" + pathFile
			image = pathFile
		}
	}
	if image == "" {
		resp.Message = "图片上传失败"
		Response(c, http.StatusOK, resp)
		return
	}
	resp.Status = true
	resp.Message = "图片上传成功"
	resp.Data = image
	Response(c, http.StatusOK, resp)
	return
}

// 删除已上传表情包图片
func DelImage(c *gin.Context) {
	file := c.Query("file")
	var resp DoutuResponse
	if file == "" {
		resp.Message = "图片路径不能为空"
		Response(c, http.StatusOK, resp)
		return
	}
	if err := os.Remove(file); err != nil {
		resp.Message = "图片删除失败"
		Response(c, http.StatusOK, resp)
		return
	}
	resp.Status = true
	resp.Message = "图片删除成功"
	Response(c, http.StatusOK, resp)
	return
}

// 用户发布表情包
func NewEmoticons(c *gin.Context) {
	var resp DoutuResponse
	var req EmoticonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Message = fmt.Sprintf("请求参数错误:%s！数据格式不合法，请检查，格式为JSON", err.Error())
		Response(c, http.StatusOK, resp)
		return
	}
	var t models.EmoticonsType
	t.ID = req.EmoticonsTypeID
	if err := t.QueryEmoticonsTypeByID(); err != nil && req.EmoticonsTypeID != 0 {
		resp.Message = fmt.Sprintf("表情包类型不存在，请重新选择：%v", err.Error())
		Response(c, http.StatusOK, resp)
		return
	}
	if req.Url == "" {
		resp.Message = fmt.Sprintf("图片地址不能为空")
		Response(c, http.StatusOK, resp)
		return
	}
	if !file.IsExist(req.Url) {
		resp.Message = "图片不存在，请确认图片地址路径，或重新上传"
		Response(c, http.StatusOK, resp)
		return
	}
	if req.OpenID == "" {
		resp.Message = "open_id不能为空"
		Response(c, http.StatusOK, resp)
		return
	}
	user := models.User{
		OpenID: req.OpenID,
	}
	// 尚未注册的用户，需要先授权注册
	if err := user.QueryUserByOpenID(); err != nil {
		resp.Message = "用户信息不存在，请先授权完成注册，在上传"
		Response(c, http.StatusOK, resp)
		return
	}
	// 百度AL识别文字图片
	wordDescription := baidu.BaiDuALPictureIdentify(req.Url)
	// 先把图片上传至七牛云，再删除本地文件
	_, nqyUrl := qiniuyun.Upload(req.Url)
	os.Remove(req.Url)
	emoticons := models.Emoticons{
		Url:             nqyUrl,
		WordDescription: wordDescription,
		EmoticonsTypeID: req.EmoticonsTypeID,
	}
	grouping := models.EmoticonsGrouping{
		Title:       req.GroupingTitle,
		Description: req.GroupingTitle,
	}
	if err := models.UserAddEmoticons(emoticons, grouping); err != nil {
		resp.Message = "表情包发布失败"
		Response(c, http.StatusOK, resp)
		return
	}
	resp.Status = true
	resp.Message = "表情包发布成功"
	Response(c, http.StatusOK, resp)
	return
}
