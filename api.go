package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Statban Server"))
}

func hourlyHandler(w http.ResponseWriter, r *http.Request) {

}

func dailyHandler(w http.ResponseWriter, r *http.Request) {

}
