package app

import (
	"github.com/emadghaffari/rest_kafka_example/controllers/ping"
)

func mapURL() {
	router.GET("/ping", ping.Ping)
}
