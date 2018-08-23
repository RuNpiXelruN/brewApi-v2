package api

import (
	"go_apps/go_api_apps/brewApi-v2/db"
	"net/http"

	"github.com/gorilla/mux"
)

type rank struct{}

func (rn rank) registerRoutes(r *mux.Router) {
	r.Path("/ranks/{rank:[1-8]}").HandlerFunc(rn.getBrewersOfRank).Methods("GET") // GET /ranks/:level ?:(limit|order|offset)
	r.Path("/ranks").HandlerFunc(rn.getRanks).Methods("GET")                      // GET /ranks?:(limit|order|offset)
}

// GET /ranks/:level ?:(limit|order|offset)
func (rn rank) getBrewersOfRank(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rank := vars["rank"]
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := db.GetBrewersOfRank(rank, limit, order, offset)
	Respond(w, result)
}

// GET /ranks?:(limit|order|offset)
func (rn rank) getRanks(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := db.GetRanks(limit, order, offset)
	Respond(w, result)
}
