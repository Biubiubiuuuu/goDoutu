package error

import (
	"net/http"

	"github.com/Biubiubiuuuu/goDoutu/controller"
	"github.com/gin-gonic/gin"
)

// 404
func NotFound(c *gin.Context) {
	response := controller.DoutuResponse{
		Message: "404 Not Found",
	}
	c.JSON(http.StatusNotFound, response)
}
