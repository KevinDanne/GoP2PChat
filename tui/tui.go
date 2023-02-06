package tui

import (
	"bufio"
	"fmt"
	"os"
	"p2pchat/connection/msg"
	"strings"
)

func PrintCommands() {
	fmt.Println("/help Print help message")
	fmt.Println("/connect <IP:PORT> <CHAT_NAME> Create new connection/chat")
	fmt.Println("/msg <CHAT_NAME> <MSG> Print help message")
	fmt.Println("/create-group <GROUP_NAME> <CHAT_NAME...> Print help message")
	fmt.Println("/msg-group <GROUP_NAME> <MSG> Print help message")
	fmt.Println("/broadcast <MSG> Print message to all connections/chats")
}

func Start(username string, chConnections chan msg.ConnectionMsg) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		command := split[0]
		args := split[1:]

		if !strings.HasPrefix(command, "/") {
			fmt.Println("Commands must start with / character")
			continue
		}

		switch command {
		case "/help":
			PrintCommands()
		case "/connect":
			if len(args) < 2 {
				fmt.Println("Missing args. See /help")
				break
			}
			address := args[0]
			chatName := args[1]
			chConnections <- msg.CreateConnection{Address: address, ChatName: chatName}
		case "/msg":
			if len(args) < 2 {
				fmt.Println("Missing args. See /help")
				break
			}
			chatName := args[0]
			message := strings.Join(args[1:], " ")
			chConnections <- msg.Outgoing{Message: message, Sender: username, ChatName: chatName}
		case "/create-group":
			if len(args) < 2 {
				fmt.Println("Missing args. See /help")
				break
			}
			groupName := args[0]
			participants := args[1:]
			chConnections <- msg.CreateGroup{GroupName: groupName, Participants: participants}
		case "/msg-group":
			if len(args) < 2 {
				fmt.Println("Missing args. See /help")
				break
			}
			groupName := args[0]
			message := strings.Join(args[1:], " ")
			chConnections <- msg.OutgoingGroup{GroupName: groupName, Sender: username, Message: message}
		case "/broadcast":
			if len(args) < 1 {
				fmt.Println("Missing args. See /help")
				break
			}
			message := strings.Join(args, " ")
			chConnections <- msg.Broadcast{Sender: username, Message: message}
		}
	}
}
