package main

import (
	"context"
	"strings"
	"time"

	"github.com/itchyny/gojq"
)

const (
	DATEFORMAT = "2006-01-02T15:04:05Z"
	CHARACTERS = "$+-.0123456789:=@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

	LEN_CHARACTERS = len(CHARACTERS)
)

func nextPos(pos string) string {
	length := len(pos)
	lastchar := pos[length-1 : length]
	lastcharindex := strings.Index(CHARACTERS, lastchar)
	if lastcharindex == LEN_CHARACTERS-1 {
		return pos + CHARACTERS[0:1]
	} else {
		return pos[0:length-1] + CHARACTERS[lastcharindex+1:lastcharindex+2]
	}
}

func runJQ(
	ctx context.Context,
	object interface{},
	filter string,
) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	query, err := gojq.Parse(filter)
	if err != nil {
		return
	}

	iter := query.RunWithContext(ctx, object)
	v, _ := iter.Next()
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}
