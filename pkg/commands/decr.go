package commands

import (
	"go-redis/pkg/resp"
	"strconv"
)

func handleDecr(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	record, ok := dataSet.Load(args[0].Bulk)
	if !ok {
		record := Record{Value: "-1"}
		dataSet.Store(args[0].Bulk, record)
		return resp.Value{DataType: resp.TypeInteger, Num: -1}
	} else {
		decrValue, err := strconv.ParseInt(record.(Record).Value, 10, 64)
		if err != nil {
			return resp.Value{DataType: resp.TypeError, Err: errNotInteger}
		}
		decrValue -= 1
		dataSet.Store(args[0].Bulk, Record{Value: strconv.FormatInt(decrValue, 10)})
		return resp.Value{DataType: resp.TypeInteger, Num: int(decrValue)}
	}
}
