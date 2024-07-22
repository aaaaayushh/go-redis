package commands

import (
	"errors"
	"fmt"
	"go-redis/pkg/resp"
	"strconv"
	"time"
)

type SetOptions struct {
	NX   bool
	XX   bool
	EX   int64
	PX   int64
	EXAT int64
	PXAT int64
}

func parseSetOptions(args []resp.Value) (SetOptions, error) {
	opts := SetOptions{
		EX:   -1,
		PX:   -1,
		EXAT: -1,
		PXAT: -1,
	}
	var timeOptionSet = false

	for i := 2; i < len(args); i++ {
		arg := args[i].Bulk
		switch arg {
		case "NX":
			opts.NX = true
		case "XX":
			opts.XX = true
		case "EX", "PX", "EXAT", "PXAT":
			if timeOptionSet {
				return opts, errors.New("only one time-based option (EX, PX, EXAT, PXAT) can be set")
			}
			if i+1 >= len(args) {
				return opts, errors.New(fmt.Sprintf(errWrongArgsCount, "set"))
			}
			value, err := strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil {
				return opts, errors.New(errNotInteger)
			}
			switch arg {
			case "EX":
				opts.EX = value
			case "PX":
				opts.PX = value
			case "EXAT":
				opts.EXAT = value
			case "PXAT":
				opts.PXAT = value
			}
			timeOptionSet = true
			i++
		default:
			return opts, errors.New(errSyntax)
		}
	}
	return opts, nil
}

func handleSet(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	opts, err := parseSetOptions(args)
	if err != nil {
		return resp.Value{DataType: resp.TypeError, Err: err.Error()}
	}

	record := Record{Value: value}

	if opts.EX > 0 {
		expirationTime := time.Now().Add(time.Duration(opts.EX) * time.Second)
		record.ExpiryTime = &expirationTime
	} else if opts.PX > 0 {
		expirationTime := time.Now().Add(time.Duration(opts.PX) * time.Millisecond)
		record.ExpiryTime = &expirationTime
	} else if opts.EXAT > 0 {
		expirationTime := time.Unix(opts.EXAT, 0)
		record.ExpiryTime = &expirationTime
	} else if opts.PXAT > 0 {
		secs := opts.PXAT / 1000
		nsecs := (opts.PXAT % 1000) * 1000000
		expirationTime := time.Unix(secs, nsecs)
		record.ExpiryTime = &expirationTime
	}

	_, exists := dataSet.Load(key)
	if (opts.NX && exists) || (opts.XX && !exists) {
		return resp.Value{DataType: resp.TypeNull, IsNull: true}
	}

	dataSet.Store(key, record)
	return resp.Value{DataType: resp.TypeString, Str: okResponse}
}
