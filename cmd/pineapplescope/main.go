package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bmonkman/PineappleScope/handlers"
	"github.com/bmonkman/PineappleScope/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-contrib/multitemplate"

	"github.com/gin-gonic/gin"
)

const version = "0.1.0"

// AddDbHandle middleware will add the db connection to the context
func AddDbHandle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
}

// AddSharedVars middleware will add shared vars to all templates
func AddSharedVars(vars map[string]string, funcs map[string]func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range vars {
			c.Set(k, v)
		}

		for k, f := range funcs {
			c.Set(k, f(c))
		}

		c.Next()
	}
}

func setupTemplates() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("list", "resources/html/base.html", "resources/html/list.html")
	r.AddFromFiles("firing", "resources/html/base.html", "resources/html/firing.html")
	r.AddFromFiles("new-firing", "resources/html/base.html", "resources/html/new-firing.html")
	r.AddFromFiles("stats", "resources/html/base.html", "resources/html/stats.html")
	return r
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	// r.SetTrustedProxies([]string{})

	dbFile := os.Getenv("DBFILE")
	if len(dbFile) == 0 {
		dbFile = "pineapplescope.db"
	}

	dbConnection, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	// Auto create these tables
	dbConnection.AutoMigrate(&models.Firing{}, &models.TemperatureReading{}, &models.Photo{}, &models.Stats{})

	if os.Getenv("DB_DEBUG") == "true" {
		dbConnection.LogMode(true)
		dbConnection.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}

	// Use middleware
	r.Use(AddDbHandle(dbConnection))
	sharedVars := map[string]string{"version": version}
	sharedFuncs := map[string]func(*gin.Context) string{
		"deviceCheckedIn": getDeviceCheckedInState,
		"currentTemp":     getCurrentTemp}

	r.Use(AddSharedVars(sharedVars, sharedFuncs))

	// Use multitemplate rendering
	r.HTMLRender = setupTemplates()

	// Setup static assets
	r.Static("/js", "./resources/js/")
	r.Static("/css", "./resources/css/")
	r.Static("/images", "./resources/images/")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Index
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/firings")
	})
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	handlers.NewFiringHandlers(r).Register()
	handlers.NewTemperatureHandlers(r).Register()
	handlers.NewStatsHandlers(r).Register()

	r.Run(":1111")
	dbConnection.Close()
}

func getDeviceCheckedInState(c *gin.Context) string {
	db, _ := c.MustGet("databaseConn").(*gorm.DB)
	statsRecord := models.Stats{}
	notFound := db.Where("created_date >= datetime('now', 'localtime', '-2 minutes')").
		First(&statsRecord).
		RecordNotFound()
	if notFound {
		return "0"
	}
	return "1"
}

func getCurrentTemp(c *gin.Context) string {
	db, _ := c.MustGet("databaseConn").(*gorm.DB)
	statsRecord := models.Stats{}
	notFound := db.Order("created_date desc").
		First(&statsRecord).
		RecordNotFound()
	if notFound {
		return "0"
	}
	return strconv.FormatFloat(statsRecord.Temperature, 'f', 2, 64)
}
