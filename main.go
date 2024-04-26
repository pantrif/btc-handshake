package main

import (
	"bitcoin-node-handshake/message"
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	nodeAddress := getEnv("NODE_ADDRESS", "127.0.0.1:18444")
	log.Printf("Attempting to connect to the Bitcoin node at %s...", nodeAddress)

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(context.Background(), "tcp", nodeAddress)
	defer conn.Close()

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}
	log.Println("Connected successfully to the Bitcoin node.")

	msg, err := message.New(message.Version)
	if err != nil {
		log.Fatalf("failed to create message: %v\n", err)
	}

	m, err := msg.Serialize()
	if err != nil {
		log.Fatalf("failed to serialize message: %v\n", err)
	}

	log.Printf("Sending the version message...\n")
	go sendMessage(conn, m)

	var wg sync.WaitGroup
	wg.Add(1)
	go receiveMessages(conn, &wg)
	wg.Wait()
}

func sendMessage(conn net.Conn, msg []byte) {
	_, err := conn.Write(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v\n", err)
	}

	log.Println("Message sent successfully.")
}

func receiveMessages(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	var responseHeader [24]byte

	// First read: Expect a 'version' response
	_, err := conn.Read(responseHeader[:])
	if err != nil {
		log.Fatalf("failed to read version response: %v\n", err)
	}
	command := extractCommand(responseHeader)
	if command == string(message.Version) {
		log.Println("received version message, sending verack.")
		msg, err := message.New(message.Verack)
		if err != nil {
			log.Fatalf("failed to create verach message: %v\n", err)
		}
		m, err := msg.Serialize()
		if err != nil {
			log.Fatalf("failed to serialize verack message: %v\n", err)
		}

		log.Printf("sending the verack message...\n")
		go sendMessage(conn, m)
	} else {
		log.Fatalf("expected version but received: %s\n", command)
	}
	startTime := time.Now()
	timeout := 10 * time.Second

	for {
		if time.Since(startTime) > timeout {
			log.Fatal("timeout waiting for verack message.")
		}

		// Expect a 'verack' after sending our verack command
		_, err = conn.Read(responseHeader[:])
		if err != nil {
			log.Fatalf("failed to read verack response: %v\n", err)
		}
		command = extractCommand(responseHeader)
		if command == string(message.Verack) {
			log.Println("received verack message, handshake completed successfully.")
			break
		}
	}
}

func extractCommand(header [24]byte) string {
	return string(bytes.TrimRight(header[4:16], "\x00"))
}
