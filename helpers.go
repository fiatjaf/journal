package main

import (
	"context"
	"time"

	"github.com/itchyny/gojq"
)

const (
	DATEFORMAT = "2006-01-02T15:04:05Z"
)

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
