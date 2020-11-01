package app

import (
	"github.com/emadghaffari/rest_kafka_example/controllers/ping"
	"github.com/emadghaffari/rest_kafka_example/controllers/twitter"
)

func mapURL() {
	router.GET("/ping", ping.Ping)

	router.POST("/search", twitter.Search)
	router.POST("/store", twitter.Store)
}
