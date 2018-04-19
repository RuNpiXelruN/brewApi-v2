package controller

import (
	"github.com/gorilla/mux"
)

var (
	beerController beer
)

func Startup(r *mux.Router) {
	beerController.registerRoutes(r)
}
