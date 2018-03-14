package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello world from api")
	router := mux.NewRouter()

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/v1/getprofiles", GetProfiles).Methods("GET")
	router.HandleFunc("/v1/login", Login).Methods("POST")
	router.HandleFunc("/v1/update/profile", UpdateProfile).Methods("POST")
	router.HandleFunc("/v1/getprofile/{user}", GetProfileByUserName).Methods("GET")
	router.HandleFunc("/v1/getmatches/{user}", GetProfileMatches).Methods("GET")
	router.HandleFunc("/v1/getusers/interest", GetProfilesWithInterest).Methods("POST")

	log.Fatal(http.ListenAndServe(":8002", router))
}
