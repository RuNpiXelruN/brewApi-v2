package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type beer struct{}

func (b beer) registerRoutes(r *mux.Router) {
	r.Path("/beers").HandlerFunc(b.getBeersHandler).Methods("GET")
}

func (b beer) getBeersHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hittt")
}
