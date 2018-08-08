package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// GetBrewerNames func
func GetBrewerNames() *utils.Result {
	names := []BasicBrewer{}
	err := db.Model(&Brewer{}).Order("first_name asc").Select([]string{"id", "first_name", "last_name"}).Scan(&names).Error
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error fetching brewer names from DB")
	}

	return dbSuccess(&names)
}

// GetBrewers func
func GetBrewers(limit, order, offset string) *utils.Result {
	brewers := []Brewer{}

	tx := db.Begin()
	if err := tx.Model(&Brewer{}).
		Limit(limit).
		Order("created_at " + order).
		Offset(offset).
		Preload("Beers").Preload("Rank").Find(&brewers).Error; err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusNotFound, "Error fetching brewers from DB")
	}
	tx.Commit()
	return dbSuccess(brewers)
}

// GetBrewer func
func GetBrewer(id, includeBeers string) *utils.Result {
	brewer := Brewer{}

	tx := db.Begin()
	if err := tx.Model(&Brewer{}).Preload("Beers").Preload("Rank").
		Where("id = ?", id).Find(&brewer).Error; err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusNotFound, "Error fetching Brewer from DB")
	}

	beers := []BasicBeer{}
	if len(includeBeers) > 0 {
		include, _ := strconv.ParseBool(includeBeers)
		if include == true {
			err := tx.Model(&Beer{}).Order("name asc").Select([]string{"id", "name"}).Scan(&beers).Error
			if err != nil {
				tx.Rollback()
				return dbWithError(err, http.StatusInternalServerError, "Error fetching Basic Beers")
			}

			tx.Commit()

			data := make(map[string]interface{})
			data["brewer"] = &brewer
			data["beers"] = &beers

			return dbSuccess(data)
		}
	}

	tx.Commit()
	return dbSuccess(&brewer)
}

// CreateBrewer func
func CreateBrewer(first, last, feat, username, rnk, beerIDs string) *utils.Result {
	var beers []Beer

	tx := db.Begin()
	if len(beerIDs) > 0 {
		bIDs := strings.Split(beerIDs, ",")
		if err := tx.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error; err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
		}
	}

	rank, err := setRankOrNull(rnk, tx)
	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching rank from DB")
	}

	ft, _ := strconv.ParseBool(feat)
	brewer := Brewer{
		FirstName: first,
		LastName:  last,
		Featured:  ft,
		Username:  setStringValOrNil(username),
		Rank:      rank,
		Beers:     beers,
	}

	if err := tx.Save(&brewer).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error creating brewer")
	}

	tx.Commit()
	return dbSuccess(&brewer)
}

// UpdateBrewer func
func UpdateBrewer(id, first, last, ft, username, rnk, beerIDs string) *utils.Result {
	brewer := Brewer{}

	tx := db.Begin()
	err := tx.Model(&Brewer{}).
		Preload("Beers").
		Preload("Rank").
		Where("id = ?", id).
		Find(&brewer).Error
	if err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusNotFound, "Error fetching brewer from DB")
	}

	rank, err := setRankOrNull(rnk, db)
	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching rank from DB")
	}

	err = tx.Model(&brewer).Updates(&Brewer{
		FirstName: first,
		LastName:  last,
	}).Error
	if err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusInternalServerError, "Error updating brewer")
	}

	if len(ft) > 0 {
		feat, _ := strconv.ParseBool(ft)
		if err := tx.Model(&brewer).Update("featured", feat).Error; err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer featured status")
		}
	}

	if len(username) > 0 {
		err := tx.Model(&brewer).Update("username", &username).Error
		if err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer username")
		}
	}

	if rank != nil {
		err := tx.Model(&brewer).Association("Rank").Replace(rank).Error
		if err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer rank")
		}
	}

	if len(beerIDs) > 0 {
		beers := []Beer{}
		bIDs := strings.Split(beerIDs, ",")
		err := tx.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error
		if err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
		}

		if err := tx.Model(&brewer).Association("Beers").Replace(&beers).Error; err != nil {
			tx.Rollback()
			return dbWithError(err, http.StatusInternalServerError, "Error updating brewer's beers")
		}
	}

	tx.Commit()
	return dbSuccess(&brewer)
}

// DeleteBrewer func
func DeleteBrewer(id string) *utils.Result {
	brewer := Brewer{}

	tx := db.Begin()
	err := tx.Model(&Brewer{}).Where("id = ?", id).Find(&brewer).Error
	if err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusNotFound, "Error fetching brewer from DB")
	}

	err = tx.Delete(&brewer).Error
	if err != nil {
		tx.Rollback()
		return dbWithError(err, http.StatusInternalServerError, "Error deleting brewer from DB")
	}

	tx.Commit()
	return dbSuccess("Successfully deleted brewer from DB")
}

// GetRankedBrewers func
func GetRankedBrewers(level, limit, order, offset string) *utils.Result {
	brewers := []Brewer{}

	err := db.Joins("JOIN ranks on ranks.id = brewers.rank_id").Where("ranks.level = ?", level).
		Limit(limit).Order("created_at " + order).Offset(offset).
		Preload("Rank").
		Find(&brewers).Error

	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching beers from DB")
	}

	return dbSuccess(brewers)
}

// GetFeaturedBrewers func
func GetFeaturedBrewers(feat, limit, order string) *utils.Result {
	brewers := []Brewer{}

	err := db.Model(&Brewer{}).Limit(limit).Order("created_at "+order).Preload("Beers").Preload("Rank").
		Where("featured = ?", feat).Find(&brewers).Error

	if err != nil {
		return dbWithError(err, http.StatusNotFound, "Error fetching featured brewers from DB")
	}

	return dbSuccess(brewers)
}

// ************************************************************ UTILITY FUNCTIONS ************************************************************ //

// setRankOrNull func
func setRankOrNull(r string, tx *gorm.DB) (*Rank, error) {
	if len(r) > 0 {
		rank := Rank{}

		err := tx.Model(&Rank{}).Where("level = ?", r).Find(&rank).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		return &rank, nil
	}
	return nil, nil
}

// setStringValOrNil func
func setStringValOrNil(v string) *string {
	if len(v) > 0 {
		return &v
	}
	return nil
}
