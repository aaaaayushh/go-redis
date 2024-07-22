package commands

import (
	"go-redis/pkg/resp"
)

func handleExists(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	var result = 0
	for _, arg := range args {
		key := arg.Bulk
		_, ok := dataSet.Load(key)
		if ok {
			result += 1
		}
	}
	return resp.Value{DataType: resp.TypeInteger, Num: result}
}
