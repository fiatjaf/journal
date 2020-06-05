package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
)

func getMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Characters string              `json:"characters"`
		Methods    map[string][]string `json:"methods"`
	}{
		"$&+-.0123456789:=@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~",
		methodsAvailable,
	})
}

func listEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")

	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		return bucket.ForEach(func(_, v []byte) error {
			w.Write(v)
			w.Write([]byte{'\n'})
			return nil
		})
	})
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{"db-list", err})
		w.WriteHeader(500)
		return
	}
}

func setEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decode
	var entry LogEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"decode", err})
		return
	}

	// time is given in the id
	id := mux.Vars(r)["id"]
	spl := strings.Split(id, "~") // id = <2006-01-02T15:04:05>~<pos-string>
	_, err = time.Parse(DATEFORMAT, spl[0])
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"time", err})
		return
	}
	entry.Time = spl[0]
	entry.Pos = spl[1]

	// validate?
	// ~

	// save
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		v, _ := json.Marshal(entry)
		return bucket.Put([]byte(id), v)
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorResponse{"db-put", err})
		return
	}

	go notifyStateUpdated()

	w.WriteHeader(200)
}

func delEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		return bucket.Delete([]byte(id))
	})
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{"db-delete", err})
		w.WriteHeader(500)
		return
	}

	go notifyStateUpdated()

	w.WriteHeader(200)
}

func queryState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jqfilter string
	if r.Method == "GET" {
		jqfilter, _ = mux.Vars(r)["jq"]
	} else if r.Method == "POST" {
		defer r.Body.Close()
		b, _ := ioutil.ReadAll(r.Body)
		jqfilter = string(b)
	}

	result := state
	if strings.TrimSpace(jqfilter) != "" {
		if result, err = runJQ(r.Context(), state, jqfilter); err != nil {
			json.NewEncoder(w).Encode(ErrorResponse{"jq", err})
			w.WriteHeader(400)
			return
		} else {
			result = state
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}
