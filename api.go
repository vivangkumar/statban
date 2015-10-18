package main

import (
	"encoding/json"
	"fmt"
	r "github.com/dancannon/gorethink"
	s "github.com/vivangkumar/statban/stats"
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
	cur, err := r.DB(db.Name).Table("hourly_summary").Filter(r.Row.Field("created_at").
		During(today, tomorrow)).Run(db.Session)
	if err != nil {
		log.Printf("Error reading batch summary: %v", err.Error())
		internalError(w, "Read error", err)
		return
	}
	defer cur.Close()

	var res []s.SummarizedBatch
	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when decoding into struct: %v", err.Error())
		internalError(w, "Decoding error", err)
		return
	}

	setHeaders(w)
	e := json.NewEncoder(w)
	e.Encode(res)
}

func dailyHandler(w http.ResponseWriter, req *http.Request) {
	db := StatbanConfig.Db
	cur, err := r.DB(db.Name).Table("daily_summary").Filter(r.Row.Field("beginning").
		Eq(today)).Run(db.Session)
	if err != nil {
		log.Printf("Error reading day summary: %v", err.Error())
		internalError(w, "Read error", err)
		return
	}
	defer cur.Close()

	var res []s.SummarizedDay
	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when decoding into struct: %v", err.Error())
		internalError(w, "Decoding error", err)
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
