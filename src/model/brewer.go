package model

import (
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var dbError *utils.Result

type chanResult struct {
	Data  interface{}
	Error error
}

func getBrewer(id string, brewCh chan chanResult) {
	brewer := Brewer{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := db.Model(&Brewer{}).Preload("Beers").Preload("Rank").Where("id = ?", id).Find(&brewer).Error; err != nil {
			brewCh <- chanResult{nil, err}
		}
		brewCh <- chanResult{brewer, nil}
	}()
	wg.Wait()
	fmt.Println("fetch brewer completed")
}

func getRank(lvl string, rankCh chan chanResult) {
	rank := Rank{}
	go func() {
		if err := db.Model(&Rank{}).Where("level = ?", lvl).Find(&rank).Error; err != nil {
			rankCh <- chanResult{nil, err}
		}
		rankCh <- chanResult{&rank, nil}
	}()
	fmt.Println("fetch rank completed")
}

func getBeers(beerIDs string, beersCh chan chanResult) {
	bIDs := strings.Split(beerIDs, ",")
	beers := []Beer{}
	go func() {
		err := db.Model(&Beer{}).Where("id in (?)", bIDs).Find(&beers).Error
		if err != nil {
			beersCh <- chanResult{nil, err}
		}
		beersCh <- chanResult{beers, nil}
	}()
	fmt.Println("fetch beers completed")
}

// CreateBrewerWithChannels func
func CreateBrewerWithChannels(first, last, ft, uname, rnk, beerIDs string) *utils.Result {
	var rank *Rank
	var beers []Beer
	feat, _ := strconv.ParseBool(ft)
	rankCh := make(chan chanResult)
	beersCh := make(chan chanResult)

	tx := db.Begin()
	if len(rnk) == 0 {
		rankCh = nil
	} else {
		go getRank(rnk, rankCh)
	}

	if len(beerIDs) == 0 {
		beersCh = nil
	} else {
		go getBeers(beerIDs, beersCh)
	}

	for rankCh != nil || beersCh != nil {
		select {
		case fetchRank := <-rankCh:
			if fetchRank.Error != nil {
				tx.Rollback()
				dbError = dbWithError(fetchRank.Error, http.StatusNotFound, "Error fetching Rank from DB")
				return dbError
			}
			rank = fetchRank.Data.(*Rank)
			rankCh = nil
		case fetchBeers := <-beersCh:
			if fetchBeers.Error != nil {
				tx.Rollback()
				dbError = dbWithError(fetchBeers.Error, http.StatusNotFound, "Error fetching beers from DB")
				return dbError
			}
			beers = fetchBeers.Data.([]Beer)
			beersCh = nil
		default:
		}
	}

	brewer := Brewer{
		FirstName: first,
		LastName:  last,
		Featured:  feat,
		Username:  &uname,
		Rank:      rank,
		Beers:     beers,
	}

	if err := tx.Save(&brewer).Error; err != nil {
		tx.Rollback()
		dbError = dbWithError(err, http.StatusInternalServerError, "Error saving Brewer in DB")
		return dbError
	}

	tx.Commit()
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   brewer,
	}
	return &result
}

// UpdateBrewerWithChannels func
func UpdateBrewerWithChannels(id, f, l, uname, ft, rnk, beerIDs string) *utils.Result {
	var brewer Brewer
	brewCh := make(chan chanResult)
	go getBrewer(id, brewCh)

	var rank *Rank
	var beers []Beer
	rankCh := make(chan chanResult)
	beersCh := make(chan chanResult)

	// all db writes pass together, otherwise all fail together
	tx := db.Begin()

brewLoop:
	for {
		select {
		case fetchBrewer := <-brewCh:
			if fetchBrewer.Error != nil {
				dbError = dbWithError(fetchBrewer.Error, http.StatusNotFound, "Error fetching Brewer from DB")
				return dbError
			}
			brewer = fetchBrewer.Data.(Brewer)

			break brewLoop
		default:
		}
	}

	err := tx.Model(&brewer).Updates(&Brewer{
		FirstName: f,
		LastName:  l,
	}).Error

	if err != nil {
		tx.Rollback()
		dbError = dbWithError(err, http.StatusInternalServerError, "Error updating Brewer")
		return dbError
	}

	if len(uname) > 0 {
		if err := tx.Model(&brewer).Update("username", &uname).Error; err != nil {
			tx.Rollback()
			dbError = dbWithError(err, http.StatusInternalServerError, "Error updating username in DB")
			return dbError
		}
	}

	if len(ft) > 0 {
		feat, _ := strconv.ParseBool(ft)
		if err := tx.Model(&brewer).Update("featured", feat).Error; err != nil {
			tx.Rollback()
			dbError = dbWithError(err, http.StatusInternalServerError, "Error updating featured status in DB")
			return dbError
		}
	}

	if len(rnk) == 0 {
		rankCh = nil
	} else {
		go getRank(rnk, rankCh)
	}

	if len(beerIDs) == 0 {
		beersCh = nil
	} else {
		go getBeers(beerIDs, beersCh)
	}

	for rankCh != nil || beersCh != nil {
		select {
		case fetchRank := <-rankCh:
			if fetchRank.Error != nil {
				dbError = dbWithError(fetchRank.Error, http.StatusNotFound, "Error fetching Rank from DB")
				return dbError
			}
			rank = fetchRank.Data.(*Rank)

			if err := tx.Model(&brewer).Association("Rank").Replace(rank).Error; err != nil {
				tx.Rollback()
				dbError = dbWithError(err, http.StatusInternalServerError, "Error updating Brewer rank")
				return dbError
			}
			rankCh = nil
		case fetchBeers := <-beersCh:
			if fetchBeers.Error != nil {
				dbError = dbWithError(fetchBeers.Error, http.StatusNotFound, "Error fetching Beers from DB")
				return dbError
			}

			beers = fetchBeers.Data.([]Beer)
			if err := tx.Model(&brewer).Association("Beers").Replace(beers).Error; err != nil {
				tx.Rollback()
				dbError = dbWithError(err, http.StatusInternalServerError, "Error updating Brewer's beers")
				return dbError
			}
			beersCh = nil
		default:
		}
	}

	tx.Commit()
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

// CreateBrewer func
func CreateBrewer(first, last, feat, uname, rank, beerIDs string) *utils.Result {
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
		Username:  &uname,
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
			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching Brewer from DB",
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
