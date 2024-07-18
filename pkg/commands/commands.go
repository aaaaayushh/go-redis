package commands

import (
	"go-redis/pkg/resp"
	"sync"
)

var dataSet = map[string]string{}
var setMutex = sync.RWMutex{}

func handleSet(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].Bulk
	value := args[1].Bulk
	setMutex.Lock()
	dataSet[key] = value
	setMutex.Unlock()
	return resp.Value{DataType: resp.TypeString, Str: "OK"}
}
func handleGet(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].Bulk
	setMutex.RLock()
	value, ok := dataSet[key]
	setMutex.RUnlock()
	if !ok {
		return resp.Value{DataType: resp.TypeNull, IsNull: true}
	}
	return resp.Value{DataType: resp.TypeString, Str: value}
}

func handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{DataType: resp.TypeString, Str: "PONG"}
	}

	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}
func handleEcho(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: "ERR: wrong number of arguments"}
	}
	return resp.Value{DataType: resp.TypeString, Str: args[0].Bulk}
}

var CommandHandler = map[string]func([]resp.Value) resp.Value{
	"PING": handlePing,
	"ECHO": handleEcho,
	"GET":  handleGet,
	"SET":  handleSet,
}
