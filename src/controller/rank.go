package controller

import (
	"go_apps/go_api_apps/brewApi-v2/src/model"
	"net/http"

	"github.com/gorilla/mux"
)

type rank struct{}

func (rn rank) registerRoutes(r *mux.Router) {
	r.Path("/ranks").Queries("level", "{level:[1-8]}").HandlerFunc(rn.getBrewersOfRank).Methods("GET")
	r.Path("/ranks").HandlerFunc(rn.getRanks).Methods("GET")
}

// GET /ranks?:(level)
func (rn rank) getBrewersOfRank(w http.ResponseWriter, req *http.Request) {
	level := req.FormValue("level")

	result := model.GetBrewersOfRank(level)
	Response(w, result)
}

// GET /ranks?:(limit|order|offset)
func (rn rank) getRanks(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := model.GetRanks(limit, order, offset)
	Response(w, result)
}
