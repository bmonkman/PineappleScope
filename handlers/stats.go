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
		fmt.Println("Couldn't read free memory")
		panic(err)
	}

	uptime, err := strconv.ParseUint(c.PostForm("uptime"), 10, 64)
	if err != nil {
		fmt.Println("Couldn't read uptime")
		panic(err)
	}

	temperature, err := strconv.ParseFloat(c.PostForm("temp"), 64)
	if err != nil {
		fmt.Println("Couldn't read temp")
		panic(err)
	}

	cpuTemp, err := strconv.ParseFloat(c.PostForm("cpuTemp"), 64)
	if err != nil {
		fmt.Println("Couldn't read cpu temp")
		panic(err)
	}

	ambientTemp, err := strconv.ParseFloat(c.PostForm("ambientTemp"), 64)
	if err != nil {
		fmt.Println("Couldn't read ambient temp")
		panic(err)
	}

	humidity, err := strconv.ParseFloat(c.PostForm("humidity"), 64)
	if err != nil {
		fmt.Println("Couldn't read humidity")
		panic(err)
	}

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	newStats := models.Stats{
		FreeMemory:         freeMemory,
		Uptime:             uptime,
		Temperature:        temperature,
		AmbientTemperature: ambientTemp,
		CPUTemperature:     cpuTemp,
		Humidity:           humidity}

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
	db.Limit(60 * 5).Order("created_date desc").Find(&stats)

	Success(c, "stats", gin.H{
		"title": "Stats",
		"stats": stats})

}
