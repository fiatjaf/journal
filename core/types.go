package core

import (
	"errors"
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

func (le LogEntry) Validate() error {
	if le.Method == "" {
		return errors.New("no method")
	}

	return nil
}
