package main

import (
	"flag"
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/config"
	"go_apps/go_api_apps/brewApi-v2/src/controller"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// TODO use db.FirstOrCreate

func init() {
	config.SetVars()
}

func main() {
	seed := flag.Bool("seed", false, "Include to seed the DB")
	migrate := flag.Bool("migrate", false, "Include to migrate DB tables")
	flag.Parse()

	db := config.SetupDatabase(*seed, *migrate)
	defer db.Close()

	r := mux.NewRouter()
	controller.Startup(r)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedHeaders:   []string{"auth_token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	handler := cors.Handler(r)

	r.Path("/").HandlerFunc(index).Methods("GET")

	go http.ListenAndServe(":8000", handler)
	go fmt.Println("..listening on port :8000")
	fmt.Scanln()
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index hit!")
}
