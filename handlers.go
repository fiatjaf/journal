package main

import (
	"encoding/json"
	"net/http"
)

func getMetadata(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Methods map[string][]string `json:"methods"`
	}{methodsAvailable})
}

func newEntry(w http.ResponseWriter, r *http.Request) {
	// db.Update()
}

func setEntry(w http.ResponseWriter, r *http.Request) {

}

func delEntry(w http.ResponseWriter, r *http.Request) {

}
