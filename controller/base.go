package controller

import (
	"github.com/gin-gonic/gin"
)

// 返回结果
type DoutuResponse struct {
	Status  bool        `json:"status"`  // 成功失败标志；true：成功 、false：失败
	Data    interface{} `json:"data"`    // 返回数据
	Message string      `json:"message"` // 提示信息
}

// 请求返回
func Response(c *gin.Context, code int, data DoutuResponse) {
	c.JSON(code, data)
}
