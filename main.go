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
	http.HandleFunc("/batch", batchHandler)
	http.HandleFunc("/daily", dailyHandler)

	log.Printf("Statban server running on %v", StatbanConfig.HttpAddress)
	http.ListenAndServe(StatbanConfig.HttpAddress, nil)
}
