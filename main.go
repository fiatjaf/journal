package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"go.etcd.io/bbolt"
)

var s Settings
var db *bbolt.DB
var err error
var log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
var router = mux.NewRouter()
var state interface{}

type Settings struct {
	Host       string `envconfig:"HOST" default:"0.0.0.0"`
	Port       string `envconfig:"PORT" required:"true"`
	ServiceURL string `envconfig:"SERVICE_URL" required:"true"`
}

const (
	JOURNAL_DB   = "journal.db"
	COMPUTE_FILE = "compute.jq"
	STATE_FILE   = "state.json"
)

func main() {
	err = envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

	// journal db
	db, err = bbolt.Open(JOURNAL_DB, 0666, nil)
	if err != nil {
		log.Fatal().Err(err).Str("path", JOURNAL_DB).Msg("couldn't open db")
	}
	defer db.Close()
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("logs"))
		return nil
	})

	// prepare stuff to compute
	prepareComputation()

	// computed state.json
	state, err = computeAll()
	if err != nil {
		es.SendEventMessage(err.Error(), "error", "")
	} else {
		jstate, _ := json.Marshal(state)
		es.SendEventMessage(string(jstate), "state", "")
		err := ioutil.WriteFile(STATE_FILE, jstate, 0666)
		if err != nil {
			log.Error().Err(err).Msg("couldn't write state to file")
		}
	}

	// api routes
	router.Path("/~/metadata").Methods("GET").HandlerFunc(getMetadata)
	router.Path("/~/entries").Methods("GET").HandlerFunc(listEntries)
	router.Path("/~/entry/{id}").Methods("PUT").HandlerFunc(setEntry)
	router.Path("/~/entry/{id}").Methods("DELETE").HandlerFunc(delEntry)
	router.Path("/~/state").Methods("GET", "POST").HandlerFunc(queryState)
	router.Path("/~/state/{jq}").Methods("GET").HandlerFunc(queryState)
	router.Path("/~~~/state").Methods("GET").HandlerFunc(serveStream)

	// js client
	http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "/static/"}))

	// start the server
	log.Info().Str("host", s.Host).Str("port", s.Port).Msg("listening")
	srv := &http.Server{
		Handler:      router,
		Addr:         s.Host + ":" + s.Port,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("error serving http")
	}
}
