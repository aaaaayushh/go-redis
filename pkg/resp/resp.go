package resp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type DataType int

const (
	TypeString DataType = iota
	TypeError
	TypeInteger
	TypeBulk
	TypeArray
	TypeNull
)

type Value struct {
	DataType DataType
	Str      string  // simple string value
	Num      int     // integer value
	Bulk     string  // bulk string value
	Err      string  // simple error string value
	Array    []Value // array value
	IsNull   bool
}

type Deserializer struct {
	reader *bufio.Reader
}
type Serializer struct {
	writer io.Writer
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
	v.DataType = TypeBulk

	strLen, _, err := d.readIntegerInLine()
	if err != nil {
		return v, err
	}

	// handle NULL
	if strLen == -1 {
		v.IsNull = true
		v.DataType = TypeNull
		return v, nil
	}

	bulkString := make([]byte, strLen)
	_, err = d.reader.Read(bulkString)
	if err != nil {
		return Value{}, err
	}
	v.Bulk = string(bulkString)

	// read the trailing CRLF
	_, _, err = d.readLine()
	if err != nil {
		return Value{}, err
	}
	return v, nil
}
func (d *Deserializer) readArray() (Value, error) {
	v := Value{}
	v.DataType = TypeArray

	// get the length of the array
	arrLen, _, err := d.readIntegerInLine()
	if err != nil {
		return v, err
	}

	// handle NULL
	if arrLen == -1 {
		v.IsNull = true
		v.DataType = TypeNull
		return v, nil
	}
	// allocate a slice of arrLen
	v.Array = make([]Value, arrLen)
	// read each subsequent entry of the array and insert it into Value[]
	for i := 0; i < arrLen; i++ {
		val, err := d.Read()
		if err != nil {
			return v, err
		}
		v.Array[i] = val
	}
	return v, nil
}
func (d *Deserializer) readSimpleString() (Value, error) {
	v := Value{}
	v.DataType = TypeString
	lineData, _, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.Str = string(lineData)
	return v, nil
}
func (d *Deserializer) readInteger() (Value, error) {
	v := Value{}
	v.DataType = TypeInteger
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

	v.Num = sign * num
	return v, nil
}
func (d *Deserializer) readError() (Value, error) {
	v := Value{}
	v.DataType = TypeError
	errorMsg, _, err := d.readLine()
	if err != nil {
		return v, err
	}
	v.Err = string(errorMsg)
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
	case TypeArray:
		return v.serializeArray()
	case TypeBulk:
		return v.serializeBulkString()
	case TypeString:
		return v.serializeString()
	case TypeInteger:
		return v.serializeInteger()
	case TypeError:
		return v.serializeError()
	case TypeNull:
		return v.serializeNull()
	default:
		return []byte{}
	}
}
func (v Value) serializeError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Err...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeNull() []byte {
	return []byte("$-1\r\n")
}
func (v Value) serializeString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeArray() []byte {
	var bytes []byte
	length := len(v.Array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = appendCRLF(bytes)
	for _, val := range v.Array {
		bytes = append(bytes, val.Serialize()...)
	}
	return bytes
}
func (v Value) serializeBulkString() []byte {
	var bytes []byte
	length := len(v.Bulk)
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = appendCRLF(bytes)
	bytes = append(bytes, v.Bulk...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (v Value) serializeInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	if v.Num > 0 {
		bytes = append(bytes, '+')
	}
	bytes = append(bytes, strconv.Itoa(v.Num)...)
	bytes = appendCRLF(bytes)
	return bytes
}
func (s *Serializer) Write(v Value) error {
	var bytes = v.Serialize()
	_, err := s.writer.Write(bytes)
	if f, ok := s.writer.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
