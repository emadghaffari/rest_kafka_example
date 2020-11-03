package app

import (
	"os"

	"github.com/emadghaffari/res_errors/logger"
	"github.com/emadghaffari/rest_kafka_example/databases/elasticsearch"
	"github.com/emadghaffari/rest_kafka_example/databases/kafka"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func init()  {
	gin.SetMode(os.Getenv("GIN_MODE"))
}


// StartApplication func
func StartApplication() {
	elasticsearch.Init()
	consumer := kafka.Consumer{}
	go consumer.Consumer()
	mapURL()
	logger.Info("about to start application")
	router.Run(":8000")
}