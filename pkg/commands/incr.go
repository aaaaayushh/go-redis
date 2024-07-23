package commands

import (
	"go-redis/pkg/resp"
	"strconv"
)

func handleIncr(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	record, ok := dataSet.Load(args[0].Bulk)
	if !ok {
		record := Record{Value: "0"}
		dataSet.Store(args[0].Bulk, record)
		return resp.Value{DataType: resp.TypeInteger, Num: 0}
	} else {
		incrValue, err := strconv.ParseInt(record.(Record).Value, 10, 64)
		if err != nil {
			return resp.Value{DataType: resp.TypeError, Err: errNotInteger}
		}
		incrValue += 1
		dataSet.Store(args[0].Bulk, Record{Value: strconv.FormatInt(incrValue, 10)})
		return resp.Value{DataType: resp.TypeInteger, Num: int(incrValue)}
	}
}
