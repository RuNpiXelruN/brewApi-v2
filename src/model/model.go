package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"mime/multipart"
	"net/http"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// ReqImage struct
type ReqImage struct {
	File   multipart.File
	Header *multipart.FileHeader
	Error  error
}

func dbWithError(err error, code int, text string) *utils.Result {
	result.Error = &utils.Error{
		Status:     code,
		StatusText: http.StatusText(code) + " - " + text + " : " + err.Error(),
	}
	return &result
}

func dbSuccess(data interface{}) *utils.Result {
	result.Success = &utils.Success{
		Status: http.StatusOK,
		Data:   &data,
	}
	return &result
}
