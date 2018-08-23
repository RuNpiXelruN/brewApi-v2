package db

import (
	"go_apps/go_api_apps/brewApi-v2/utils"
	"net/http"
)

// BasicRank struct
type BasicRank struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// GetBrewersOfRank func
func GetBrewersOfRank(rnk, limit, order, offset string) *utils.Result {
	rank := Rank{}

	err := db.Model(&Rank{}).
		Limit(limit).Order("created_at "+order).Offset(offset).
		Where("level = ?", rnk).Preload("Brewers.Beers").Find(&rank).Error
	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching ranked brewers from db")
	}

	return dbSuccess(&rank)
}

// GetRanks func
func GetRanks(limit, order, offset string) *utils.Result {
	ranks := []BasicRank{}
	err := db.Model(&Rank{}).Limit(limit).Order("created_at " + order).Offset(offset).Select([]string{"id", "name", "level"}).Scan(&ranks).Error
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error fetching ranks from DB")
	}

	return dbSuccess(&ranks)
}
