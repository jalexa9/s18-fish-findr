package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Clemson-CPSC-4910/s18-fish-findr/go/fisher"
	"github.com/Clemson-CPSC-4910/s18-fish-findr/go/postgres"
	"github.com/gorilla/mux"
)

// GetIndex returns the index page for the webapp.
func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data, err := ioutil.ReadFile("/webapp/index.html") // Open the file.
	if err != nil {
		log.Printf("There was an error opening the index.html file.")
	}
	_, _ = w.Write(data)
}

// Login will take a json packet with login info and login user or deney login.
// JSON should appear like below for this:
/*
	{
		"user_name": "user1ex",
		"password": "passwordex"
	}
*/
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userSubmit fisherLogin
	err := json.NewDecoder(r.Body).Decode(&userSubmit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error decoding submit from json. %v\n", err)
		return
	}

	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	_, err = db.GetIfLogin(userSubmit.UserName, userSubmit.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting from database. %v\n", err)
		return
	}
	err = db.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error closing database. %v\n", err)
		return
	}

	// TODO, find a way to set session on server and client, maybe use whats below
	/*
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "login cookie", Value: profile.UserName, Expires: expiration}
		http.SetCookie(w, &cookie)
		fmt.Printf("%v\n", cookie)
	*/
	w.WriteHeader(http.StatusOK)
}

type fisherLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// GetProfiles will return the profiles in json format.
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	profiles, err := db.GetProfiles()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting from database. %v\n", err)
		return
	}
	err = db.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error closing database. %v\n", err)
		return
	}

	respondJSON, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error marshalling json. %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respondJSON)
}

// GetProfileMatches returns the matches in order for the passed profile.
func GetProfileMatches(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	profiles, err := db.GetProfileMatches(vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting from database. %v\n", err)
		return
	}

	err = db.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error closing database. %v\n", err)
		return
	}

	respondJSON, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error marshalling json. %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respondJSON)
}

// GetProfileByUserName returns a profile it gets by the user name.
func GetProfileByUserName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	profile, err := db.GetProfileByUserName(vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting from database. %v\n", err)
		return
	}
	err = db.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error closing database. %v\n", err)
		return
	}

	respondJSON, err := json.Marshal(profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error marshalling json. %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respondJSON)
}

// UpdateProfile will accept a profile in json and then update it.
// JSON should come in like below for this:
/*
	{
		"first_name": "firstex",
		"last_name": "lastex",
		"user_name": "userex",
		"password": "passwordex",
		"phone_number": "123421",
		"email_address": "example@gmail.com",
		"facebook_profile": "www.facebook.com",
		"bio": "example bio",
		"interest": [
			{
				"interest_type": "typeex1"
			},
			{
				"interest_type": "typeex2"
			} ...
		]
	}
*/
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userSubmit fisher.Profile
	err := json.NewDecoder(r.Body).Decode(&userSubmit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error decoding submit from json. %v\n", err)
		return
	}

	//fmt.Printf("%v\n", userSubmit)
	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	err = db.UpdateProfile(userSubmit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error updating database. %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetProfilesWithInterest returns a list of profiles with the passed interest.
// They will be sorted by most matching intrest to least matching interest.
// JSON should come in for this like below:
/*
	{
		"user_name": "user name of requestor"
		"interest": [
			{
				"interest_type": "typeex1"
			},
			{
				"interest_type": "typeex2"
			} ...
		]
	}
*/
func GetProfilesWithInterest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userSubmit findUserInterest
	err := json.NewDecoder(r.Body).Decode(&userSubmit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error decoding submit from json. %v\n", err)
		return
	}

	db, err := postgres.CreateDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error starting database. %v\n", err)
		return
	}
	profiles, err := db.GetProfilesWithInterest(userSubmit.UserName, userSubmit.Interest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting from database. %v\n", err)
		return
	}

	err = db.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error closing database. %v\n", err)
		return
	}

	respondJSON, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error marshalling json. %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respondJSON)
}

type findUserInterest struct {
	UserName string            `json:"user_name"`
	Interest []fisher.Interest `json:"interest"`
}
