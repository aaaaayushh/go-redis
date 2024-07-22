package commands

import "go-redis/pkg/resp"

func handleDelete(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	var numKeysDeleted = 0
	for _, arg := range args {
		key := arg.Bulk
		_, ok := dataSet.Load(key)
		if !ok {
			continue
		}
		dataSet.Delete(key)
		numKeysDeleted++
	}
	return resp.Value{DataType: resp.TypeInteger, Num: numKeysDeleted}
}
