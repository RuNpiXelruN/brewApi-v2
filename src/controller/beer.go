package controller

import (
	"net/http"

	"go_apps/go_api_apps/brewApi-v2/src/model"

	"github.com/gorilla/mux"
)

type beer struct{}

func (b beer) registerRoutes(r *mux.Router) {
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(b.getBeer).Methods("GET")
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(b.updateBeer).Methods("PUT", "PATCH")
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(b.deleteBeer).Methods("DELETE")
	r.Path("/beers").Queries("status", "{status:(?:upcoming|brewing|active-full|active-empty|past)}").HandlerFunc(b.getStatusBeers).Methods("GET")
	r.Path("/beers").Queries("featured", "{featured:(?:true|false)}").HandlerFunc(b.getFeaturedBeers).Methods("GET")
	r.Path("/beers").HandlerFunc(b.getBeersHandler).Methods("GET")
	r.Path("/beers").HandlerFunc(b.createBeerHandler).Methods("POST")
}

// TODO validate in middleware that name is unique
// TODO add S3 image upload

// PUT/PATCH /beers/:id
func (b beer) updateBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	name := req.FormValue("name")
	desc := req.FormValue("description")
	stat := req.FormValue("status")
	alc := req.FormValue("alcohol_content")
	feat := req.FormValue("featured")
	brewIDs := req.FormValue("brewer_ids")

	result := model.UpdateBeer(id, name, desc, stat, alc, feat, brewIDs)
	Response(w, result)
}

// GET /beers?:status
func (b beer) getStatusBeers(w http.ResponseWriter, req *http.Request) {
	status := req.FormValue("status")

	result := model.GetBeersWithStatus(status)
	Response(w, result)
}

// GET /beers?:featured
func (b beer) getFeaturedBeers(w http.ResponseWriter, req *http.Request) {
	feat := req.FormValue("featured")

	result := model.GetFeaturedBeers(feat)
	Response(w, result)
}

// POST /beers
func (b beer) createBeerHandler(w http.ResponseWriter, req *http.Request) {

	// CHECK UNIQUE NAME FUNCTION

	name := req.FormValue("name")
	desc := req.FormValue("description")
	alc := req.FormValue("alcohol_content")
	feat := req.FormValue("featured")
	brewerIDs := req.FormValue("brewer_ids")

	result := model.CreateBeer(name, desc, alc, feat, brewerIDs)
	Response(w, result)
}

// GET /beers/:id
func (b beer) getBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := model.GetBeer(id)
	Response(w, result)
}

// DELETE /beers/:id
func (b beer) deleteBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := model.DeleteBeer(id)
	Response(w, result)
}

// GET /beers
func (b beer) getBeersHandler(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := model.GetBeers(limit, order, offset)
	Response(w, result)
}
