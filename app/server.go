package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	addr string
	port int
}

func handleClientConnection(clientConnection net.Conn) {
	// defer clientConnection.Close()

	// Read data
	buf := make([]byte, 1024)
	dataLength, err := clientConnection.Read(buf)
	if err != nil {
		return
	}
	log.Println("Received Data", buf[:dataLength])

	answer := []byte("+PONG\r\n")
	clientConnection.Write(answer)
}

func main() {
	serverConfig := Server{addr: "0.0.0.0", port: 6379}
	server := fmt.Sprintf("%s:%d", serverConfig.addr, serverConfig.port)

	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", server)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on ", listener.Addr())

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		handleClientConnection(connection)
	}
}
