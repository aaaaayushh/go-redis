package commands

import (
	"go-redis/pkg/resp"
	"strconv"
	"sync"
	"time"
)

type Record struct {
	Value      string
	ExpiryTime *time.Time
}

var dataSet = map[string]Record{}
var setMutex = sync.RWMutex{}

func handleSet(args []resp.Value) resp.Value {
	key := args[0].Bulk
	value := args[1].Bulk
	var nx, xx bool
	var ex int64 = -1
	var px int64 = -1
	var exat int64 = -1
	var pxat int64 = -1
	var err error

	for i := 2; i < len(args); i++ {
		arg := args[i].Bulk
		switch arg {
		case "NX":
			nx = true
			break
		case "XX":
			xx = true
			break
		case "EX":
			if i+1 >= len(args) {
				return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'set' command"}
			}
			ex, err = strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil {
				return resp.Value{DataType: resp.TypeError, Err: "Value is not an Integer or is out of range"}
			}
			i++
			break
		case "PX":
			if i+1 >= len(args) {
				return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'set' command"}
			}
			px, err = strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil {
				return resp.Value{DataType: resp.TypeError, Err: "Value is not an Integer or is out of range"}
			}
			i++
			break
		case "EXAT":
			if i+1 >= len(args) {
				return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'set' command"}
			}
			exat, err = strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil {
				return resp.Value{DataType: resp.TypeError, Err: "Value is not an Integer"}
			}
			i++
			break
		case "PXAT":
			if i+1 >= len(args) {
				return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'set' command"}
			}
			pxat, err = strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil {
				return resp.Value{DataType: resp.TypeError, Err: "Value is not an Integer"}
			}
			i++
			break
		default:
			return resp.Value{DataType: resp.TypeError, Err: "ERR syntax error"}
		}
	}
	setMutex.Lock()
	defer setMutex.Unlock()

	exists := false
	if _, ok := dataSet[key]; ok {
		exists = true
	}
	if (nx && exists) || (xx && !exists) {
		return resp.Value{DataType: resp.TypeNull, IsNull: true}
	}
	record := Record{Value: value}

	if ex > 0 {
		expirationTime := time.Now().Add(time.Duration(ex) * time.Second)
		record.ExpiryTime = &expirationTime
	} else if px > 0 {
		expirationTime := time.Now().Add(time.Duration(px) * time.Millisecond)
		record.ExpiryTime = &expirationTime
	} else if exat > 0 {
		expirationTime := time.Unix(exat, 0)
		record.ExpiryTime = &expirationTime
	} else if pxat > 0 {
		secs := pxat / 1000
		nsecs := (pxat % 1000) * 1000000

		expirationTime := time.Unix(secs, nsecs)
		record.ExpiryTime = &expirationTime
	}
	dataSet[key] = record

	return resp.Value{DataType: resp.TypeString, Str: "OK"}
}
func handleGet(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{DataType: resp.TypeError, Err: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].Bulk
	setMutex.RLock()
	record, ok := dataSet[key]
	setMutex.RUnlock()
	if !ok {
		return resp.Value{DataType: resp.TypeNull, IsNull: true}
	} else if record.ExpiryTime != nil {
		if record.ExpiryTime.Before(time.Now()) {
			delete(dataSet, key)
			return resp.Value{DataType: resp.TypeNull, IsNull: true}
		}
	}
	return resp.Value{DataType: resp.TypeString, Str: record.Value}
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
