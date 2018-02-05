package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-contrib/multitemplate"

	"github.com/gin-gonic/gin"
)

const newFiringThreshold = 3 // hour threshold for new measurements to be considered part of a firing

type Firing struct {
	ID                   uint
	StartDate            time.Time
	EndDate              time.Time
	StartDateAmbientTemp float64
	Cone                 uint
	Name                 string
	Notes                string

	TemperatureReadings []TemperatureReading
	Photos              []Photo
}

type TemperatureReading struct {
	ID          uint
	CreatedDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	FiringID    uint      `gorm:"index"`
	Inner       float64
	Outer       float64
}

type Photo struct {
	ID          uint
	FiringID    uint
	CreatedDate time.Time
	photoURL    string
}

// AddDbHandle middleware will add the db connection to the context
func AddDbHandle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
}

// GetOrCreateFiring If there was a recent temperature reading, return that reading's session, otherwise create a new session
func GetOrCreateFiring(db *gorm.DB) uint {
	var temperature = TemperatureReading{}

	whereString := fmt.Sprintf("created_date >= datetime('now')-%d", newFiringThreshold*60*60)
	db.Where(whereString).
		Order("created_date desc").
		First(&temperature)

	if temperature.ID == 0 {
		newFiring := Firing{StartDate: time.Now(), EndDate: time.Now(), Name: "New Firing"}
		db.Create(&newFiring)
		return newFiring.ID
	}

	return temperature.FiringID
}

func createMyRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("list", "resources/html/base.html", "resources/html/list.html")
	r.AddFromFiles("firing", "resources/html/base.html", "resources/html/firing.html")
	r.AddFromFiles("new-firing", "resources/html/base.html", "resources/html/new-firing.html")
	return r
}

func main() {
	r := gin.Default()

	dbConnection, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	// Auto create these tables
	dbConnection.AutoMigrate(&Firing{}, &TemperatureReading{}, &Photo{})

	// Use middleware
	r.Use(AddDbHandle(dbConnection))

	// Use multitemplate rendering
	r.HTMLRender = createMyRender()

	// Setup static assets
	r.Static("/js", "./resources/js/")
	r.Static("/css", "./resources/css/")

	// Index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list", gin.H{"title": "Firings"})
	})
	r.GET("/firing", func(c *gin.Context) {
		c.HTML(http.StatusOK, "firing", gin.H{"title": "Firing"})
	})
	r.GET("/new-firing", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new-firing", gin.H{"title": "New Firing"})
	})

	// Create a new temperature reading
	r.POST("/temperature", func(c *gin.Context) {
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

		firingID := GetOrCreateFiring(db)
		newReading := TemperatureReading{FiringID: firingID, Inner: innerTemperature, Outer: outerTemperature}
		db.Create(&newReading)

		c.Status(200)
	})

	// Get details about a specific firing
	r.GET("/firing/:firingId", func(c *gin.Context) {
		firingID := c.Param("firingId")

		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var firing Firing
		db.Where("id = ?", firingID).Find(&firing)

		c.JSON(200, firing)

	})

	// List all firings
	r.GET("/firings", func(c *gin.Context) {

		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var firings []Firing
		db.Find(&firings)

		c.JSON(200, firings)

	})

	// List all readings for a firing
	r.GET("/firing/:firingId/readings", func(c *gin.Context) {
		firingID := c.Param("firingId")

		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var temperatureReadings []TemperatureReading
		db.Where("firing_id = ?", firingID).Find(&temperatureReadings)

		c.JSON(200, temperatureReadings)

	})

	r.Run(":8177")
	dbConnection.Close()

}
