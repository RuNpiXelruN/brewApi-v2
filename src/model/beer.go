package model

import (
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var result utils.Result

func getBrewers(brewerIDs string) <-chan chanResult {
	out := make(chan chanResult)
	brewers := []Brewer{}
	bIDs := strings.Split(brewerIDs, ",")

	go func() {
		if err := db.Model(&Brewer{}).Preload("Rank").Where("id in (?)", bIDs).Find(&brewers).Error; err != nil {
			out <- chanResult{nil, err}
		}
		out <- chanResult{brewers, nil}
	}()
	return out
}

// CreateBeerWithChannels func
func CreateBeerWithChannels(name, desc, status, alcContent, ft, brewerIDs string, image ReqImage) *utils.Result {
	// var s3Img string
	var s3Result chan chanResult

	if image.Error != nil {
		s3Result = nil
	} else {
		// s3Result = uploadToS3(image.File, image.Header)
		uploadToS3(image.File, image.Header, s3Result)
	}

	var brewers []Brewer
	var brewersResult <-chan chanResult

	if len(brewerIDs) == 0 {
		brewersResult = nil
	} else {
		brewersResult = getBrewers(brewerIDs)
	}

	fmt.Println("here?", <-s3Result)

	for s3Result != nil || brewersResult != nil {
		select {
		case s3Res := <-s3Result:
			fmt.Println("hithit")
			if s3Res.Error != nil {
				return dbWithError(s3Res.Error, http.StatusInternalServerError, "Error uploading image to S3")
			}
			fmt.Printf("%+v\n", s3Res)
			// s3Img = s3Res.Data.(string)
			s3Result = nil
		case brewResult := <-brewersResult:
			if brewResult.Error != nil {
				return dbWithError(brewResult.Error, http.StatusNotFound, "Error fetching brewers from DB")
			}
			brewers = brewResult.Data.([]Brewer)
			brewersResult = nil
		default:
		}
	}

	// alc, _ := strconv.ParseFloat(alcContent, 64)

	// beer := Beer{
	// 	Name:        name,
	// 	Description: desc,
	// 	Status: status,
	// 	AlcoholContent: alc,
	// }
	return dbSuccess(brewers)
}

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

func s3ImageUpload(img ReqImage) (interface{}, error) {
	if img.Error != nil {
		return nil, img.Error
	}
	newFile, err := os.Create(img.Header.Filename)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, img.File)
	if err != nil {
		return nil, err
	}

	err = newFile.Sync()
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-southeast-2")}))
	url, err := uploader.Upload(&s3manager.UploadInput{
		Body:        newFile,
		Bucket:      aws.String("brew-site"),
		Key:         aws.String(newFile.Name()),
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(newFile.Name()))),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return nil, err
	}

	return url.Location, nil
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

// CreateBeer func
// func CreateBeer(name, desc, alc, feat, brewIDs, imageURL string) *utils.Result {
// 	al, _ := strconv.ParseFloat(alc, 64)
// 	ft, _ := strconv.ParseBool(feat)

// 	bIDs := strings.Split(brewIDs, ",")

// 	var brUintIDs []uint
// 	for _, b := range bIDs {
// 		intID, _ := strconv.Atoi(b)
// 		brUintIDs = append(brUintIDs, uint(intID))
// 	}

// 	brewers := []Brewer{}
// 	if err := db.Model(&Brewer{}).Where(brUintIDs).Find(&brewers).Error; err != nil {
// 		result.Error = &utils.Error{
// 			Status:     http.StatusNotFound,
// 			StatusText: http.StatusText(http.StatusNotFound) + " - Error fetching brewers from DB",
// 		}
// 		return &result
// 	}

// 	beer := Beer{
// 		Name:           name,
// 		Description:    desc,
// 		AlcoholContent: al,
// 		Featured:       ft,
// 		ImageURL:       imageURL,
// 		Brewers:        brewers,
// 		CreatedAt:      time.Now(),
// 	}

// 	if err := db.Save(&beer).Error; err != nil {
// 		result.Error = &utils.Error{
// 			Status:     http.StatusInternalServerError,
// 			StatusText: http.StatusText(http.StatusInternalServerError) + " - Error saving beer to DB",
// 		}
// 		return &result
// 	}

// 	result.Success = &utils.Success{
// 		Status: http.StatusCreated,
// 		Data:   &beer,
// 	}
// 	return &result
// }

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

func uploadToS3(file multipart.File, header *multipart.FileHeader, ch chan chanResult) {
	// out := make(chan chanResult)
	// signal := make(chan bool)
	// var wg sync.WaitGroup

	// wg.Add(1)
	go func() {
		// defer wg.Done()
		newFile, err := os.Create(header.Filename)
		if err != nil {
			ch <- chanResult{nil, err}
			return
			// signal <- true
		}
		defer file.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			ch <- chanResult{nil, err}
			return
			// signal <- true
		}

		err = newFile.Sync()
		if err != nil {
			ch <- chanResult{nil, err}
			return
			// signal <- true
		}
		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-southeast-2")}))
		url, err := uploader.Upload(&s3manager.UploadInput{
			Body:        newFile,
			Bucket:      aws.String("brew-site"),
			Key:         aws.String(newFile.Name()),
			ContentType: aws.String(mime.TypeByExtension(filepath.Ext(newFile.Name()))),
			ACL:         aws.String("public-read"),
		})
		if err != nil {
			ch <- chanResult{nil, err}
			return
			// signal <- true
		}
		ch <- chanResult{
			Data:  url.Location,
			Error: nil,
		}

		return
		// signal <- true
	}()

	// go func() {
	// 	wg.Wait()
	// 	fmt.Println("finnn")
	// }()

}
