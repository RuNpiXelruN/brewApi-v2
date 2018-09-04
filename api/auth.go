package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go_apps/go_api_apps/brewApi-v2/db"
	"go_apps/go_api_apps/brewApi-v2/utils"

	"golang.org/x/oauth2/google"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type auth struct{}

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

var (
	clientID     = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/api/auth/callback",
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

func (a auth) registerRoutes(r *mux.Router) {
	r.Path("/auth").HandlerFunc(a.handleAuth).Methods("GET", "POST")
	r.Path("/auth/callback").HandlerFunc(a.authCallback).Methods("GET", "POST")
}

func (a auth) handleAuth(w http.ResponseWriter, req *http.Request) {
	fmt.Println("AUTH HIT!!!")

	// check if user email already has an active session / check cookie
	result, ok := alreadyLoggedIn(req)
	if ok {
		Respond(w, result)
		return
	}

	// no active session / cookie
	oauthStateString := uniuri.New()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)

	// type urlData struct {
	// 	URL string `json:"url"`
	// }
	// responseURL := urlData{
	// 	URL: url,
	// }

	// Respond(w, dbSuccess(responseURL, nil))
}

func alreadyLoggedIn(req *http.Request) (*utils.Result, bool) {
	sessionVal := req.Header.Get("brew_token")
	if len(sessionVal) < 1 {
		return nil, false
	}

	sessionUser, sessionID, ok := db.GetSessionUser(sessionVal)
	if !ok {
		return nil, false
	}

	result := dbSuccess(sessionUser, sessionID)
	return result, true
}

func (a auth) authCallback(w http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")
	token, _ := googleOauthConfig.Exchange(oauth2.NoContext, code)

	res, err := http.Get(`https://www.googleapis.com/oauth2/v2/userinfo?access_token=` + token.AccessToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer res.Body.Close()

	contents, _ := ioutil.ReadAll(res.Body)
	var user *googleUser
	_ = json.Unmarshal(contents, &user)

	result := db.HandleAuthCallback(user.Email, w)
	Respond(w, result)
}
