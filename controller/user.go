package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/Biubiubiuuuu/goDoutu/models"
	"github.com/gin-gonic/gin"
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
	respData := HttpGet(url)
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
		resp.Message = "请求参数数据格式不合法，请检查，格式为JSON！"
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
	resp.Message = "用户注册成功"
	Response(c, http.StatusOK, resp)
	return
}
