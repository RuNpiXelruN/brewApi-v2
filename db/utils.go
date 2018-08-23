package db

import (
	"go_apps/go_api_apps/brewApi-v2/utils"
	"mime/multipart"
	"net/http"
)

// ReqImage struct
type ReqImage struct {
	File   multipart.File
	Header *multipart.FileHeader
	Error  error
}

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

func dbWithError(err error, code int, text string) *utils.Result {
	result := utils.Result{}
	result.Error = &utils.Error{
		StatusCode: code,
		StatusText: http.StatusText(code) + " - " + text + " : " + err.Error(),
	}
	return &result
}

func dbSuccess(data interface{}) *utils.Result {
	result := utils.Result{}

	result.Success = &utils.Success{
		StatusCode: http.StatusOK,
		Data:       &data,
	}

	return &result
}
