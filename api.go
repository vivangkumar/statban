package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Statban": "API"})
}

func hourlyHandler(c *gin.Context) {
	res, err := StatbanConfig.Db.GetHourlyStats()
	if err != nil {
		internalError(c, "Error when reading hourly stats", err)
		return
	}

	sendResponse(c, res)
}

func dailyHandler(c *gin.Context) {
	res, err := StatbanConfig.Db.GetDailyStats()
	if err != nil {
		internalError(c, "Error when reading daily stats", err)
		return
	}

	sendResponse(c, res)
}

func graphHandler(c *gin.Context) {
	res, err := StatbanConfig.Db.GetDailyStats()
	if err != nil {
		internalError(c, "Error when reading daily stats", err)
		return
	}

	data, _ := json.Marshal(res)
	c.HTML(http.StatusOK, "graph.tmpl", gin.H{
		"title": "Kanban graph",
		"data":  string(data),
	})
}

func sendResponse(c *gin.Context, res interface{}) {
	setHeaders(c)
	c.JSON(http.StatusOK, res)
}

func setHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
}

func fail(c *gin.Context, status int, msg string, err error) {
	err = fmt.Errorf("%s: %s", msg, err.Error())
	log.Printf(err.Error())
	c.JSON(status, gin.H{"error": err.Error()})
}

func internalError(c *gin.Context, msg string, err error) {
	fail(c, http.StatusInternalServerError, msg, err)
}

func clientError(c *gin.Context, msg string, err error) {
	fail(c, http.StatusBadRequest, msg, err)
}
