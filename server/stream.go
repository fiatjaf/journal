package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

var es eventsource.EventSource = eventsource.New(
	&eventsource.Settings{
		Timeout:        5 * time.Second,
		CloseOnTimeout: true,
		IdleTimeout:    300 * time.Minute,
	},
	func(r *http.Request) [][]byte {
		return [][]byte{
			[]byte("X-Accel-Buffering: no"),
			[]byte("Cache-Control: no-cache"),
			[]byte("Content-Type: text/event-stream"),
			[]byte("Connection: keep-alive"),
			[]byte("Access-Control-Allow-Origin: *"),
		}
	},
)

var stateUpdated = make(chan bool)

func serveStream(w http.ResponseWriter, r *http.Request) {
	go func() {
		for {
			time.Sleep(25 * time.Second)
			es.SendEventMessage("", "keepalive", "")
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		es.SendRetryMessage(3 * time.Second)
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)

		s, _ := json.Marshal(state)
		es.SendEventMessage(string(s), "state", "")

		for _ = range stateUpdated {
			state, err := computeAll()
			if err != nil {
				log.Warn().Err(err).Msg("error computing state")
				m, _ := json.Marshal(ErrorResponse{"compute", err})
				es.SendEventMessage(string(m), "error", "")
			} else {
				s, _ := json.Marshal(state)
				es.SendEventMessage(string(s), "state", "")
			}
		}
	}()

	es.ServeHTTP(w, r)
}

func notifyStateUpdated() {
	stateUpdated <- true
}
