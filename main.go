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
const version = "0.0.1"

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
	CreatedDate time.Time `gorm:"default:(datetime('now','localtime'))"`
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

// Calculate the current cone number based on the temperature
func CalculateCone(temperature float64) uint {
	/*var cones char[float] = {
		'1': 2077.0
		'2': 2088
		'3': 2106
		'4': 2120
		'5': 2163
		'6': 2228
		'7': 2259
	}*/
	return 1
}

// GetOrCreateFiring If there was a recent temperature reading, return that reading's session, otherwise create a new session
func GetOrCreateFiring(db *gorm.DB) uint {
	var temperature = TemperatureReading{}

	whereString := fmt.Sprintf("created_date >= datetime('now', 'localtime', '-%d hours')", newFiringThreshold)
	db.Where(whereString).
		Order("created_date desc").
		First(&temperature)

	var firing Firing
	if temperature.ID == 0 || temperature.FiringID == 0 {
		firing = Firing{StartDate: time.Now(), EndDate: time.Now(), Name: "New Firing"}
		db.Save(&firing)
	} else {
		db.First(&firing, temperature.FiringID)
		firing.EndDate = time.Now()

		currentCone := CalculateCone(temperature.Inner)
		if firing.Cone < currentCone {
			firing.Cone = currentCone
		}

		db.Save(&firing)
	}

	return firing.ID
}

func setupTemplates() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("list", "resources/html/base.html", "resources/html/list.html")
	r.AddFromFiles("firing", "resources/html/base.html", "resources/html/firing.html")
	r.AddFromFiles("new-firing", "resources/html/base.html", "resources/html/new-firing.html")
	return r
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	dbConnection, err := gorm.Open("sqlite3", "/var/db/pineapplescope.db")
	if err != nil {
		panic(err)
	}

	// Auto create these tables
	dbConnection.AutoMigrate(&Firing{}, &TemperatureReading{}, &Photo{})

	// Use middleware
	r.Use(AddDbHandle(dbConnection))

	// Use multitemplate rendering
	r.HTMLRender = setupTemplates()

	// Setup static assets
	r.Static("/js", "./resources/js/")
	r.Static("/css", "./resources/css/")

	// Index
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/firings")
	})
	r.GET("/new-firing", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new-firing", gin.H{"title": "New Firing"})
	})
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": version})
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
		db.Save(&newReading)

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
		db.First(&firing, firingID)

		var temperatureReadings []TemperatureReading
		db.Where("firing_id = ?", firingID).Find(&temperatureReadings)

		c.HTML(http.StatusOK, "firing", gin.H{"title": "Firing: " + firing.Name, "firing": firing, "temperatureReadings": temperatureReadings})

	})

	// Show form to edit details about a specific firing
	r.GET("/firing/:firingId/edit", func(c *gin.Context) {
		firingID := c.Param("firingId")

		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var firing Firing
		db.First(&firing, firingID)

		c.HTML(http.StatusOK, "new-firing", gin.H{"title": "Edit Firing: " + firing.Name, "firing": firing})

	})

	// Edit details about a specific firing
	r.POST("/firing/:firingId", func(c *gin.Context) {
		firingID := c.Param("firingId")
		name := c.PostForm("name")
		notes := c.PostForm("notes")
		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var firing Firing
		db.Where("id = ?", firingID).Find(&firing)
		db.First(&firing, firingID)
		firing.Name = name
		firing.Notes = notes
		fmt.Println(firing)
		db.Save(&firing)

		c.Redirect(301, "/firing/"+firingID)
	})

	// List all firings
	r.GET("/firings", func(c *gin.Context) {

		db, ok := c.MustGet("databaseConn").(*gorm.DB)
		if !ok {
			return
		}

		var firings []Firing
		db.Find(&firings)

		c.HTML(http.StatusOK, "list", gin.H{"title": "Firings", "firings": firings})

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

	r.Run(":1111")
	dbConnection.Close()

}
