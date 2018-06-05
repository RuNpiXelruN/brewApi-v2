package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"
)

// GetBrewers func
func GetBrewers(limit, order, offset string) *utils.Result {
	brewers := []Brewer{}
	if err := db.Model(&Brewer{}).
		Limit(limit).
		Order("created_at " + order).
		Offset(offset).
		Preload("Beers").Preload("Rank").Find(&brewers).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
	}

	return dbSuccess(brewers)
}

// GetBrewer func
func GetBrewer(id string) *utils.Result {
	brewer := Brewer{}
	if err := db.Model(&Brewer{}).Preload("Beers").Preload("Rank").
		Where("id = ?", id).Find(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching Brewer from DB")
	}

	return dbSuccess(brewer)
}

// CreateBrewer func
func CreateBrewer(first, last, feat, username, rnk, beerIDs string) *utils.Result {
	var beers []Beer
	if len(beerIDs) > 0 {
		bIDs := strings.Split(beerIDs, ",")
		if err := db.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error; err != nil {
			return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
		}
	}

	ft, _ := strconv.ParseBool(feat)
	brewer := Brewer{
		FirstName: first,
		LastName:  last,
		Featured:  ft,
		Username:  &username,
		Beers:     beers,
	}

	if len(rnk) > 0 {
		rank := Rank{}
		if err := db.Model(&Rank{}).Where("level = ?", rnk).Find(&rank).Error; err != nil {
			return dbWithError(err, http.StatusNotFound, "Error fetching rank from DB")
		}

		brewer.Rank = &rank
	}

	if err := db.Save(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error creating brewer")
	}

	return dbSuccess(&brewer)
}

// UpdateBrewer func
func UpdateBrewer(id, first, last, ft, username, rnk, beerIDs string) *utils.Result {
	brewer := Brewer{}

	if err := db.Model(&Brewer{}).
		Preload("Beers").
		Preload("Rank").
		Where("id = ?", id).
		Find(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewer from DB")
	}

	if err := db.Model(&brewer).Updates(&Brewer{
		FirstName: first,
		LastName:  last,
	}).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error updating brewer")
	}

	if len(ft) > 0 {
		feat, _ := strconv.ParseBool(ft)
		if err := db.Model(&brewer).Update("featured", feat).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer featured status")
		}
	}

	if len(username) > 0 {
		if err := db.Model(&brewer).Update("username", &username).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer username")
		}
	}

	if len(beerIDs) > 0 {
		beers := []Beer{}
		bIDs := strings.Split(beerIDs, ",")
		if err := db.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error; err != nil {
			return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
		}

		if err := db.Model(&brewer).Association("Beers").Replace(&beers).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer's beers")
		}
	}

	if len(rnk) > 0 {
		rank := Rank{}
		if err := db.Model(&Rank{}).Where("level = ?", rnk).Find(&rank).Error; err != nil {
			return dbWithError(err, http.StatusNotFound, "Error fetching rank from DB")
		}

		if err := db.Model(&brewer).Association("Rank").Replace(&rank).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer rank")
		}
	}

	return dbSuccess(&brewer)
}

// DeleteBrewer func
func DeleteBrewer(id string) *utils.Result {
	brewer := Brewer{}
	if err := db.Model(&Brewer{}).Where("id = ?", id).Find(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewer from DB")
	}

	if err := db.Delete(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error deleting brewer from DB")
	}

	return dbSuccess("Successfully deleted brewer from DB")
}

// GetRankedBrewers func
func GetRankedBrewers(level string) *utils.Result {
	brewers := []Brewer{}

	if err := db.Joins("JOIN ranks on ranks.id = brewers.rank_id").Where("ranks.level = ?", level).Preload("Rank").Find(&brewers).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
	}

	return dbSuccess(brewers)
}

// GetFeaturedBrewers func
func GetFeaturedBrewers(feat string) *utils.Result {
	brewers := []Brewer{}

	err := db.Model(&Brewer{}).Preload("Beers").Preload("Rank").
		Where("featured = ?", feat).Find(&brewers).Error

	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching featured brewers from DB")
	}

	return dbSuccess(brewers)
}
