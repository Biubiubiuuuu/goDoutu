package controller

import (
	"bytes"
	"io"
	"net/http"
	"time"

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
