package api

import (
	"net/http"

	"go_apps/go_api_apps/brewApi-v2/db"
	"go_apps/go_api_apps/brewApi-v2/utils"

	"github.com/gorilla/mux"
)

type utilRoutes struct{}

func (ur utilRoutes) registerRoutes(r *mux.Router) {
	r.Path("/image-uploader").HandlerFunc(utils.Adapt(ur.uploadImage, utils.CheckToken())).Methods("POST")
}

func (ur utilRoutes) uploadImage(w http.ResponseWriter, req *http.Request) {
	mfile, mheader, err := req.FormFile("image")
	image := db.ReqImage{
		File:   mfile,
		Header: mheader,
		Error:  err,
	}

	result := db.UploadImage(image)
	Respond(w, result)
}
