package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bmonkman/PineappleScope/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const newFiringThreshold = 3 // hour threshold for new measurements to be considered part of a firing

// TemperatureHandlers handles routing for temperature endpoints
type TemperatureHandlers struct {
	renderer *gin.Engine
}

// NewTemperatureHandlers returns an instance of TemperatureHandlers
func NewTemperatureHandlers(r *gin.Engine) *TemperatureHandlers {
	return &TemperatureHandlers{renderer: r}
}

// Register firing handlers
func (f *TemperatureHandlers) Register() {
	f.renderer.POST("/temperature", createTemperatureReading)
}

// Create a new temperature reading
func createTemperatureReading(c *gin.Context) {
	innerTemperature, err := strconv.ParseFloat(c.PostForm("inner"), 64)
	if err != nil {
		panic(err)
	}

	outerTemperature, err := strconv.ParseFloat(c.PostForm("outer"), 64)
	if err != nil {
		panic(err)
	}

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	firingID := getOrCreateFiring(db)

	newReading := models.TemperatureReading{FiringID: firingID, Inner: innerTemperature, Outer: outerTemperature}
	db.Save(&newReading)

	c.Status(http.StatusOK)
}

// If there was a recent temperature reading, return that reading's session, otherwise create a new session
func getOrCreateFiring(db *gorm.DB) uint {
	var temperature = models.TemperatureReading{}

	whereString := fmt.Sprintf("created_date >= datetime('now', 'localtime', '-%d hours')", newFiringThreshold)
	db.Where(whereString).
		Order("created_date desc").
		First(&temperature)

	var firing models.Firing
	if temperature.ID == 0 || temperature.FiringID == 0 {
		firing = models.Firing{StartDate: time.Now(), EndDate: time.Now(), StartDateAmbientTemp: temperature.Outer, Name: "New Firing"}
		db.Save(&firing)
	} else {
		db.First(&firing, temperature.FiringID)
		firing.EndDate = time.Now()

		db.Save(&firing)
	}

	return firing.ID
}
