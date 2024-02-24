package main

const (
	UserConnected SocketMessageKind = iota + 1
	UserDisconnected
	UserCommand
	UserConnectionAcknowledged
)

type SocketMessageKind uint8

type ToolName string

type User struct {
	UserId   UserId `json:"userId"`
	UserName string `json:"userName"`
}

type UserCommandMessage struct {
	X              float64  `json:"x"`
	Y              float64  `json:"y"`
	IsDrag         bool     `json:"isDrag"`
	ToolName       ToolName `json:"toolName"`
	ToolAttributes any      `json:"toolAttributes"`
}

// Define SocketMessage struct
type SocketMessage struct {
	Type    SocketMessageKind  `json:"type"`
	User    User               `json:"user,omitempty"`
	UserID  UserId             `json:"userId,omitempty"`
	Command UserCommandMessage `json:"command,omitempty"`
	Users   []User             `json:"users,omitempty"`
}
