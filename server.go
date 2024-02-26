package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type UserId string

type Server struct {
	mu    sync.Mutex
	conns map[UserId]Client
}

func NewServer() *Server {
	return &Server{
		mu:    sync.Mutex{},
		conns: map[UserId]Client{},
	}
}

func (s *Server) addConn(newUser User, ws *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[newUser.UserId] = Client{
		user: User{
			UserId:   newUser.UserId,
			UserName: newUser.UserName,
		},
		conn: ws,
	}
}

func (s *Server) getConn(userId UserId) *websocket.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conns[userId].conn
}

func (s *Server) getConnectedUsers() []User {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := make([]User, len(s.conns))

	i := 0
	for _, client := range s.conns {
		users[i] = client.user
		i++;
	}

	return users
}

func (s *Server) removeConn(userId UserId) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.conns, userId)
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client", ws.RemoteAddr())
	newUser := GenerateUser()
	s.addConn(newUser, ws)
	s.readLoop(newUser.UserId)
}

func (s *Server) readLoop(userId UserId) {
	ws := s.getConn(userId)
	defer ws.Close()
	buf := make([]byte, 1<<10)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]

		switch msg[0] {
		case byte(UserConnected + '0'):
			msg = msg[1:]
			s.handleUserConnection(userId, &msg)
		case byte(UserDisconnected + '0'):
			msg = msg[1:]
			s.removeConn(userId)
		}

		// fmt.Println(string(msg))

		s.broadCast(msg, userId)
	}
}

func (s *Server) handleUserConnection(userId UserId, msg *[]byte) {
	var socketMessage SocketMessage
	err := json.Unmarshal(*msg, &socketMessage)
	if err != nil {
		fmt.Println("Error decoding connection message")
		return
	}

	socketMessage.User = s.conns[userId].user

	*msg = getSocketMessage(socketMessage)

	socketMessage.Type = UserConnectionAcknowledged
	socketMessage.Users = s.getConnectedUsers()

	s.sendMessage(userId, getSocketMessage(socketMessage))
}

func (s *Server) handleUserDisconnection(userId UserId) {
	socketMessage := SocketMessage{
		Type:   UserDisconnected,
		UserID: userId,
	}

	s.removeConn(userId)

	msg := getSocketMessage(socketMessage)
	s.broadCast(msg, userId)
}

// Send message to a single user
func (s *Server) sendMessage(userId UserId, msg []byte) {
	ws := s.getConn(userId)
	ws.Write(msg)
}

// Broadcast message to all the connections
func (s *Server) broadCast(b []byte, senderId UserId) {
	for userId, client := range s.conns {
		if userId == senderId {
			continue
		}
		go func(userId UserId, ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				s.handleUserDisconnection(userId)
				fmt.Println("Write error:", err)
			}
		}(userId, client.conn)
	}
}

// Receives a socketMessage type and returns a message byte buffer after marshalling
func getSocketMessage(socketMessage SocketMessage) (msg []byte) {
	msg, err := json.Marshal(socketMessage)
	if err != nil {
		fmt.Println("Error encoding message to json")
		return
	}
	return
}
