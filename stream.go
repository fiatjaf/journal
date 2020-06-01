package main

import (
	"net/http"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

func stateStream() eventsource.EventSource {
	es := eventsource.New(
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
	go func() {
		for {
			time.Sleep(25 * time.Second)
			es.SendEventMessage("", "keepalive", "")
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		es.SendRetryMessage(3 * time.Second)
	}()

	return es
}
