package router

import (
	"fmt"

	"github.com/Biubiubiuuuu/goDoutu/controller"
	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/Biubiubiuuuu/goDoutu/middleware/cross"
	"github.com/Biubiubiuuuu/goDoutu/middleware/error"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func Init() *gin.Engine {
	// 设置模式
	if config.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// 静态资源路径 /static 开头 或者 取自定义配置
	router.Static(config.Static, fmt.Sprintf(".%s", config.Static))
	// 允许跨域请求
	router.Use(cross.Cors())
	// 小程序端
	InitUser(router)
	// 后台后续处理
	//404
	router.NoRoute(error.NotFound)
	return router
}

// 小程序端
func InitUser(router *gin.Engine) {
	api := router.Group("api/v1")
	// 微信授权
	api.POST("auth", controller.WechatAuth)
	// 新增或更新用户信息
	api.POST("user", controller.NewUser)
}
