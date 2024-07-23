# go-redis

go-redis is a lightweight Redis-like key-value store implementation in Go. It supports various Redis commands and includes a custom RESP (Redis Serialization Protocol) parser and serializer, as well as a TCP server implementation.

## Features

- In-memory key-value store
- RESP (Redis Serialization Protocol) implementation
- TCP server implementation
- Support for various Redis commands:
    - PING
    - ECHO
    - GET
    - SET (with options: NX, XX, EX, PX, EXAT, PXAT)
    - EXISTS
    - DEL
    - INCR
    - DECR
- Thread-safe operations using `sync.Map`

## Project Structure

- `main.go`: TCP server implementation and connection handling
- `pkg/resp/resp.go`: RESP serializer and deserializer implementation
- `pkg/commands/`:
    - `commands.go`: Command handler definitions and main data structure
    - `set.go`: Implementation of the SET command
    - `get.go`: Implementation of the GET command
    - `delete.go`: Implementation of the DEL command
    - `incr.go`: Implementation of the INCR command
    - `decr.go`: Implementation of the DECR command
    - `exists.go`: Implementation of the EXISTS command
    - `echo.go`: Implementation of the ECHO command
    - `ping.go`: Implementation of the PING command

## Running the Server

To run the go-redis server:

1. Ensure you have Go installed on your system.
2. Navigate to the project root directory.
3. Run the following command:

```
go run main.go
```

The server will start and listen on port 6379 (the default Redis port).

## Connecting to the Server

You can connect to the go-redis server using any Redis client. For example, using the `redis-cli`:

```
redis-cli -p 6379
```

## Supported Commands

### PING [message]
Returns PONG if no argument is provided, otherwise returns the message.

### ECHO message
Returns the message sent by the client.

### GET key
Get the value of a key.

### SET key value [NX] [XX] [EX seconds] [PX milliseconds] [EXAT timestamp-seconds] [PXAT timestamp-milliseconds]
Set the value of a key with optional parameters:
- NX: Only set the key if it does not already exist
- XX: Only set the key if it already exists
- EX: Set the specified expire time, in seconds
- PX: Set the specified expire time, in milliseconds
- EXAT: Set the specified Unix time at which the key will expire, in seconds
- PXAT: Set the specified Unix time at which the key will expire, in milliseconds

### EXISTS key [key ...]
Check if one or more keys exist. Returns the number of keys that exist.

### DEL key [key ...]
Delete one or more keys. Returns the number of keys that were removed.

### INCR key
Increment the integer value of a key by one. If the key does not exist, it is set to 0 before performing the operation.

### DECR key
Decrement the integer value of a key by one. If the key does not exist, it is set to 0 before performing the operation.

## Error Handling

The server returns error messages in the following cases:
- Wrong number of arguments for a command
- Invalid integer value for INCR and DECR operations
- Syntax errors in command arguments

## Contributing

Contributions to go-redis are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).