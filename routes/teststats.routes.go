package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
)

var wpmThresholds = map[int64]float64{
	1:  50.0,
	2:  60.0,
	3:  70.0,
	4:  80.0,
	5:  90.0,
	6:  100.0,
	7:  110.0,
	8:  120.0,
	9:  140.0,
	10: 150.0,
}

var accThresholds = map[int64]float64{
	1:  80.0,
	2:  85.0,
	3:  90.0,
	4:  95.0,
	5:  96.0,
	6:  97.0,
	7:  98.0,
	8:  99.0,
	9:  100.0,
	10: 100.0,
}

func checkIfUserLvlUp(u m.User) bool {
	var totalWordsWritten int64
	var avgWPM float64
	var avgAccuracy float64
	var err error

	err = db.DB.Model(&m.TestStats{}).Where("user_id = ?", u.ID).Select("SUM(words)").Scan(&totalWordsWritten).Error
	if err != nil {
		return false
	}
	err = db.DB.Model(&m.TestStats{}).Where("user_id = ?", u.ID).Select("AVG(accuracy)").Scan(&avgAccuracy).Error
	if err != nil {
		return false
	}
	err = db.DB.Model(&m.TestStats{}).Where("user_id = ?", u.ID).Select("AVG(wpm)").Scan(&avgWPM).Error
	if err != nil {
		return false
	}

	newLevel := calculateUserLevel(totalWordsWritten, avgWPM, avgAccuracy)

	return newLevel > u.Level
}

func calculateUserLevel(totalWordsWritten int64, avgWPM float64, avgAccuracy float64) int64 {
	switch {
	case totalWordsWritten > 50000 && avgWPM >= 150 && avgAccuracy >= 98:
		return 10
	case totalWordsWritten > 40000 && avgWPM >= 140 && avgAccuracy >= 97:
		return 9
	case totalWordsWritten > 30000 && avgWPM >= 120 && avgAccuracy >= 96:
		return 8
	case totalWordsWritten > 20000 && avgWPM >= 110 && avgAccuracy >= 95:
		return 7
	case totalWordsWritten > 15000 && avgWPM >= 100 && avgAccuracy >= 94:
		return 6
	case totalWordsWritten > 10000 && avgWPM >= 90 && avgAccuracy >= 93:
		return 5
	case totalWordsWritten > 2000 && avgWPM >= 80 && avgAccuracy >= 90:
		return 4
	case totalWordsWritten > 1000 && avgWPM >= 70 && avgAccuracy >= 85:
		return 3
	case totalWordsWritten > 10 && avgWPM >= 60 && avgAccuracy >= 80:
		return 2
	default:
		return 1
	}
}

func SaveTestStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := &m.TestStats{}
	if err := json.NewDecoder(r.Body).Decode(stats); err != nil {
		WriteError(w, "Invalid request payload")
		return
	}

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

	u := &m.User{}
	db.DB.First(u, id)

	stats.UserID = u.ID
	stats.CreatedAt = time.Now()
	stats.User = u
	levelUp := checkIfUserLvlUp(*u)

	if err := db.DB.Create(&stats).Error; err != nil {
		WriteError(w, "Failed to save statistics")
		return
	}

	response := map[string]bool{"levelUp": levelUp}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetTestStatsHandler(w http.ResponseWriter, r *http.Request) {
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

	type tests struct {
		Tests []m.TestStats `json:"tests"`
	}

	stats := make([]m.TestStats, 0)
	db.DB.Where("user_id = ?", id).Find(&stats)

	response := &tests{Tests: stats}
	json.NewEncoder(w).Encode(response)
}

func ValidatelevelUpTest(w http.ResponseWriter, r *http.Request) {
	stats := &m.TestStats{}
	if err := json.NewDecoder(r.Body).Decode(stats); err != nil {
		WriteError(w, "Invalid request payload")
		return
	}

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

	u := &m.User{}
	db.DB.First(u, id)

	wpmthreshold, exists := wpmThresholds[u.Level]
	if !exists {
		WriteError(w, "Invalid user level")
		return
	}

	accthreshold, exists := accThresholds[u.Level]
	if !exists {
		WriteError(w, "Invalid user level")
		return
	}

	if stats.WPM >= wpmthreshold && stats.Accuracy >= accthreshold {
		u.Level += 1
		if result := db.DB.Save(u); result.Error != nil {
			WriteError(w, "Failed to update user level")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Level up successful!",
		})
	} else {
		w.WriteHeader(http.StatusConflict)
		WriteError(w, "User did not pass the test without errors")
	}
}
