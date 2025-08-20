package main

import (
	"github/paolovls/prejects/go_bookstore/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/paolovls/Prejects/go_bookstore/pkg/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
