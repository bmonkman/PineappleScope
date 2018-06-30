package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bmonkman/PineappleScope/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// FiringHandlers handles routing for firing endpoints
type FiringHandlers struct {
	renderer *gin.Engine
}

// NewFiringHandlers returns an instance of FiringHandlers
func NewFiringHandlers(r *gin.Engine) *FiringHandlers {
	return &FiringHandlers{renderer: r}
}

// Register firing handlers
func (f *FiringHandlers) Register() {
	f.renderer.GET("/firing/:firingId", f.getFiring)
	f.renderer.GET("/firing/:firingId/edit", f.showEditFiring)
	f.renderer.POST("/firing/:firingId", f.editFiring)
	f.renderer.GET("/firings", f.listFirings)
	f.renderer.GET("/firing/:firingId/readings", f.getReadingsForFiring)
	f.renderer.DELETE("/firing/:firingId", f.deleteFiring)
}

// Get details about a specific firing
func (f *FiringHandlers) getFiring(c *gin.Context) {
	firingID := c.Param("firingId")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var firing models.Firing
	db.First(&firing, firingID)

	var temperatureReadings []models.TemperatureReading
	db.Where("firing_id = ?", firingID).Find(&temperatureReadings)

	var peakTemperature = 0.0
	for _, temp := range temperatureReadings {
		if temp.Inner > peakTemperature {
			peakTemperature = temp.Inner
		}
	}

	Success(c, "firing", gin.H{
		"title":               "Firing: " + firing.Name,
		"firing":              firing,
		"temperatureReadings": temperatureReadings,
		"peakTemperature":     peakTemperature})

}

// Show form to edit details about a specific firing
func (f *FiringHandlers) showEditFiring(c *gin.Context) {
	firingID := c.Param("firingId")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var firing models.Firing
	db.First(&firing, firingID)

	if firing.Name == "New Firing" {
		firing.Name = ""
	}
	Success(c, "new-firing", gin.H{
		"title":  "Edit Firing: " + firing.Name,
		"firing": firing})
}

// Edit details about a specific firing
func (f *FiringHandlers) editFiring(c *gin.Context) {
	firingID := c.Param("firingId")
	name := c.PostForm("name")
	coneNumber := c.PostForm("coneNumber")
	notes := c.PostForm("notes")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	if name == "" {
		name = "New Firing"
	}

	var firing models.Firing
	db.Where("id = ?", firingID).Find(&firing)
	db.First(&firing, firingID)
	firing.Name = name
	firing.Notes = notes
	firing.ConeNumber = coneNumber
	fmt.Println(firing)
	db.Save(&firing)

	c.Redirect(http.StatusMovedPermanently, "/firing/"+firingID)
}

// List all firings
func (f *FiringHandlers) listFirings(c *gin.Context) {
	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var firings []models.Firing
	db.Order("end_date DESC").Find(&firings)

	Success(c, "list", gin.H{
		"title":                  "Firings",
		"firings":                firings,
		"currentFiringThreshold": time.Now().Add(-3 * time.Hour)})

}

// List all readings for a firing
func (f *FiringHandlers) getReadingsForFiring(c *gin.Context) {
	firingID := c.Param("firingId")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var temperatureReadings []models.TemperatureReading
	db.Where("firing_id = ?", firingID).Order("created_date").Find(&temperatureReadings)

	c.JSON(http.StatusOK, temperatureReadings)

}

// Delete a firing
func (f *FiringHandlers) deleteFiring(c *gin.Context) {
	firingID := c.Param("firingId")

	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	db.Delete(models.Firing{}, "ID = ?", firingID)
	db.Delete(models.TemperatureReading{}, "firing_id = ?", firingID)
	db.Delete(models.Photo{}, "firing_id = ?", firingID)
	c.Status(http.StatusOK)
}
