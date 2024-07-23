package commands

import (
	"go-redis/pkg/resp"
)

func handleLPush(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}

	key := args[0].Bulk
	elements := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		elements[i-1] = args[i].Bulk
	}

	var list []string
	if value, ok := dataSet.Load(key); ok {
		record := value.(Record)
		if record.Type != TypeList {
			return resp.Value{DataType: resp.TypeError, Err: errWrongType}
		}
		list = record.Value.([]string)
	}

	// Prepend new elements to the list in reverse order
	newList := make([]string, len(elements)+len(list))
	for i := len(elements) - 1; i >= 0; i-- {
		newList[len(elements)-1-i] = elements[i]
	}
	copy(newList[len(elements):], list)

	// Store the updated list
	dataSet.Store(key, Record{Type: TypeList, Value: newList})

	// Return the new length of the list
	return resp.Value{DataType: resp.TypeInteger, Num: len(newList)}
}
