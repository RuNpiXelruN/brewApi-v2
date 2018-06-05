package controller

import (
	"log"
	"net/http"

	"go_apps/go_api_apps/brewApi-v2/src/model"
	"go_apps/go_api_apps/brewApi-v2/src/utils"

	"github.com/gorilla/mux"
)

// TODO add errors struct to return with success object for non-critical errors eg. image upload fail

type beer struct{}

func (b beer) registerRoutes(r *mux.Router) {
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.getBeer)).Methods("GET")
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(utils.Adapt(b.updateBeer, model.CheckBeerNameIsUnique())).Methods("PUT", "PATCH")
	r.Path("/beers/{id:[0-9]+}").HandlerFunc(b.deleteBeer).Methods("DELETE")
	r.Path("/beers").Queries("status", "{status:(?:upcoming|brewing|active-full|active-empty|past)}").HandlerFunc(b.getStatusBeers).Methods("GET")
	r.Path("/beers").Queries("featured", "{featured:(?:true|false)}").HandlerFunc(b.getFeaturedBeers).Methods("GET")
	r.Path("/beers").HandlerFunc(utils.Adapt(b.getBeersHandler, model.SayHi(), utils.SayHi())).Methods("GET")
	// r.Path("/beers").HandlerFunc(utils.Adapt(b.createBeerWithChannels, model.SayHi())).Methods("POST")
	r.Path("/beers").HandlerFunc(b.createBeer).Methods("POST")
	// r.Path("/beers").HandlerFunc(utils.Adapt(b.createBeerHandler, model.CheckBeerNameIsUnique())).Methods("POST")
}

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

func (b beer) createBeerWithChannels(w http.ResponseWriter, req *http.Request) {

	name := req.FormValue("name")
	description := req.FormValue("description")
	status := req.FormValue("status")
	alc := req.FormValue("alc")
	feat := req.FormValue("feat")
	brewerIDs := req.FormValue("brewer_ids")
	mfile, mheader, err := req.FormFile("image")
	image := model.ReqImage{
		File:   mfile,
		Header: mheader,
		Error:  err,
	}

	result := model.CreateBeerWithChannels(name, description, status, alc, feat, brewerIDs, image)
	Response(w, result)
}

// POST /beers
// func (b beer) createBeerHandler(w http.ResponseWriter, req *http.Request) {

// 	// CHECK UNIQUE NAME FUNCTION

// 	done := make(chan interface{})
// 	defer fmt.Println("done channel closed.")
// 	defer close(done)

// 	var result *utils.Result
// 	var imageURL string

// 	name := req.FormValue("name")
// 	desc := req.FormValue("description")
// 	alc := req.FormValue("alcohol_content")
// 	feat := req.FormValue("featured")
// 	brewerIDs := req.FormValue("brewer_ids")

// 	multiFile, multiHeader, err := req.FormFile("image")
// 	if err != nil {
// 		log.Println("Error getting image from request ->", err.Error())

// 		result = model.CreateBeer(name, desc, alc, feat, brewerIDs, imageURL)
// 		Response(w, result)
// 		return
// 	}

// 	s3ResultChan := utils.UploadToS3(done, multiFile, multiHeader)

// 	select {
// 	case s3Result := <-s3ResultChan:
// 		if s3Result.Error != nil {
// 			log.Println("Error uploading image to s3 ->", s3Result.Error.StatusText)
// 			break
// 		}
// 		imageURL = s3Result.Success.Data.(string)
// 		result = model.CreateBeer(name, desc, alc, feat, brewerIDs, imageURL)
// 		Response(w, result)
// 		return
// 	}

// 	result = model.CreateBeer(name, desc, alc, feat, brewerIDs, imageURL)
// 	Response(w, result)
// }

// PUT/PATCH /beers/:id
func (b beer) updateBeer(w http.ResponseWriter, req *http.Request) {
	var result *utils.Result
	done := make(chan interface{})
	defer log.Println("done cannel closed.")
	defer close(done)

	var imageURL string

	vars := mux.Vars(req)
	id := vars["id"]
	name := req.FormValue("name")
	desc := req.FormValue("description")
	stat := req.FormValue("status")
	alc := req.FormValue("alcohol_content")
	feat := req.FormValue("featured")
	brewIDs := req.FormValue("brewer_ids")

	multifile, multiheader, err := req.FormFile("image")
	if err != nil {
		// There is no image for the update
		log.Println("Error getting image from request ->", err.Error())

		result = model.UpdateBeer(id, name, desc, stat, alc, feat, brewIDs, imageURL)
		Response(w, result)
		return
	}

	s3ResultChan := utils.UploadToS3(done, multifile, multiheader)

	select {
	case s3Result := <-s3ResultChan:
		if s3Result.Error != nil {
			log.Println("Error uploading image to S3:", s3Result.Error.StatusText)
			break
		}
		imageURL = s3Result.Success.Data.(string)
		result = model.UpdateBeer(id, name, desc, stat, alc, feat, brewIDs, imageURL)
		Response(w, result)
		return
	}

	result = model.UpdateBeer(id, name, desc, stat, alc, feat, brewIDs, imageURL)
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
