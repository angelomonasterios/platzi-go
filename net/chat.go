package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Client chan<- string

var (
	incomingCLients = make(chan Client)
	levingClients   = make(chan Client)
	messages        = make(chan string)
)

var (
	hostName   = flag.String("h", "localhost", "hostname")
	portNumber = flag.Int("p", 3090, "port")
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	message := make(chan string)
	go messageWrite(conn, message)
	clientName := conn.RemoteAddr().String()

	message <- fmt.Sprintf("welcome to the server, your name %s\n", clientName)

	messages <- fmt.Sprintf("New client is here, name %s\n", clientName)

	incomingCLients <- message

	inputMessage := bufio.NewScanner(conn)
	for inputMessage.Scan() {
		message <- fmt.Sprintf("%s: %s\n", clientName, inputMessage.Text())
	}

	levingClients <- message

	messages <- fmt.Sprintf("%s said goodbye\n", clientName)

}

func messageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		fmt.Fprintln(conn, message)

	}
}

func BroadCast() {
	clients := make(map[Client]bool)
	for {
		select {
		case message := <-messages:
			for client := range clients {
				client <- message
			}
		case newClient := <-incomingCLients:
			clients[newClient] = true
		case levingClient := <-levingClients:
			delete(clients, levingClient)
			close(levingClient)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *hostName, *portNumber))
	if err != nil {
		log.Fatal(err)
	}

	go BroadCast()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		go HandleConnection(conn)

	}

}
