package commands

import (
	"go-redis/pkg/resp"
	"strconv"
)

func handleLRange(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{DataType: resp.TypeError, Err: errWrongArgsCount}
	}

	key := args[0].Bulk
	start, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Value{DataType: resp.TypeError, Err: errNotInteger}
	}
	end, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return resp.Value{DataType: resp.TypeError, Err: errNotInteger}
	}

	value, ok := dataSet.Load(key)
	if !ok {
		return resp.Value{DataType: resp.TypeArray, Array: []resp.Value{}}
	}

	record := value.(Record)
	if record.Type != TypeList {
		return resp.Value{DataType: resp.TypeError, Err: errWrongType}
	}

	list := record.Value.([]string)
	listLen := len(list)

	// Handle negative indices
	if start < 0 {
		start = listLen + start
	}
	if end < 0 {
		end = listLen + end
	}

	// Bound check
	if start < 0 {
		start = 0
	}
	if end >= listLen {
		end = listLen - 1
	}

	// If start is greater than end or start is beyond the list, return empty array
	if start > end || start >= listLen {
		return resp.Value{DataType: resp.TypeArray, Array: []resp.Value{}}
	}

	result := make([]resp.Value, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, resp.Value{DataType: resp.TypeBulk, Bulk: list[i]})
	}

	return resp.Value{DataType: resp.TypeArray, Array: result}
}
