package main

import (
	"fmt"
	"go-redis/pkg/commands"
	"go-redis/pkg/resp"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	deserializer := resp.NewDeserializer(conn)
	serializer := resp.NewSerializer(conn)

	greetingMsg := "-REDIS 0.0.1 go-redis-server 00000000:0 standalone"
	err := serializer.Write(resp.Value{DataType: resp.TypeString, Str: greetingMsg})
	if err != nil {
		log.Println("Error sending greeting:", err)
		return
	}

	for {
		value, err := deserializer.Read()
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		if value.DataType != resp.TypeArray {
			log.Println("unexpected data type:", value.DataType, " ,expected array")
			continue
		}

		if len(value.Array) == 0 {
			log.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := commands.CommandHandler[command]
		if !ok {
			log.Println("Invalid command:", command)
			//err := serializer.Write(resp.Value{DataType: resp.TypeError, Err: "ERR unknown command"})
			if err != nil {
				log.Println("Error writing response:", err)
			}
			continue
		}

		result := handler(args)
		err = serializer.Write(result)
		if err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

func main() {
	fmt.Println("***********Go-Redis-Server***********")
	// start a server on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Panicln(err)
	}
	defer l.Close()

	fmt.Println("Server is listening on port 6379")
	for {
		fmt.Println("Waiting for a connection...")
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New connection accepted")
		go handleConnection(conn)
	}
}
