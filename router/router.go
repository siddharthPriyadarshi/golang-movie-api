package router

import (
	"github.com/gorilla/mux"
	"github.com/siddharthPriyadarshi/mongoapi/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.DeleteOneMovie).Methods("DELETE")
	router.HandleFunc("/api/movies", controller.DeleteAllMovie).Methods("DELETE")
	router.HandleFunc("/api/movie/{id}", controller.MarkedAsWatcher).Methods("PUT")

	return router
}
