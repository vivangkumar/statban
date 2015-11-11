package main

import (
	"github.com/gin-gonic/gin"
	s "github.com/vivangkumar/statban/stats"
	"log"
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

	router.GET("/", rootHandler)
	router.GET("/hourly", hourlyHandler)
	router.GET("/daily", dailyHandler)

	log.Printf("Statban server running on %v", StatbanConfig.Port)
	router.Run(StatbanConfig.Port)
}
