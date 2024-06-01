package main

import (
	"log"
	"net/http"

	"github.com/Backend/go-jwt-gin-gorm/controllers/authcontroller"
	"github.com/Backend/go-jwt-gin-gorm/controllers/photocontroller"
	"github.com/Backend/go-jwt-gin-gorm/database"
	"github.com/Backend/go-jwt-gin-gorm/middleware"
	"github.com/gorilla/mux"
)

func main() {

	database.ConnectDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/photos", photocontroller.Index).Methods("GET")
	api.HandleFunc("/photos/{id}", photocontroller.Show).Methods("GET")
	api.HandleFunc("/photos", photocontroller.Upload).Methods("POST")
	api.HandleFunc("/photos/{id}", photocontroller.Update).Methods("PUT")
	api.HandleFunc("/photos/{id}", photocontroller.Delete).Methods("DELETE")
	api.Use(middleware.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
