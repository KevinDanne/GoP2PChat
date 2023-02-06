package connection

import (
	"fmt"
	"net"
	"p2pchat/connection/msg"
	"p2pchat/server"
)

func HandleAllConnections(chConnections chan msg.ConnectionMsg) {
	connections := make(map[string]*net.TCPConn)
	groups := make(map[string][]string)

	for {
		conMsg, ok := <-chConnections
		if !ok {
			break
		}

		switch payload := conMsg.(type) {
		case msg.CreateConnection:
			tcpAddr, err := net.ResolveTCPAddr("tcp", payload.Address)
			if err != nil {
				fmt.Println("Error while resolving tcp address", err.Error())
				break
			}

			con, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Println("Error while connecting to tcp socket", err.Error())
				break
			}

			go server.Connect(con, payload.ChatName, chConnections)
		case msg.AcceptedConnection:
			connections[payload.ChatName] = payload.Connection
		case msg.CreateGroup:
			groups[payload.GroupName] = payload.Participants
		case msg.Incoming:
			if _, ok := connections[payload.ChatName]; ok {
				fmt.Println(payload.Message)
			} else {
				fmt.Println("No chat found for name " + payload.ChatName)
			}
		case msg.Outgoing:
			con, ok := connections[payload.ChatName]
			if !ok {
				fmt.Println("No chat found with name " + payload.ChatName)
				break
			}
			msg := fmt.Sprintf("[CHAT: %s] %s > %s", payload.ChatName, payload.Sender, payload.Message)
			con.Write([]byte(msg))
			fmt.Println(msg)
		case msg.OutgoingGroup:
			participants, ok := groups[payload.GroupName]
			if !ok {
				fmt.Println("No group found with name " + payload.GroupName)
				break
			}

			sendingMsg := fmt.Sprintf("[GROUP: %s] %s > %s", payload.GroupName, payload.Sender, payload.Message)
			for _, chatName := range participants {
				con, ok := connections[chatName]
				if !ok {
					fmt.Println("No chat found with name " + chatName)
					continue
				}

				con.Write([]byte(sendingMsg))
				fmt.Println(sendingMsg)
			}
		case msg.Broadcast:
			for chatName, con := range connections {
				sendingMsg := fmt.Sprintf("[CHAT: %s] %s > %s", chatName, payload.Sender, payload.Message)
				con.Write([]byte(sendingMsg))
				fmt.Println(sendingMsg)
			}
		}
	}
}
