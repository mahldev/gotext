package routes

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
	"gorm.io/gorm"
)

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		WriteError(w, "Unauthorized")
		return
	}

	userID, ok := (*claims)["userID"].(string)
	if !ok {
		WriteError(w, "Invalid token claims")
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		WriteError(w, "Invalid uuid")
		return
	}

	var user m.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			WriteError(w, "User not found")
		} else {
			WriteError(w, "Database error")
		}
		return
	}

	var totalTests int64
	var totalWordsWritten int64
	var totalErrors int64

	if err := db.DB.Model(&m.TestStats{}).Where("user_id = ?", id).Count(&totalTests).Error; err != nil {
		WriteError(w, "Failed to retrieve total tests")
		return
	}

	if err := db.DB.Model(&m.TestStats{}).Where("user_id = ?", id).Select("SUM(words)").Scan(&totalWordsWritten).Error; err != nil {
		WriteError(w, "Failed to retrieve total words written")
		return
	}

	if err := db.DB.Model(&m.TestStats{}).Where("user_id = ?", id).Select("SUM(errors)").Scan(&totalErrors).Error; err != nil {
		WriteError(w, "Failed to retrieve total errors")
		return
	}

	level := user.Level

	stats := m.GeneralStats{
		TotalTests:        int(totalTests),
		TotalWordsWritten: int(totalWordsWritten),
		TotalErrors:       int(totalErrors),
		Level:             level,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		WriteError(w, "Failed to encode statistics")
		return
	}

	w.WriteHeader(http.StatusOK)
}
