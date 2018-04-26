package controller

import (
	"github.com/gorilla/mux"
)

var (
	beerController   beer
	brewerController brewer
	rankController   rank
)

// Startup func to register model routes
func Startup(r *mux.Router) {
	beerController.registerRoutes(r)
	brewerController.registerRoutes(r)
	rankController.registerRoutes(r)
}
