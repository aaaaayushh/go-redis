package commands

import (
	"go-redis/pkg/resp"
	"time"
)

func handleGet(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}
	key := args[0].Bulk
	if record, ok := dataSet.Load(key); ok {
		r := record.(Record)
		if r.ExpiryTime != nil && r.ExpiryTime.Before(time.Now()) {
			dataSet.Delete(key)
			return resp.Value{DataType: resp.TypeNull, IsNull: true}
		}

		switch r.Type {
		case TypeString:
			return resp.Value{DataType: resp.TypeBulk, Bulk: r.Value.(string)}
		case TypeList, TypeSet, TypeZSet, TypeHash:
			return resp.Value{DataType: resp.TypeError, Err: errWrongType}
		default:
			return resp.Value{DataType: resp.TypeError, Err: "ERR unknown data type"}
		}
	}
	return resp.Value{DataType: resp.TypeNull, IsNull: true}
}
