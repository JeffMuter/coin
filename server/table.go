package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Room struct {
	clients   []net.Conn
	messages  chan string
	mu        sync.Mutex
	roomId    string
	maxPeople int
}

// Newroom creates a new room instance.
func NewRoom(roomId string) *Room {
	return &Room{
		clients:   []net.Conn{},
		messages:  make(chan string),
		roomId:    roomId,
		maxPeople: 4,
	}
}

// AddClient adds a client to the room.
func (r *Room) AddClient(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients = append(r.clients, conn)
}

// BroadcastMessages sends messages to all clients in the room.
func (r *Room) BroadcastMessages() {
	for msg := range r.messages {
		r.mu.Lock()
		for _, client := range r.clients {
			_, err := client.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
		r.mu.Unlock()
	}
}

// HandleClient reads messages from a client and broadcasts them.
func (room *Room) HandleClientMessages(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			room.RemoveClient(conn)
			return fmt.Errorf("error reading message from client (msg: %s): %w,", msg, err)
		}
		if strings.HasPrefix(msg, "/") { // with a / prefix, it's assumed to be a command
			err = room.handleCommand(conn, msg)
			if err != nil {
				return fmt.Errorf("error handling detected commend (cmd: %s): %w,", msg, err)
			}
		} else {
			// broadcast the message, assume it's a normal message
			room.messages <- fmt.Sprintf("%s: %s\n", conn.RemoteAddr(), msg)
		}
	}
}

// RemoveClient removes a client from the room.
func (r *Room) RemoveClient(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, client := range r.clients {
		if client == conn {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			break
		}
	}
	conn.Close()

	if len(r.clients) == 0 {
		close(r.messages)
	}
}

func (room *Room) handleCommand(conn net.Conn, msg string) error {
	msg = strings.TrimSpace(msg)

	switch msg {
	case "/new":
		// handle make new room.
		NewRoom()
		return nil
	}

	return nil
}

func listRooms() ([]string, error) {
	var roomList []string

	for _, room := range RoomServer.rooms {
		roomList = append(roomList, room.roomId)
	}

	return nil, roomList
}
