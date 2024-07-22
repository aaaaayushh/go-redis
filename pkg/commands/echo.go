package commands

import "go-redis/pkg/resp"

func handleEcho(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: "ERR: wrong number of arguments"}
	}
	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}
