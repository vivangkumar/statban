package main

import (
	"github.com/gin-gonic/gin"
	s "github.com/vivangkumar/statban/stats"
)

func init() {
	initialize()
}

func main() {
	cfg, err := StatbanConfig.Db.Setup()
	if err != nil {
		panic(err.Error())
	}

	go s.RunCollector(cfg, GithubConfig)

	router := gin.Default()
	router.LoadHTMLGlob("public/*")

	router.GET("/", rootHandler)
	router.GET("/hourly", hourlyHandler)
	router.GET("/daily", dailyHandler)

	router.GET("/graphs", graphHandler)

	router.Run(StatbanConfig.Port)
}
