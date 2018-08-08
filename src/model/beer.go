package model

import (
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// GetBasicBeers func
func GetBasicBeers() *utils.Result {
	beers := []BasicBeer{}
	if err := db.Model(&Beer{}).Select([]string{"id", "name"}).Scan(&beers).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Cannot fetch basic beers")
	}
	return dbSuccess(beers)
}

// GetBeer func
func GetBeer(id string) *utils.Result {
	beer := Beer{}
	if err := db.Model(&Beer{}).Preload("Brewers").Where("id = ?", id).Find(&beer).Error; err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beer from DB")
	}

	return dbSuccess(&beer)
}

// GetBeers func
func GetBeers(limit, order, offset string) *utils.Result {
	beers := []Beer{}
	err := db.Model(&Beer{}).Preload("Brewers.Rank").Limit(limit).Offset(offset).Order("updated_at desc").Find(&beers).Error
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error fetching beers from Db")
	}
	return dbSuccess(&beers)
}

// CreateBeer func
func CreateBeer(name, desc, status, alc, feat, brewerIDs string, image ReqImage) *utils.Result {

	var imgURL string

	if image.Error == nil {
		imageURL, err := s3ImageUpload(image)
		if err != nil {
			fmt.Println("Error uploading to S3:", err.Error())
		}
		imgURL = imageURL.(string)
	}

	tx := db.Begin()
	brewers, err := setBrewers(brewerIDs, tx)
	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
	}

	alcFl, _ := strconv.ParseFloat(alc, 64)
	featured, _ := strconv.ParseBool(feat)

	beer := Beer{
		Name:           name,
		Description:    desc,
		Status:         status,
		AlcoholContent: alcFl,
		Featured:       featured,
		Brewers:        *brewers,
		ImageURL:       imgURL,
	}

	if err := tx.Save(&beer).Error; err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusInternalServerError, "Error saving beer to DB")
	}

	tx.Commit()
	return dbSuccess(&beer)
}

// UpdateBeer func
func UpdateBeer(id, name, desc, status, alc, feat, brewerIDs string, image ReqImage) *utils.Result {
	beer := Beer{}
	var imgURL string

	tx := db.Begin()
	if err := tx.Model(&beer).Preload("Brewers.Rank").Where("id = ?", id).Find(&beer).Error; err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusNotFound, "Error fetching beer from DB")
	}

	imageURL, err := s3ImageUpload(image)
	if err == nil {
		imgURL = imageURL.(string)
	}

	brewers, err := setBrewers(brewerIDs, tx)
	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
	}

	alcFl, _ := strconv.ParseFloat(alc, 64)

	if err := tx.Model(&beer).Updates(&Beer{
		Name:           name,
		Description:    desc,
		Status:         status,
		AlcoholContent: alcFl,
		ImageURL:       imgURL,
	}).Error; err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusInternalServerError, "Error updating beer")
	}

	if len(feat) > 0 {
		featured, _ := strconv.ParseBool(feat)
		err := tx.Model(&beer).Update("featured", featured).Error
		if err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating featured status of beer")
		}
	}

	if len(*brewers) > 0 {
		err := tx.Model(&beer).Association("Brewers").Replace(*brewers).Error
		if err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating beers brewers association")
		}
	}

	tx.Commit()
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
func GetBeersWithStatus(status, limit, ord string) *utils.Result {
	order := "desc"
	if len(ord) > 0 {
		order = ord
	}
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

// ************************************************************ UTILITY FUNCTIONS ************************************************************ //

// setBrewers func
func setBrewers(bIDs string, tx *gorm.DB) (*[]Brewer, error) {
	brewers := []Brewer{}
	if len(bIDs) > 0 {
		brewerIDs := strings.Split(bIDs, ",")

		if err := tx.Model(&Brewer{}).Preload("Rank").Where("id in (?)", brewerIDs).Find(&brewers).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return &brewers, nil
}
