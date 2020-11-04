package app

import (
	"github.com/emadghaffari/res_errors/logger"
	"github.com/emadghaffari/rest_kafka_example/databases/elasticsearch"
	"github.com/emadghaffari/rest_kafka_example/databases/kafka"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	router = gin.Default()
)

func init()  {
	gin.SetMode(viper.GetString("gin.mode"))
}


// StartApplication func
func StartApplication() {
	elasticsearch.Init()
	kafka.Init()
	consumer := kafka.Consumer{}
	go consumer.Consumer()
	mapURL()
	logger.Info("about to start application")
	router.Run(viper.GetString("gin.port"))
}