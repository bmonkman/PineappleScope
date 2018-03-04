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
	f.renderer.GET("/firings", f.getFirings)
	f.renderer.GET("/firing/:firingId/readings", f.getReadingsForFiring)
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

	c.HTML(http.StatusOK, "firing", gin.H{
		"title":               "Firing: " + firing.Name,
		"firing":              firing,
		"temperatureReadings": temperatureReadings})

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

	c.HTML(http.StatusOK, "new-firing", gin.H{
		"title":  "Edit Firing: " + firing.Name,
		"firing": firing})

}

// Edit details about a specific firing
func (f *FiringHandlers) editFiring(c *gin.Context) {
	firingID := c.Param("firingId")
	name := c.PostForm("name")
	notes := c.PostForm("notes")
	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var firing models.Firing
	db.Where("id = ?", firingID).Find(&firing)
	db.First(&firing, firingID)
	firing.Name = name
	firing.Notes = notes
	fmt.Println(firing)
	db.Save(&firing)

	c.Redirect(301, "/firing/"+firingID)
}

// List all firings
func (f *FiringHandlers) getFirings(c *gin.Context) {
	db, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		return
	}

	var firings []models.Firing
	db.Order("end_date DESC").Find(&firings)

	c.HTML(http.StatusOK, "list", gin.H{
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
	db.Where("firing_id = ?", firingID).Find(&temperatureReadings)

	c.JSON(200, temperatureReadings)

}
