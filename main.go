package main

import (
	s "github.com/vivangkumar/statban/stats"
	"log"
	"net/http"
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

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hourly", hourlyHandler)
	http.HandleFunc("/daily", dailyHandler)

	log.Printf("Statban server running on %v", StatbanConfig.Port)
	http.ListenAndServe(StatbanConfig.Port, nil)
}
