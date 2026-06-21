package handlers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/bmonkman/pineapplescope/internal/models"
)

// statsRetentionDays is how long full-resolution device stats are kept. Stats
// are reported every 60s and rarely reviewed, so old rows are deleted to keep
// the table from growing unbounded.
const statsRetentionDays = 90

// PruneOldStats deletes stats records older than the retention window. The
// cutoff is expressed as a SQLite datetime so it matches how created_date is
// stored (datetime('now','localtime')) rather than a Go time.Time.
func PruneOldStats(db *gorm.DB) {
	whereString := fmt.Sprintf("created_date < datetime('now', 'localtime', '-%d days')", statsRetentionDays)
	result := db.Where(whereString).Delete(&models.Stats{})
	if result.Error != nil {
		fmt.Println("stats cleanup failed:", result.Error)
		return
	}
	fmt.Printf("stats cleanup: removed %d rows older than %d days\n", result.RowsAffected, statsRetentionDays)
}

// StartStatsCleanup prunes old stats immediately, then once every 24 hours.
// Running on startup means a restart doesn't have to wait a full day.
func StartStatsCleanup(db *gorm.DB) {
	PruneOldStats(db)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			PruneOldStats(db)
		}
	}()
}
