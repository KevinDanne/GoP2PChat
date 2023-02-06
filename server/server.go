package server

import (
	"fmt"
	"net"
	"p2pchat/connection/msg"
)

func Listen(addr string, chConnections chan msg.ConnectionMsg) {
	fmt.Printf("Listening on address: %s\n", addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("Error while resolving tcp address", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error while listening on tcp socket", err)
		return
	}

	for {
		con, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Error while accepting tcp connection", err)
			continue
		}
		if chatName, err := Validate(con); err == nil {
			go Connect(con, chatName, chConnections)
		}
	}
}

func Validate(con *net.TCPConn) (string, error) {
	fmt.Println("Validating incoming connection from " + con.RemoteAddr().String())

	// TODO enter validation logic
	fmt.Println("Connection validated")

	return con.RemoteAddr().String(), nil
}

func Connect(con *net.TCPConn, chatName string, chConnections chan msg.ConnectionMsg) {
	fmt.Println("Connecting to " + con.RemoteAddr().String())

	chConnections <- msg.AcceptedConnection{Connection: con, ChatName: chatName}

	for {
		incomingBuffer := make([]byte, 1024)
		_, err := con.Read(incomingBuffer)
		if err != nil {
			fmt.Println("Error while reading from stream. Closing chat " + chatName)
			break
		}

		chConnections <- msg.Incoming{ChatName: chatName, Message: string(incomingBuffer)}
	}
}
