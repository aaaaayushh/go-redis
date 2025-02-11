package commands

import "go-redis/pkg/resp"

func handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{DataType: resp.TypeString, Str: "PONG"}
	}

	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}
