package main

import (
	"net/http"

	"wizzomafizzo/steamgrid-proxy/config"
	"wizzomafizzo/steamgrid-proxy/controller"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cnf := *config.Cnf
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/steamgriddb/api").Subrouter()
	apiRouter.HandleFunc("/search/{gameName}", controller.Search).Methods("GET")
	apiRouter.HandleFunc("/image/{url:.*}", controller.ImageProxy).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})

	http.ListenAndServe(":"+cnf.Port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
