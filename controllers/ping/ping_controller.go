package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping func for get status ping from service
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
