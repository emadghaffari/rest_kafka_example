package main

import (
	"github.com/emadghaffari/rest_kafka_example/app"
	"github.com/joho/godotenv"
)

func init()  {
	err := godotenv.Load()
	if err != nil {
		panic("ENV file noy found")
	}
}

func main()  {
	app.StartApplication()
}
