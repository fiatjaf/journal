package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
)

func getMetadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Characters string              `json:"characters"`
		Methods    map[string][]string `json:"methods"`
	}{
		CHARACTERS,
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

func batchEntryOps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decode
	var actions []BatchAction
	err := json.NewDecoder(r.Body).Decode(&actions)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{"decode", err})
		return
	}

	errorType, err := save(actions)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{errorType, err})
		return
	}

	go notifyStateUpdated()

	w.WriteHeader(200)
}

func newEntry(w http.ResponseWriter, r *http.Request) {
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

	// pos will be "", we will set it so this entry goes to the last pos in this date
	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		c := bucket.Cursor()
		k, _ := c.Seek([]byte(entry.Time))
		if k == nil || string(k) != entry.Time {
			// first at this date
			entry.Pos = CHARACTERS[0:1]
		} else {
			for {
				nextk, _ := c.Next()
				if nextk != nil || strings.Split(string(nextk), "~")[0] != entry.Time {
					// we reached the next date after the one we're in
					pos := strings.Split(string(k), "~")[1]
					entry.Pos = nextPos(pos)
				}
				k = nextk
			}
		}

		return nil
	})

	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{"db-list", err})
		w.WriteHeader(500)
		return
	}

	// save
	errorType, err := save([]BatchAction{
		{
			Id:  entry.Id(),
			Set: &entry,
		},
	})
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{errorType, err})
		return
	}

	go notifyStateUpdated()

	w.WriteHeader(200)
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

	errorType, err := save([]BatchAction{
		{
			Id:  mux.Vars(r)["id"],
			Set: &entry,
		},
	})
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{errorType, err})
		return
	}

	go notifyStateUpdated()

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(entry)
}

func delEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorType, err := save([]BatchAction{
		{
			Id:     mux.Vars(r)["id"],
			Delete: true,
		},
	})
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{errorType, err})
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
