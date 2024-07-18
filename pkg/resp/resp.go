package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	DataType string
	str      string  // simple string value
	num      int     // integer value
	bulk     string  // bulk string value
	err      string  // simple error string value
	array    []Value // array value
	isNull   bool
}

type Deserializer struct {
	reader *bufio.Reader
}
type Serializer struct {
	writer *bufio.Writer
}

func appendCRLF(data []byte) []byte {
	return append(data, '\r', '\n')
}

func NewDeserializer(r io.Reader) *Deserializer {
	return &Deserializer{bufio.NewReader(r)}
}
func NewSerializer(w io.Writer) *Serializer {
	return &Serializer{bufio.NewWriter(w)}
}

// read until we find \r\n
// return line read and number of bytes read
func (d *Deserializer) readLine() (line []byte, numBytes int, err error) {
	for {
		b, err := d.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		numBytes += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], numBytes, nil
}
func (d *Deserializer) readIntegerInLine() (num int, numBytes int, err error) {
	line, numBytes, err := d.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, numBytes, err
	}
	return int(i64), numBytes, nil
}
func (d *Deserializer) readBulk() (Value, error) {
	v := Value{}
	v.DataType = "BULK"

	strLen, _, err := d.readIntegerInLine()
	if err != nil {
		return v, err
	}

	// handle NULL
	if strLen == -1 {
		v.isNull = true
		v.DataType = "NULL"
		return v, nil
	}

	bulkString := make([]byte, strLen)
	_, err = d.reader.Read(bulkString)
	if err != nil {
		return Value{}, err
	}
	v.bulk = string(bulkString)

	// read the trailing CRLF
	_, _, err = d.readLine()
	if err != nil {
		return Value{}, err
	}
	return v, nil
}
func (d *Deserializer) readArray() (Value, error) {
	v := Value{}
	v.DataType = "ARRAY"

	// get the length of the array
	arrLen, _, err := d.readIntegerInLine()
	if err != nil {
		return v, err
	}
	// allocate a slice of arrLen
	v.array = make([]Value, arrLen)
	// read each subsequent entry of the array and insert it into Value[]
	for i := 0; i < arrLen; i++ {
		val, err := d.Read()
		if err != nil {
			return v, err
		}
		v.array[i] = val
	}
	return v, nil
}
func (d *Deserializer) readSimpleString() (Value, error) {
	v := Value{}
	v.DataType = "STRING"
	lineData, _, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.str = string(lineData)
	return v, nil
}
func (d *Deserializer) readInteger() (Value, error) {
	v := Value{}
	v.DataType = "INTEGER"
	line, _, err := d.readLine()
	if err != nil {
		return v, err
	}

	// Convert line to string for easier manipulation
	numStr := string(line)

	// Check if the number is negative
	sign := 1
	if len(numStr) > 0 && numStr[0] == '-' {
		sign = -1
		numStr = numStr[1:] // Remove the sign
	} else if len(numStr) > 0 && numStr[0] == '+' {
		numStr = numStr[1:] // Remove the sign
	}

	// Parse the number
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return v, err
	}

	v.num = sign * num
	return v, nil
}
func (d *Deserializer) readError() (Value, error) {
	v := Value{}
	v.DataType = "ERROR"
	errorMsg, _, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.err = string(errorMsg)
	return v, nil
}
func (d *Deserializer) Read() (Value, error) {
	dataType, err := d.reader.ReadByte() // read the first byte of the RESP message to get the datatype
	if err != nil {
		return Value{}, err
	}
	switch dataType {
	case STRING:
		return d.readSimpleString()
	case INTEGER:
		return d.readInteger()
	case ERROR:
		return d.readError()
	case BULK:
		return d.readBulk()
	case ARRAY:
		return d.readArray()
	default:
		return Value{}, errors.New("invalid data type")
	}
}

func (v Value) Serialize() []byte {
	switch v.DataType {
	case "ARRAY":
		return v.serializeArray()
	case "BULK":
		return v.serializeBulkString()
	case "STRING":
		return v.serializeString()
	case "INTEGER":
		return v.serializeInteger()
	case "ERROR":
		return v.serializeError()
	case "NULL":
		return v.serializeNull()
	default:
		return []byte{}
	}
}
func (v Value) serializeError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.err...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeNull() []byte {
	return []byte("$-1\r\n")
}
func (v Value) serializeString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeArray() []byte {
	var bytes []byte
	length := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = appendCRLF(bytes)
	for _, val := range v.array {
		bytes = append(bytes, val.Serialize()...)
	}
	return bytes
}
func (v Value) serializeBulkString() []byte {
	var bytes []byte
	length := len(v.bulk)
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = appendCRLF(bytes)
	bytes = append(bytes, v.bulk...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	if v.num > 0 {
		bytes = append(bytes, '+')
	}
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = appendCRLF(bytes)
	return bytes
}
