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
)

var s Settings
var err error
var log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
var router = mux.NewRouter()
var httpPublic = &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}

type Settings struct {
	Host       string `envconfig:"HOST" default:"0.0.0.0"`
	Port       string `envconfig:"PORT" required:"true"`
	ServiceURL string `envconfig:"SERVICE_URL" required:"true"`
}

func main() {
	err = envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

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
	router.Path("/~/entries").Methods("PATCH").HandlerFunc(batchEntryOps)
	router.Path("/~/entries").Methods("POST").HandlerFunc(newEntry)
	router.Path("/~/entry/{id}").Methods("PUT").HandlerFunc(setEntry)
	router.Path("/~/entry/{id}").Methods("DELETE").HandlerFunc(delEntry)
	router.Path("/~/state").Methods("GET", "POST").HandlerFunc(queryState)
	router.Path("/~/state/{jq}").Methods("GET").HandlerFunc(queryState)
	router.Path("/~~~/state").Methods("GET").HandlerFunc(serveStream)

	// client
	router.PathPrefix("/static/").Methods("GET").Handler(http.FileServer(httpPublic))
	router.PathPrefix("/").Methods("GET").HandlerFunc(serveClient)

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

func serveClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexf, err := httpPublic.Open("static/index.html")
	if err != nil {
		log.Error().Err(err).Str("file", "static/index.html").Msg("make sure you generated bindata.go without -debug")
		return
	}
	fstat, _ := indexf.Stat()
	http.ServeContent(w, r, "static/index.html", fstat.ModTime(), indexf)
	return
}
