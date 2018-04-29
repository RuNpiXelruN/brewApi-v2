package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"log"
	"net/http"
)

// CheckUsernameIsUnique func for unique Brewer Usernames
func CheckUsernameIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			username := req.FormValue("username")
			if len(username) > 0 {
				brewer := Brewer{}
				newUsername := db.Model(&Brewer{}).Where(&Brewer{Username: &username}).First(&brewer).RecordNotFound()
				if newUsername != true {
					result.Error = &utils.Error{
						Status:     http.StatusBadRequest,
						StatusText: http.StatusText(http.StatusBadRequest) + ": Username already exists.",
					}
					utils.Response(w, &result)
					return
				}
			}

			h.ServeHTTP(w, req)
		}
	}
}

// CheckBeerNameIsUnique func
func CheckBeerNameIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			name := req.FormValue("name")
			beer := Beer{}
			newName := db.Model(&Beer{}).Where(&Beer{Name: name}).First(&beer).RecordNotFound()
			if newName != true {
				result.Error = &utils.Error{
					Status:     http.StatusBadRequest,
					StatusText: http.StatusText(http.StatusBadRequest) + ": Beer name already exists.",
				}
				utils.Response(w, &result)
				return
			}

			h.ServeHTTP(w, req)
		}
	}
}

// SayHi func to test working middleware
func SayHi() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			log.Println("Hi from Model Middleware!")

			h.ServeHTTP(w, req)
		}
	}
}
