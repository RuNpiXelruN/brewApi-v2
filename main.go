package main

import (
	"flag"
	"fmt"
	"go_apps/go_api_apps/brewApi-v2/config"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO use db.FirstOrCreate

func init() {
	config.SetVars()
}

func main() {
	seed := flag.String("seed", "false", "Set to true to call the seed functions")
	seedShort := flag.String("S", "false", "Set to true to call the seed functions <shorthand>")
	flag.Parse()

	db := config.SetupDatabase(*seed, *seedShort)
	defer db.Close()

	r := mux.NewRouter()

	fmt.Println("hi from v2 :0)")
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", r)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index hit!")
}
