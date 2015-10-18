package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	today    time.Time
	tomorrow time.Time
)

func init() {
	today = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	tomorrow = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.UTC)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Statban Server"))
}

func batchHandler(w http.ResponseWriter, req *http.Request) {
	db := StatbanConfig.Db
	res, err := db.GetBatchStats()
	if err != nil {
		internalError(w, "Error when reading batch stats", err)
		return
	}

	setHeaders(w)
	e := json.NewEncoder(w)
	e.Encode(res)
}

func dailyHandler(w http.ResponseWriter, req *http.Request) {
	db := StatbanConfig.Db
	res, err := db.GetDailyStats()
	if err != nil {
		internalError(w, "Error when reading daily stats", err)
		return
	}

	setHeaders(w)
	e := json.NewEncoder(w)
	e.Encode(res)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func fail(w http.ResponseWriter, status int, msg string, err error) {
	w.WriteHeader(status)
	err = fmt.Errorf("%s: %s", msg, err.Error())
	log.Printf(err.Error())
	w.Write([]byte(err.Error()))
}

func internalError(w http.ResponseWriter, msg string, err error) {
	fail(w, http.StatusInternalServerError, msg, err)
}

func clientError(w http.ResponseWriter, msg string, err error) {
	fail(w, http.StatusBadRequest, msg, err)
}
