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

	_, err = client.Write([]byte("PING\r\n"))
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

func TestServerHandleMutliplePing(t *testing.T) {
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

	_, err = client.Write([]byte("PING\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("(client) Error writing to server: %v", err)
	}

	expected := "+PONG\r\n"
	response1 := make([]byte, len(expected))
	if _, err := client.Read(response1); err != nil {
		t.Fatal(err)
	}
	if string(response1) != expected {
		t.Errorf("Expected %q, got %q", expected, string(response1))
	}
	response2 := make([]byte, len(expected))
	if _, err := client.Read(response2); err != nil {
		t.Fatal(err)
	}
	if string(response2) != expected {
		t.Errorf("Expected %q, got %q", expected, string(response2))
	}
}
