package main

import (
	"flag"
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/api"
	"go_apps/go_api_apps/brewApi-v2/config"
	"go_apps/go_api_apps/brewApi-v2/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func init() {
	config.SetVars()
}

func main() {
	seed := flag.Bool("seed", false, "Include to seed the DB")
	migrate := flag.Bool("migrate", false, "Include to migrate DB tables")
	flag.Parse()

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	database := db.Init(*seed, *migrate)
	defer database.Close()

	api.Startup(s)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"auth_token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	handler := cors.Handler(r)

	r.Path("/").HandlerFunc(index).Methods("GET")

	fmt.Println("..listening on port 8000")
	http.ListenAndServe("0.0.0.0:8000", handler)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Index Hit")
}
