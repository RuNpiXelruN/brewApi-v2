package controller

import (
	"net/http"

	"go_apps/go_api_apps/brewApi-v2/src/model"

	"github.com/gorilla/mux"
)

type utilRoutes struct{}

func (ur utilRoutes) registerRoutes(r *mux.Router) {
	r.Path("/image-uploader").HandlerFunc(ur.uploadImage).Methods("POST")
}

func (ur utilRoutes) uploadImage(w http.ResponseWriter, req *http.Request) {
	mfile, mheader, err := req.FormFile("image")
	image := model.ReqImage{
		File:   mfile,
		Header: mheader,
		Error:  err,
	}

	result := model.UploadImage(image)
	Response(w, result)
}
