package model

import (
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"
)

var result utils.Result

// GetBeer func
func GetBeer(id string) *utils.Result {
	beer := Beer{}
	if err := db.Model(&Beer{}).Preload("Brewers.Rank").Where("id = ?", id).Find(&beer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beer from DB")
	}

	return dbSuccess(&beer)
}

// GetBeers func
func GetBeers(limit, order, offset string) *utils.Result {
	beers := []Beer{}
	err := db.Model(&Beer{}).Preload("Brewers.Rank").Limit(limit).Offset(offset).Order("created_at " + order).Find(&beers).Error
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error fetching beers from Db")
	}
	return dbSuccess(&beers)
}

// CreateBeer func
func CreateBeer(name, desc, status, alc, feat, brewerIDs string, image ReqImage) *utils.Result {
	var imgURL string
	imageURL, err := s3ImageUpload(image)
	if err != nil {
		fmt.Println("Error uploading to S3:", err.Error())
	}
	imgURL = imageURL.(string)

	bIDs := strings.Split(brewerIDs, ",")
	brewers := []Brewer{}
	if err := db.Model(&Brewer{}).Where("id in (?)", bIDs).Find(&brewers).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
	}

	alcFl, _ := strconv.ParseFloat(alc, 64)
	ft, _ := strconv.ParseBool(feat)

	beer := Beer{
		Name:           name,
		Description:    desc,
		Status:         status,
		AlcoholContent: alcFl,
		Featured:       ft,
		Brewers:        brewers,
		ImageURL:       imgURL,
	}

	if err := db.Save(&beer).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error saving beer to DB")
	}
	return dbSuccess(&beer)
}

// UpdateBeer func
func UpdateBeer(id, name, desc, status, alc, feat, brewerIDs string, image ReqImage) *utils.Result {
	beer := Beer{}
	if err := db.Model(&beer).Preload("Brewers.Rank").Where("id = ?", id).Find(&beer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beer from DB")
	}

	var imgURL string
	imageURL, err := s3ImageUpload(image)
	if err == nil {
		imgURL = imageURL.(string)
	}

	var brewers []Brewer
	if len(brewerIDs) > 0 {
		bIDs := strings.Split(brewerIDs, ",")
		if err := db.Model(&Brewer{}).Where("id in (?)", bIDs).Find(&brewers).Error; err != nil {
			return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
		}
	}

	alcFl, _ := strconv.ParseFloat(alc, 64)
	ft, _ := strconv.ParseBool(feat)

	if err := db.Model(&beer).Updates(&Beer{
		Name:           name,
		Description:    desc,
		Status:         status,
		AlcoholContent: alcFl,
		ImageURL:       imgURL,
	}).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error updating beer")
	}

	if err := db.Model(&beer).Update("featured", ft).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error updating featured status of beer")
	}

	if len(brewers) > 0 {
		if err := db.Model(&beer).Association("Brewers").Replace(&brewers).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error updating beers brewers association")
		}
	}

	return dbSuccess(&beer)
}

// DeleteBeer func
func DeleteBeer(id string) *utils.Result {
	beer := Beer{}
	if err := db.Model(&Beer{}).Where("id = ?", id).Find(&beer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error finding beer in DB")
	}
	if err := db.Delete(beer).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error deleting beer from DB")
	}

	return dbSuccess("Beer succesfully deleted")
}

// GetBeersWithStatus func
func GetBeersWithStatus(status, limit, order string) *utils.Result {
	beers := []Beer{}
	if err := db.Model(&Beer{}).Limit(limit).Order("created_at "+order).
		Preload("Brewers.Rank").
		Where("status LIKE ?", status+"%").
		Find(&beers).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
	}

	return dbSuccess(&beers)
}

// GetFeaturedBeers func
func GetFeaturedBeers(feat, limit, order string) *utils.Result {
	beers := []Beer{}
	if err := db.Model(&Beer{}).Limit(limit).Order("created_at "+order).
		Preload("Brewers.Rank").
		Where("featured = ?", feat).
		Find(&beers).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching featured beers from DB")
	}

	return dbSuccess(&beers)
}
