package api

import (
	"go_apps/go_api_apps/brewApi-v2/db"
	"go_apps/go_api_apps/brewApi-v2/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type brewer struct{}

func (b brewer) registerRoutes(r *mux.Router) {
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(b.getBrewer).Methods("GET")                                               // GET /brewers/:id
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.deleteBrewer, utils.CheckToken())).Methods("DELETE")        // DELETE /brewers/:id
	r.Path("/brewers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.updateBrewer, utils.CheckToken())).Methods("PUT", "PATCH")  // PUT/PATCH /brewers/:id
	r.Path("/brewers/basic").HandlerFunc(b.getBrewerNames).Methods("GET")                                                // GET /brewers/basic
	r.Path("/brewers").Queries("featured", "{featured:(?:true|false)}").HandlerFunc(b.getFeaturedBrewers).Methods("GET") // GET /brewers?:featured
	r.Path("/brewers").Queries("rank", "{rank:[1-8]}").HandlerFunc(b.getRankedBrewers).Methods("GET")                    // GET /brewers/:rank
	r.Path("/brewers").HandlerFunc(b.getBrewers).Methods("GET")                                                          // GET /brewers
	r.Path("/brewers").HandlerFunc(utils.Adapt(b.createBrewer, utils.CheckToken())).Methods("POST")                      // POST /brewers
}

// GET /brewers/basic
func (b brewer) getBrewerNames(w http.ResponseWriter, req *http.Request) {
	result := db.GetBrewerNames()
	Respond(w, result)
}

// GET /brewers
func (b brewer) getBrewers(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := db.GetBrewers(limit, order, offset)
	Respond(w, result)
}

// GET /brewers/:id
func (b brewer) getBrewer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	includeBeers := req.FormValue("include_beers")

	result := db.GetBrewer(id, includeBeers)
	Respond(w, result)
}

// POST /brewers
func (b brewer) createBrewer(w http.ResponseWriter, req *http.Request) {
	first := req.FormValue("first_name")
	last := req.FormValue("last_name")
	feat := req.FormValue("featured")
	username := req.FormValue("username")
	rank := req.FormValue("rank")
	beerIDs := req.FormValue("beer_ids")

	result := db.CreateBrewer(first, last, feat, username, rank, beerIDs)
	Respond(w, result)
}

// PUT/PATCH /brewers/:id
func (b brewer) updateBrewer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	first := req.FormValue("first_name")
	last := req.FormValue("last_name")
	feat := req.FormValue("featured")
	username := req.FormValue("username")
	rank := req.FormValue("rank")
	beerIDs := req.FormValue("beer_ids")

	result := db.UpdateBrewer(id, first, last, feat, username, rank, beerIDs)
	Respond(w, result)
}

// DELETE /brewers/:id
func (b brewer) deleteBrewer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := db.DeleteBrewer(id)
	Respond(w, result)
}

// GET /brewers/:rank
func (b brewer) getRankedBrewers(w http.ResponseWriter, req *http.Request) {
	rankLevel := req.FormValue("rank")
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := db.GetRankedBrewers(rankLevel, limit, order, offset)
	Respond(w, result)
}

// GET /brewers?:featured
func (b brewer) getFeaturedBrewers(w http.ResponseWriter, req *http.Request) {
	feat := req.FormValue("featured")
	limit := req.FormValue("limit")
	order := req.FormValue("order")

	result := db.GetFeaturedBrewers(feat, limit, order)
	Respond(w, result)
}
