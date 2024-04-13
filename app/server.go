package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
}

func NewServer(serverAddress string) (*Server, error) {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", serverAddress, err)
	}

	return &Server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if netErr, ok := err.(*net.OpError); ok && !netErr.Temporary() {
				log.Println("Listener closed, shutting down accept routine.")
				return
			}
			log.Println("Error accepting connection:", err)
			continue
		}
		s.connection <- conn
	}
}

func (s *Server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			// In case we get a new connection run a go routine to handle it
			go s.handleClientConnection(conn)
		}
	}
}

func (s *Server) Start() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()
}

func (s *Server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}
}

func (s *Server) handleClientConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		dataLength, err := conn.Read(buf)
		if err == io.EOF {
			break // Client closed connection
		}
		if err != nil {
			log.Println("Error reading:", err)
			return
		}
		request := string(buf[:dataLength])
		fmt.Println(request)
		// if strings.TrimSpace(request) == "*1\r\n$4\r\nping\r\n" {
		conn.Write([]byte("+PONG\r\n"))
		// }
	}
}

func main() {
	serverAddress := "0.0.0.0:6379"
	server, err := NewServer(serverAddress)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	server.Start()
	fmt.Println("Running server on ", serverAddress)

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	server.Stop()
	fmt.Println("Server stopped.")
}
