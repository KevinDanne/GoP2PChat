package msg

import "net"

type ConnectionMsg interface{}

type CreateConnection struct {
	Address  string
	ChatName string
}

type AcceptedConnection struct {
	Connection *net.TCPConn
	ChatName   string
}

type CreateGroup struct {
	GroupName    string
	Participants []string
}

type Incoming struct {
	ChatName string
	Message  string
}

type Outgoing struct {
	Message  string
	Sender   string
	ChatName string
}

type OutgoingGroup struct {
	Message   string
	Sender    string
	GroupName string
}

type Broadcast struct {
	Message string
	Sender  string
}
