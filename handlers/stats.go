package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmonkman/PineappleScope/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// StatsHandlers handles routing for stats endpoints
type StatsHandlers struct {
	renderer *gin.Engine
}

// NewStatsHandlers returns an instance of StatsHandlers
func NewStatsHandlers(r *gin.Engine) *StatsHandlers {
	return &StatsHandlers{renderer: r}
}

// Register stats handlers
func (h *StatsHandlers) Register() {
	h.renderer.GET("/stats", h.getStats)
	h.renderer.POST("/stats", h.createStatsRecord)
}

// Create a new stats record
func (h *StatsHandlers) createStatsRecord(c *gin.Context) {

	freeMemory, err := strconv.ParseUint(c.PostForm("freeMemory"), 10, 64)
	if err != nil {
		panic(err)
	}

	uptime, err := strconv.ParseUint(c.PostForm("uptime"), 10, 64)
	if err != nil {
		panic(err)
	}

	event := c.PostForm("event")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	newStats := models.Stats{FreeMemory: freeMemory, Uptime: uptime, Event: event}
	fmt.Println(newStats)
	db.Save(&newStats)

	c.Status(http.StatusOK)
}

// List most recent stats entries
func (h *StatsHandlers) getStats(c *gin.Context) {
	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var stats []models.Stats
	db.Limit(50).Find(&stats)

	Success(c, "stats", gin.H{
		"title": "Stats",
		"stats": stats})

}
