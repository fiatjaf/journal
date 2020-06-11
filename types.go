package main

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type BatchAction struct {
	Id     string    `json:"id"`
	Set    *LogEntry `json:"set"`
	Delete bool      `json:"delete"`
}

type LogEntry struct {
	Time   string                 `json:"time"` // as DATEFORMAT
	Pos    string                 `json:"pos"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func (le *LogEntry) ApplyId(id string) error {
	spl := strings.Split(id, "~") // id = <2006-01-02T15:04:05>~<pos-string>
	_, err = time.Parse(DATEFORMAT, spl[0])
	if err != nil {
		return err
	}
	le.Time = spl[0]
	le.Pos = spl[1]
	return nil
}

func (le LogEntry) Id() string { return le.Time + "~" + le.Pos }

func (le LogEntry) Validate() error {
	if le.Method == "" {
		return errors.New("no method")
	}

	return nil
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
