package main

import "encoding/json"

type LogEntry struct {
	Time   string                 `json:"time"` // as DATEFORMAT
	Pos    string                 `json:"pos"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type ErrorResponse struct {
	Type  string
	Error error
}

func (err ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}{err.Type, err.Error.Error()})
}
