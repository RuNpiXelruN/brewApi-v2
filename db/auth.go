package db

import (
	"go_apps/go_api_apps/brewApi-v2/utils"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// GetSessionUser func
func GetSessionUser(sessionValue string) (*User, *string, bool) {
	var user *User
	session := Session{}

	if err := db.Model(&Session{}).Preload("User").Where(Session{Value: sessionValue}).Find(&session).Error; err != nil {
		return nil, nil, false
	}

	// session expired, re-login
	if time.Since(session.Expires) > 0 {
		return nil, nil, false
	}

	user = session.User

	return user, &session.Value, true
}

// HandleAuthCallback func
func HandleAuthCallback(email string, w http.ResponseWriter) *utils.Result {

	tokenHeader, err := utils.GetToken(email)
	if err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error creating jwt")
	}

	var user User
	expires := time.Now().Add(1 * time.Hour)
	sID := uuid.NewV4()

	// check if new user
	if db.Model(&User{}).Preload("Session").Where("email = ?", email).Find(&user).RecordNotFound() {
		user = User{
			Email: &email,
			Session: &Session{
				Value:   sID.String(),
				Expires: expires,
			},
		}
		if err := db.Save(&user).Error; err != nil {
			return dbWithError(err, http.StatusInternalServerError, "Error saving new user to DB")
		}

		return dbSuccess(&user, tokenHeader)
	}

	// otherwise existing user but no active cookie / session,
	// create new one
	session := Session{
		Value:   sID.String(),
		Expires: expires,
	}

	if err := db.Model(&user).Association("Session").Replace(&session).Error; err != nil {
		return dbWithError(err, http.StatusInternalServerError, "Error replacing session in DB")
	}

	return dbSuccess(&user, tokenHeader)
}
