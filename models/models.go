package models

import (
	"time"
)

type Firing struct {
	ID                   uint
	StartDate            time.Time
	EndDate              time.Time
	StartDateAmbientTemp float64
	ConeNumber           string
	Name                 string
	Notes                string

	HighNotificationTemp float64
	LowNotificationTemp  float64 `gorm:"default:100.0"`
	HighNotificationSent bool    `gorm:"default:0"`
	LowNotificationSent  bool    `gorm:"default:0"`

	TemperatureReadings []TemperatureReading
	Photos              []Photo
}

func (f *Firing) Duration() int64 {
	return (f.EndDate.Unix() - f.StartDate.Unix()) / 60 / 60
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

type Stats struct {
	ID                 uint
	CreatedDate        time.Time `gorm:"default:(datetime('now','localtime'))"`
	Uptime             uint64
	FreeMemory         uint64
	Temperature        float64
	AmbientTemperature float64
	CPUTemperature     float64
	Humidity           float64
}
