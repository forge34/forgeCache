package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("Error listening: ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("Shutting down...")
		listener.Close()
	}()

	log.Println("Server listening on port 4000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Listener closed, stopping accept loop.")
			break
		}
		go handleConnection(conn)
	}

	log.Println("Server exited cleanly")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		fmt.Printf("Received: %s", string(buf[:n]))

	}
}
