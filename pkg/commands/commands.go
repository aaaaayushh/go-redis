package commands

import (
	"go-redis/pkg/resp"
	"sync"
	"time"
)

const (
	errWrongArgsCount = "ERR wrong number of arguments for command"
	errSyntax         = "ERR syntax error"
	errNotInteger     = "ERR value is not an integer or out of range"
	okResponse        = "OK"
)

type Record struct {
	Value      string
	ExpiryTime *time.Time
}

var dataSet sync.Map

var CommandHandler = map[string]func([]resp.Value) resp.Value{
	"PING":   handlePing,
	"ECHO":   handleEcho,
	"GET":    handleGet,
	"SET":    handleSet,
	"EXISTS": handleExists,
	"DEL":    handleDelete,
	"INCR":   handleIncr,
}
