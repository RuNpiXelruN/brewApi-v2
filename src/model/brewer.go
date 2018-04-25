package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"
)

// UpdateBrewer func
func UpdateBrewer(id, f, l, ft, rnk, beerIDs string) *utils.Result {
	feat, _ := strconv.ParseBool(ft)
	var bIDs []string
	var rank Rank
	var beers []Beer

	brewer := Brewer{}
	if err := db.Model(&Brewer{}).Preload("Rank").Preload("Beers").Where("id = ?", id).Find(&brewer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching brewer from DB",
		}
		return &result
	}

	err := db.Model(&brewer).Updates(&Brewer{
		FirstName: f,
		LastName:  l,
	}).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error updating brewer in DB",
		}
		return &result
	}

	if len(rnk) > 0 {
		if err := db.Model(&Rank{}).Where("level = ?", rnk).Find(&rank).Error; err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusNotFound,
				StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching rank from DB",
			}
			return &result
		}

		if err := db.Model(&brewer).Association("Rank").Replace(&rank).Error; err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError) + " - Error replacing brewer rank in DB",
			}
			return &result
		}
	}

	if len(beerIDs) > 0 {
		bIDs = strings.Split(beerIDs, ",")
		err := db.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error
		if err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusNotFound,
				StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beers from DB",
			}
			return &result
		}

		if err := db.Model(&brewer).Association("Beers").Replace(&beers).Error; err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError) + " - Error replacing brewer beers in DB",
			}
			return &result
		}
	}

	if err := db.Model(&brewer).Update("featured", feat).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error updating featured status in DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewer,
	}

	return &result
}

// DeleteBrewer func
func DeleteBrewer(id string) *utils.Result {
	brewer := Brewer{}
	if err := db.Model(&Brewer{}).Where("id = ?", id).Find(&brewer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching brewer from DB",
		}
		return &result
	}
	if err := db.Delete(&brewer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error deleting brewer from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   http.StatusText(http.StatusOK) + " - Successfully deleted brewer from DB",
	}
	return &result
}

// CreateBrewer
func CreateBrewer(first, last, feat, rank, beerIDs string) *utils.Result {
	bIDs := strings.Split(beerIDs, ",")
	f, _ := strconv.ParseBool(feat)

	beers := []Beer{}
	if err := db.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beers from DB",
		}
		return &result
	}
	r := Rank{}
	if err := db.Model(&Rank{}).Where("level = ?", rank).Find(&r).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching rank from DB",
		}
		return &result
	}

	brewer := Brewer{
		FirstName: first,
		LastName:  last,
		Featured:  f,
		Rank:      &r,
		Beers:     beers,
	}

	if err := db.Save(&brewer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error saving brewer to DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewer,
	}
	return &result
}

// GetRankedBrewers func
func GetRankedBrewers(level string) *utils.Result {
	brewers := []Brewer{}

	err := db.Model(&Brewer{}).Preload("Rank").Joins("inner join ranks on ranks.brewer_id = brewers.id").Where("ranks.level = ?", level).Find(&brewers).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beers from DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewers,
	}
	return &result
}

// GetFeaturedBrewers func
func GetFeaturedBrewers(feat string) *utils.Result {
	brewers := []Brewer{}

	err := db.Model(&Brewer{}).Preload("Beers").Preload("Rank").
		Where("featured = ?", feat).Find(&brewers).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching featured brewers from DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewers,
	}
	return &result
}

// GetBrewer func
func GetBrewer(id string) *utils.Result {
	brewer := Brewer{}

	err := db.Model(&Brewer{}).Preload("Beers").Preload("Rank").
		Where("id = ?", id).Find(&brewer).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beer from DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewer,
	}
	return &result
}

// GetBrewers func
func GetBrewers(lim, ord, offs string) *utils.Result {
	brewers := []Brewer{}

	err := db.Model(&Brewer{}).
		Limit(lim).
		Order("created_at " + ord).
		Offset(offs).
		Preload("Beers").Preload("Rank").Find(&brewers).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching brewers from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &brewers,
	}
	return &result
}
