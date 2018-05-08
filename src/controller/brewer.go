package controller

import (
	"go_apps/go_api_apps/brewApi-v2/src/model"
	"go_apps/go_api_apps/brewApi-v2/src/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type brewer struct{}

func (b brewer) registerRoutes(r *mux.Router) {
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(b.getBrewer).Methods("GET")
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(b.deleteBrewer).Methods("DELETE")
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.updateBrewerWithChannels, model.SayHi())).Methods("PUT", "PATCH")
	r.Path("/brewers").Queries("featured", "{featured:(?:true|false)}").HandlerFunc(b.getFeaturedBrewers).Methods("GET")
	r.Path("/brewers").Queries("rank", "{rank:[1-8]}").HandlerFunc(b.getRankedBrewers).Methods("GET")
	r.Path("/brewers").HandlerFunc(b.getBrewers).Methods("GET")
	r.Path("/brewers").HandlerFunc(utils.Adapt(b.createBrewerWithChannels, model.CheckUsernameIsUnique(), utils.SayHi())).Methods("POST")
}

// POST /brewers
func (b brewer) createBrewerWithChannels(w http.ResponseWriter, req *http.Request) {
	first := req.FormValue("first_name")
	last := req.FormValue("last_name")
	feat := req.FormValue("featured")
	username := req.FormValue("username")
	rank := req.FormValue("rank")
	beerIDs := req.FormValue("beer_ids")

	result := model.CreateBrewerWithChannels(first, last, feat, username, rank, beerIDs)
	Response(w, result)
}

// PUT|PATCH /brewers/:id
func (b brewer) updateBrewerWithChannels(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	first := req.FormValue("first_name")
	last := req.FormValue("last_name")
	feat := req.FormValue("featured")
	username := req.FormValue("username")
	rank := req.FormValue("rank")
	beerIDs := req.FormValue("beer_ids")

	result := model.UpdateBrewerWithChannels(id, first, last, username, feat, rank, beerIDs)
	Response(w, result)
}

// DELETE /brewers/:id
func (b brewer) deleteBrewer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := model.DeleteBrewer(id)
	Response(w, result)
}

// GET /brewers/:rank
func (b brewer) getRankedBrewers(w http.ResponseWriter, req *http.Request) {
	rankLevel := req.FormValue("rank")

	result := model.GetRankedBrewers(rankLevel)
	Response(w, result)
}

// GET /brewers?:featured
func (b brewer) getFeaturedBrewers(w http.ResponseWriter, req *http.Request) {
	feat := req.FormValue("featured")

	result := model.GetFeaturedBrewers(feat)
	Response(w, result)
}

// GET /brewer/:id
func (b brewer) getBrewer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := model.GetBrewer(id)
	Response(w, result)
}

// GET /brewers
func (b brewer) getBrewers(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := model.GetBrewers(limit, order, offset)
	Response(w, result)
}
