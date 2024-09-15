package server

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type Room struct {
	users     map[string]*User
	messages  chan string
	mu        sync.Mutex
	name      string
	maxPeople int
}

// Newroom creates a new room instance.
func newRoom(roomName string) *Room {
	return &Room{
		users:     make(map[string]*User),
		messages:  make(chan string, 100),
		name:      roomName,
		maxPeople: 4,
	}
}

// AddClient adds a client to the room.
func (room *Room) addUser(user *User) {
	room.mu.Lock()
	room.users[user.id] = user
	room.mu.Unlock()
}

// broadcastMessages prints the room's existing messages to the user's screen
func (room *Room) broadcastMessages() {
	fmt.Println("bcCurrentMsgsSta")

	for msg := range room.messages {
		room.mu.Lock()
		for _, user := range room.users {
			_, err := user.connection.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("error broadcasting to user: ", user.name)
				room.removeClient(user)
				return
			}
		}
		room.mu.Unlock()
	}
	fmt.Println("bcCurrentMsgsEnd")
}

// permanently listening, and adds new messages when a user sends one
func (room *Room) handleNewMessage(user *User) {
	reader := bufio.NewReader(user.connection)
	fmt.Println("input received in handleNewMsg")

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from user input: ", err)
			room.removeClient(user)
			return
		}
		msg = strings.TrimSpace(msg)

		if strings.HasPrefix(msg, "/") {
			// add a command handler here for commands starting with "/"
			continue
		}
		room.messages <- user.name + " " + msg
	}
}

// RemoveClient removes a client from the room.
func (room *Room) removeClient(user *User) {
	room.mu.Lock()
	defer room.mu.Unlock()

	delete(room.users, user.id)
	user.connection.Close()

	if len(room.users) == 0 {
		close(room.messages)
		// TODO: need to handle removing room from the server
	}
}

func listRooms(roomServer *RoomServer) ([]string, error) {
	var roomList []string

	for _, room := range roomServer.rooms {
		roomList = append(roomList, room.name)
	}
	return roomList, nil
}
