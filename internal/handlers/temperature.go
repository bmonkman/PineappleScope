package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bmonkman/PineappleScope/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const newFiringThreshold = 3 // hour threshold for new measurements to be considered part of a firing
// the temp must go past the low notification amount times this number before dropping back down
const lowTempNotificationThresholdModifier = 1.1

const iftttNotificationName = "kiln_temperature_reached"

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

	handleNotifications(db, firingID, innerTemperature)

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

func handleNotifications(db *gorm.DB, firingID uint, temperature float64) {
	firing := models.Firing{}
	db.First(&firing, firingID)

	if firing.HighNotificationTemp > 0 && temperature > firing.HighNotificationTemp && firing.HighNotificationSent == false {
		sendNotification("🔥 high", firingID, temperature)
		firing.HighNotificationSent = true
		db.Save(&firing)
	}

	if firing.LowNotificationTemp > 0 && temperature < firing.LowNotificationTemp && firing.LowNotificationSent == false {
		temperatureRecord := models.TemperatureReading{}

		// Make sure the temperature has previously gone above the low notification threshold
		whereString := fmt.Sprintf("firing_id = %d AND inner > %f", firingID, firing.LowNotificationTemp*lowTempNotificationThresholdModifier)
		found := !db.Where(whereString).
			Order("created_date desc").
			First(&temperatureRecord).
			RecordNotFound()

		if found {
			sendNotification("❄️ low", firingID, temperature)
			firing.LowNotificationSent = true
			db.Save(&firing)
		}
	}
}

func sendNotification(notificationType string, firingID uint, temperature float64) {
	iftttURL := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", iftttNotificationName, os.Getenv("IFTTT_KEY"))
	values := url.Values{
		"value1": {notificationType},
		"value2": {strconv.FormatFloat(temperature, 'f', 2, 64)},
		"value3": {strconv.FormatUint(uint64(firingID), 10)},
	}
	resp, err := http.PostForm(iftttURL, values)
	if err != nil {
		fmt.Println("Failed to send notification due to error: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to send notification, status code: ", resp.StatusCode)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

	}
}
