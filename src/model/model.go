package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func dbWithError(err error, code int, text string) *utils.Result {
	result.Error = &utils.Error{
		Status:     code,
		StatusText: http.StatusText(code) + " - " + text + " : " + err.Error(),
	}
	return &result
}
