package model

import (
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ******************************** BREWER MIDDLEWARE ******************************** //

// CheckPresenceOfFirstNameOrUsername func for ensuring names included
func CheckPresenceOfFirstNameOrUsername() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if len(req.FormValue("first_name")) < 1 && len(req.FormValue("username")) < 1 {
				result := utils.Result{}
				result.Error = &utils.Error{
					Status:     http.StatusBadRequest,
					StatusText: http.StatusText(http.StatusBadRequest) + ": Must include either first name or last name.",
				}
				utils.Response(w, &result)
				return
			}

			h.ServeHTTP(w, req)
		}
	}
}

// CheckUsernameIsUnique func for unique Brewer Usernames
func CheckUsernameIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			username := req.FormValue("username")
			if len(username) > 0 {
				brewer := Brewer{}
				newUsername := db.Model(&Brewer{}).Where(&Brewer{Username: &username}).First(&brewer).RecordNotFound()
				if newUsername != true {
					result := utils.Result{}
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

// CheckBrewerUsernameIsUnique func
func CheckBrewerUsernameIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			username := req.FormValue("username")
			if len(username) > 0 {
				brewer := Brewer{}
				isNewUsername := db.Model(&Brewer{}).Where(&Brewer{Username: &username}).Find(&brewer).RecordNotFound()
				if isNewUsername != true {
					result := utils.Result{}
					result.Error = &utils.Error{
						Status:     http.StatusBadRequest,
						StatusText: http.StatusText(http.StatusBadRequest) + ": Username already taken.",
					}
					utils.Response(w, &result)
					return
				}
			}

			h.ServeHTTP(w, req)
		}
	}
}

// ******************************** BEER MIDDLEWARE ******************************** //

// CheckBeerNameIsUnique func
func CheckBeerNameIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			name := req.FormValue("name")
			beer := Beer{}
			newName := db.Model(&Beer{}).Where(&Beer{Name: name}).First(&beer).RecordNotFound()
			if newName != true {
				result := utils.Result{}
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

// CheckBeerNameUpdateIsUnique func
func CheckBeerNameUpdateIsUnique() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)
			id := vars["id"]
			dbBeer := Beer{}
			var dbName []string
			if err := db.Model(&Beer{}).Where("id = ?", id).Find(&dbBeer).Pluck("name", &dbName).Error; err != nil {
				result := utils.Result{}
				result.Error = &utils.Error{
					Status:     http.StatusBadRequest,
					StatusText: http.StatusText(http.StatusBadRequest) + ": Beer with that ID not found",
				}
				utils.Response(w, &result)
				return
			}

			name := req.FormValue("name")
			var beer Beer

			if len(name) != 0 && name != dbName[0] {
				newName := db.Model(&Beer{}).Where(&Beer{Name: name}).First(&beer).RecordNotFound()
				if newName != true {
					result := utils.Result{}
					result.Error = &utils.Error{
						Status:     http.StatusBadRequest,
						StatusText: http.StatusText(http.StatusBadRequest) + ": Beer name already exists",
					}

					utils.Response(w, &result)
					return
				}
			}

			h.ServeHTTP(w, req)
		}
	}
}

// ******************************** AUTH MIDDLEWARE ******************************** //

// CheckUserAuth func
func CheckUserAuth() utils.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			log.Println("Auth Passed ☺")

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
