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
	ID          uint
	CreatedDate time.Time `gorm:"default:(datetime('now','localtime'))"`
	Uptime      uint64
	FreeMemory  uint64
	Event       string
}
