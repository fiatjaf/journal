package main

import "time"

type LogEntry struct {
	Time   time.Time              `json:"time"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}
