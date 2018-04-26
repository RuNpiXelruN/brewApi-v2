package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
)

// BasicRank struct
type BasicRank struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// GetBrewersOfRank func
func GetBrewersOfRank(level string) *utils.Result {
	rank := Rank{}
	if err := db.Model(&Rank{}).Preload("Brewers").Where("level = ?", level).Find(&rank).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error fetching ranked brewers from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &rank,
	}
	return &result
}

// GetRanks func
func GetRanks(limit, order, offset string) *utils.Result {
	ranks := []BasicRank{}
	if err := db.Model(&Rank{}).Limit(limit).Order("created_at " + order).Offset(offset).Select([]string{"id", "name", "level"}).Scan(&ranks).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error fetching ranks from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &ranks,
	}
	return &result
}
