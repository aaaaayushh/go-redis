package commands

import "go-redis/pkg/resp"

func handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{DataType: resp.TypeString, Str: "PONG"}
	}

	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}
func handleEcho(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeString, Str: "ERR: wrong number of arguments"}
	}
	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}

var CommandHandler = map[string]func([]resp.Value) resp.Value{
	"PING": handlePing,
	"ECHO": handleEcho,
}
