package core

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/tidwall/buntdb"
)

var db *buntdb.DB
var err error
var log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
var state interface{}

const (
	JOURNAL_DB  = "journal.db"
	REDUCE_FILE = "reduce.jq"
	STATE_FILE  = "state.json"
)

func init() {
	db, err = buntdb.Open(JOURNAL_DB)
	if err != nil {
		log.Fatal().Err(err).Str("path", JOURNAL_DB).Msg("couldn't open db")
	}

	// prepare index
	err = db.CreateIndex("datepos", "*",
		buntdb.IndexJSON("date"), buntdb.IndexJSON("pos"))
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't create index")
	}

	// prepare stuff to compute
	prepareComputation()
}
