package app

import (
	"github.com/emadghaffari/res_errors/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApplication func
func StartApplication() {
	mapURL()
	logger.Info("about to start application")
	router.Run(":8080")
}