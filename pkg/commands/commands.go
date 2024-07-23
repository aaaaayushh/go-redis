package commands

import (
	"go-redis/pkg/resp"
	"sync"
	"time"
)

type DataType int

const (
	TypeString DataType = iota
	TypeList
	TypeSet
	TypeZSet
	TypeHash
)

type Record struct {
	Type       DataType
	Value      interface{}
	ExpiryTime *time.Time
}

var dataSet sync.Map

const (
	errWrongArgsCount = "ERR wrong number of arguments for command"
	errSyntax         = "ERR syntax error"
	errNotInteger     = "ERR value is not an integer or out of range"
	errWrongType      = "WRONGTYPE Operation against a key holding the wrong kind of value"
	okResponse        = "OK"
)

var CommandHandler = map[string]func([]resp.Value) resp.Value{
	"PING":   handlePing,
	"ECHO":   handleEcho,
	"GET":    handleGet,
	"SET":    handleSet,
	"EXISTS": handleExists,
	"DEL":    handleDelete,
	"INCR":   handleIncr,
	"DECR":   handleDecr,
	"LPUSH":  handleLPush,
	"RPUSH":  handleRPush,
	"LRANGE": handleLRange,
}
