package commands

import (
	"go-redis/pkg/resp"
	"strconv"
)

func handleIncr(args []resp.Value) resp.Value {
	return handleIncrBy(args, 1)
}

func handleDecr(args []resp.Value) resp.Value {
	return handleIncrBy(args, -1)
}

func handleIncrBy(args []resp.Value, increment int64) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	key := args[0].Bulk
	var value int64
	if record, ok := dataSet.Load(key); ok {
		r := record.(Record)
		if r.Type != TypeString {
			return resp.Value{DataType: resp.TypeError, Err: errWrongType}
		}
		var err error
		value, err = strconv.ParseInt(r.Value.(string), 10, 64)
		if err != nil {
			return resp.Value{DataType: resp.TypeError, Err: errNotInteger}
		}
	}
	value += increment
	dataSet.Store(key, Record{Type: TypeString, Value: strconv.FormatInt(value, 10)})
	return resp.Value{DataType: resp.TypeInteger, Num: int(value)}
}
