package main

import (
	"flag"
	"fmt"
	"os"
	"p2pchat/connection"
	"p2pchat/connection/msg"
	"p2pchat/server"
	"p2pchat/tui"
)

func main() {
	port := flag.Uint("port", 7878, "The port number for the tcp server")
	flag.Parse()

	fmt.Print("Enter username: ")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		fmt.Println("Invalid username entered (" + err.Error() + ")")
		os.Exit(1)
	}

	chConnections := make(chan msg.ConnectionMsg)

	go server.Listen(fmt.Sprintf("127.0.0.1:%d", *port), chConnections)
	go connection.HandleAllConnections(chConnections)
	tui.Start(username, chConnections)
}
