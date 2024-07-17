package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
	str      string
	num      int
	bulk     string
	array    []Value
}

type Deserializer struct {
	reader *bufio.Reader
}

func NewDeserializer(r io.Reader) *Deserializer {
	return &Deserializer{bufio.NewReader(r)}
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
func (d *Deserializer) readInteger() (num int, numBytes int, err error) {
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

	strLen, _, err := d.readInteger()
	if err != nil {
		return v, err
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
	arrLen, _, err := d.readInteger()
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

func (d *Deserializer) Read() (Value, error) {
	dataType, err := d.reader.ReadByte() // read the first byte of the RESP message to get the datatype
	if err != nil {
		return Value{}, err
	}
	switch dataType {
	//case STRING:
	//	return Value{}, nil
	//case INTEGER:
	//	return Value{}, nil
	//case ERROR:
	//	return Value{}, nil
	case BULK:
		return d.readBulk()
	case ARRAY:
		return d.readArray()
	default:
		return Value{}, errors.New("Invalid data type")
	}
}

// Deserialize
// input will be a RESP string
// output will be a Value struct

func main() {
	fmt.Println("***********Go-Redis-Server***********")
	// start a server on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Panicln(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	for {
		//buf := make([]byte, 1024)
		//_, err := conn.Read(buf)
		//if err != nil {
		//	log.Panicln(err)
		//}
		//fmt.Println(string(buf))
		deserializer := NewDeserializer(conn)
		v, err := deserializer.Read()
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(v)

		conn.Write([]byte("+OK\r\n"))
	}
}
