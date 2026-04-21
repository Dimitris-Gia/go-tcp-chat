package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"net-cat/internal/connectionhandling"
	"net-cat/internal/parser"
	"net-cat/internal/server"
)

func main() {
	port, err := parser.GetPortNumber(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Shared state for the whole process — one hub, one history
	hub := server.NewHub()
	history := server.NewHistory()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port :%d\n", port)

	// Accept connections in a loop; each gets its own goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go connectionhandling.HandleConnection(conn, hub, history)
	}
}
