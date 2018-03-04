package main

import (
	"net/http"
	"os"

	"github.com/bmonkman/PineappleScope/handlers"
	"github.com/bmonkman/PineappleScope/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-contrib/multitemplate"

	"github.com/gin-gonic/gin"
)

const version = "0.0.2"

// AddDbHandle middleware will add the db connection to the context
func AddDbHandle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
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

	dbFile := os.Getenv("DBFILE")
	if dbFile == "" {
		dbFile = "pineapplescope.db"
	}
	dbConnection, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	// Auto create these tables
	dbConnection.AutoMigrate(&models.Firing{}, &models.TemperatureReading{}, &models.Photo{})

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

	handlers.NewFiringHandlers(r).Register()
	handlers.NewTemperatureHandlers(r).Register()

	r.Run(":1111")
	dbConnection.Close()
}
