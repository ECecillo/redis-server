package main

import (
	"net"
	"testing"
)

func TestServerAnswerToPong(t *testing.T) {
	testServer, err := NewServer("localhost:6370")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	testServer.Start()
	defer testServer.Stop()

	client, err := net.Dial("tcp", testServer.listener.Addr().String())
	if err != nil {
		t.Fatalf("(client) Error connecting to server: %v", err)
	}
	defer client.Close()

	_, err = client.Write([]byte("*1\r\n$4\r\nping\r\n"))
	if err != nil {
		t.Fatalf("(client) Error writing to server: %v", err)
	}

	expected := "+PONG\r\n"
	response := make([]byte, len(expected))
	if _, err := client.Read(response); err != nil {
		t.Fatal(err)
	}
	if string(response) != expected {
		t.Errorf("Expected %q, got %q", expected, string(response))
	}
}
