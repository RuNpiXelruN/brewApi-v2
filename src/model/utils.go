package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"
)

// UploadImage func
func UploadImage(image ReqImage) *utils.Result {
	var imgURL string

	if image.Error != nil {
		return dbWithError(image.Error, http.StatusBadRequest, "Image not valid.")
	}

	imageURL, err := s3ImageUpload(image)
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error uploading image to S3")
	}

	imgURL = imageURL.(string)
	return dbSuccess(imgURL)
}
