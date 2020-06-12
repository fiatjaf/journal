package main

import (
	"encoding/json"
)

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
