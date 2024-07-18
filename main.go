package main

import (
	"fmt"
	"log"
	"net"
)

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
		deserializer := NewDeserializer(conn)
		v, err := deserializer.Read()
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(v)

		conn.Write([]byte("+OK\r\n"))
	}
}
