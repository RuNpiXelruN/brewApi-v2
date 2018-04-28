package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var result utils.Result

// UpdateBeer func
func UpdateBeer(id, name, desc, stat, alc, ft, brewIDs, imageURL string) *utils.Result {
	feat, _ := strconv.ParseBool(ft)
	alcInt, _ := strconv.ParseFloat(alc, 64)
	var bIDs []string
	var brewers []Brewer

	beer := Beer{}

	err := db.Model(&Beer{}).Preload("Brewers.Rank").Where("id = ?", id).Find(&beer).Error
	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beer from DB",
		}
		return &result
	}

	err = db.Model(&beer).Updates(&Beer{
		Name:           name,
		Description:    desc,
		Status:         stat,
		ImageURL:       imageURL,
		AlcoholContent: alcInt,
	}).Error

	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error updating beer in DB",
		}
		return &result
	}

	if len(brewIDs) > 0 {
		bIDs = strings.Split(brewIDs, ",")
		if err := db.Model(&Brewer{}).Where("id in (?)", bIDs).Find(&brewers).Error; err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError) + " - Error fetching brewers from DB",
			}
			return &result
		}

		if err := db.Model(&beer).Association("Brewers").Replace(&brewers).Error; err != nil {
			result.Error = &utils.Error{
				Status:     http.StatusInternalServerError,
				StatusText: http.StatusText(http.StatusInternalServerError) + " - Error replacing beers brewers in DB",
			}
			return &result
		}
	}

	if err := db.Model(&beer).Update("featured", feat).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error updating featured status in DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &beer,
	}

	return &result
}

// GetBeersWithStatus func
func GetBeersWithStatus(status string) *utils.Result {
	beers := []Beer{}

	if err := db.Model(&Beer{}).Preload("Brewers.Rank").Where("status LIKE ?", status+"%").Find(&beers).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beers from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &beers,
	}
	return &result
}

// GetFeaturedBeers func
func GetFeaturedBeers(feat string) *utils.Result {
	beers := []Beer{}

	if err := db.Model(&Beer{}).Preload("Brewers.Rank").Where("featured = ?", feat).Find(&beers).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " -Error fetching beers from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &beers,
	}
	return &result
}

// DeleteBeer func
func DeleteBeer(id string) *utils.Result {
	beer := Beer{}

	if err := db.Model(&beer).Where("id = ?", id).Find(&beer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error finding beer in DB",
		}
		return &result
	}

	if err := db.Delete(&beer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error deleting beer from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   http.StatusText(http.StatusOK) + " - Beer successfully deleted",
	}
	return &result
}

// CreateBeer func
func CreateBeer(name, desc, alc, feat, brewIDs, imageURL string) *utils.Result {
	al, _ := strconv.ParseFloat(alc, 64)
	ft, _ := strconv.ParseBool(feat)

	bIDs := strings.Split(brewIDs, ",")

	var brUintIDs []uint
	for _, b := range bIDs {
		intID, _ := strconv.Atoi(b)
		brUintIDs = append(brUintIDs, uint(intID))
	}

	brewers := []Brewer{}
	if err := db.Model(&Brewer{}).Where(brUintIDs).Find(&brewers).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching brewers from DB",
		}
		return &result
	}

	beer := Beer{
		Name:           name,
		Description:    desc,
		AlcoholContent: al,
		Featured:       ft,
		ImageURL:       imageURL,
		Brewers:        brewers,
		CreatedAt:      time.Now(),
	}

	if err := db.Save(&beer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error saving beer to DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusCreated,
		Data:   &beer,
	}
	return &result
}

// GetBeer func
func GetBeer(id string) *utils.Result {
	beer := Beer{}
	if err := db.Model(&Beer{}).Preload("Brewers.Rank").Where("id = ?", id).Find(&beer).Error; err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beer from DB",
		}
		return &result
	}
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &beer,
	}
	return &result
}

// GetBeers func
func GetBeers(limit, order, offset string) *utils.Result {
	beers := []Beer{}

	err := db.Model(&Beer{}).Limit(limit).Order("created_at " + order).Offset(offset).Preload("Brewers.Rank").Find(&beers).Error
	if err != nil {
		result.Error = &utils.Error{
			Status:     http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching beers from DB",
		}
		return &result
	}

	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &beers,
	}
	return &result
}
