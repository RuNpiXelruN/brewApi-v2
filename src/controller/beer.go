package controller

import (
	"net/http"

	"go_apps/go_api_apps/brewApi-v2/src/model"
	"go_apps/go_api_apps/brewApi-v2/src/utils"

	"github.com/gorilla/mux"
)

// TODO add errors struct to return with success object for non-critical errors eg. image upload fail

type beer struct{}

func (b beer) registerRoutes(r *mux.Router) {
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.getBeer)).Methods("GET")                                                                         // GET /beers/:id
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.updateBeer, model.CheckBeerNameUpdateIsUnique(), model.CheckUserAuth())).Methods("PUT", "PATCH") // PUT/PATCH /beers/:id
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.deleteBeer, model.CheckUserAuth())).Methods("DELETE")                                            // DELETE /beer/:id
	r.Path("/beers").Queries("status", "{status:(?:upcoming|brewing|active|past)}").HandlerFunc(b.getBeersWithStatus).Methods("GET")                        // GET /beers?:status
	r.Path("/beers").Queries("featured", "{featured:(?:true|false)}").HandlerFunc(b.getFeaturedBeers).Methods("GET")                                        // GET /beers?:featured
	r.Path("/beers").HandlerFunc(b.getBeers).Methods("GET")                                                                                                 // GET /beers
	r.Path("/beers").HandlerFunc(utils.Adapt(b.createBeer, model.CheckBeerNameIsUnique(), model.CheckUserAuth())).Methods("POST")                           // POST /beers
}

// GET /beers/:id
func (b beer) getBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	includeBrewers := req.FormValue("include_brewers")

	result := model.GetBeer(id, includeBrewers)
	Response(w, result)
}

// GET /beers
func (b beer) getBeers(w http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	order := req.FormValue("order")
	offset := req.FormValue("offset")

	result := model.GetBeers(limit, order, offset)
	Response(w, result)
}

// POST /beers
func (b beer) createBeer(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	description := req.FormValue("description")
	status := req.FormValue("status")
	alc := req.FormValue("alcohol_content")
	feat := req.FormValue("featured")
	brewerIDs := req.FormValue("brewer_ids")
	mfile, mheader, err := req.FormFile("image")
	image := model.ReqImage{
		File:   mfile,
		Header: mheader,
		Error:  err,
	}

	result := model.CreateBeer(name, description, status, alc, feat, brewerIDs, image)
	Response(w, result)
}

// PUT/PATCH /beers/:id
func (b beer) updateBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	name := req.FormValue("name")
	description := req.FormValue("description")
	status := req.FormValue("status")
	alc := req.FormValue("alcohol_content")
	feat := req.FormValue("featured")
	brewerIDs := req.FormValue("brewer_ids")
	mfile, mheader, err := req.FormFile("image")
	image := model.ReqImage{
		File:   mfile,
		Header: mheader,
		Error:  err,
	}
	result := model.UpdateBeer(id, name, description, status, alc, feat, brewerIDs, image)
	Response(w, result)
}

// DELETE /beer/:id
func (b beer) deleteBeer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	result := model.DeleteBeer(id)
	Response(w, result)
}

// GET /beers?:status
func (b beer) getBeersWithStatus(w http.ResponseWriter, req *http.Request) {
	status := req.FormValue("status")
	limit := req.FormValue("limit")
	order := req.FormValue("order")

	result := model.GetBeersWithStatus(status, limit, order)
	Response(w, result)
}

// GET /beers?:featured
func (b beer) getFeaturedBeers(w http.ResponseWriter, req *http.Request) {
	feat := req.FormValue("featured")
	limit := req.FormValue("limit")
	order := req.FormValue("order")

	result := model.GetFeaturedBeers(feat, limit, order)
	Response(w, result)
}
