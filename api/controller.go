package api

import (
	"encoding/json"
	"go_apps/go_api_apps/brewApi-v2/utils"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	beerController   beer
	brewerController brewer
	rankController   rank
	utilsController  utilRoutes
)

// Startup func to register model routes
func Startup(r *mux.Router) {
	beerController.registerRoutes(r)
	brewerController.registerRoutes(r)
	rankController.registerRoutes(r)
	utilsController.registerRoutes(r)
}

// Respond func
func Respond(w http.ResponseWriter, result *utils.Result) {
	if result.Error != nil {
		w.WriteHeader(result.Error.StatusCode)
		w.Header().Set("Status", http.StatusText(result.Error.StatusCode))

		data, _ := json.Marshal(result.Error)
		w.Write(data)
		return
	}

	w.WriteHeader(result.Success.StatusCode)
	w.Header().Set("Status", http.StatusText(result.Success.StatusCode))

	data, _ := json.Marshal(result.Success.Data)
	w.Write(data)
}
