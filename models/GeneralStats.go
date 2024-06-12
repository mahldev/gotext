package models

type GeneralStats struct {
	TotalTests        int   `json:"totalTests"`
	TotalWordsWritten int   `json:"totalWordsWritten"`
	TotalErrors       int   `json:"totalErrors"`
	Level             int64 `json:"level"`
}
