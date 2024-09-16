package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type RoomServer struct {
	rooms map[string]*Room
	users map[string]*User
	mu    sync.Mutex
}

func NewRoomServer() *RoomServer {
	return &RoomServer{
		rooms: make(map[string]*Room),
		users: make(map[string]*User),
	}
}

func (roomServer *RoomServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error listening to tcp at address (addr: %s): %w,", addr, err)
	}
	defer listener.Close()
	fmt.Println("server listening on: ", addr)

	for { // run loop, always listening for a user to connect
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err)
			continue
		}
		go roomServer.handleNewUserConnection(conn)
	}
}

func (roomServer *RoomServer) addRoom(roomName string, room *Room) {
	roomServer.mu.Lock()
	roomServer.rooms[roomName] = room
	roomServer.mu.Unlock()
}

func (roomServer *RoomServer) addUser(user *User) {
	roomServer.mu.Lock()
	roomServer.users[user.id] = user
	roomServer.mu.Unlock()
}

func (roomServer *RoomServer) handleNewUserConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("welcome to coin...\ncreate a username:"))

	//create new user
	name := getUserName(conn)
	user := newUserFromConn(conn, name)

	// add user to roomServer
	roomServer.addUser(user)

	for { // infinite loop until meinMenu returns non nil err
		room, err := mainMenu(user, roomServer)
		if err != nil {
			fmt.Println("error in main menu: %w,", err)
			break
		} else if room != nil {
			fmt.Println("userCount: " + strconv.Itoa(len(room.users)))
			handleRoom(user, room)
			break
		}
	}
}

// removes room from roomserver
func (roomServer *RoomServer) removeRoom(room *Room) {
	roomServer.mu.Lock()
	defer roomServer.mu.Unlock()

	delete(roomServer.rooms, room.name)
}

// Remove user from from server
func (roomServer *RoomServer) removeUser(user *User) {
	roomServer.mu.Lock()
	defer roomServer.mu.Unlock()

	delete(roomServer.users, user.id)
	user.connection.Close()
}

func mainMenu(user *User, roomServer *RoomServer) (*Room, error) {
	fmt.Println("start main menu")
	reader := bufio.NewReader(user.connection)

	// list out the rooms
	if len(roomServer.rooms) == 0 {
		user.connection.Write([]byte("currently no rooms active...\n"))
	} else {
		user.connection.Write([]byte("Available rooms:\n"))
		for _, room := range roomServer.rooms {
			user.connection.Write([]byte(room.name + "\n"))
		}
	}

	user.connection.Write([]byte("choose an option from the menu:\n"))
	user.connection.Write([]byte("0: create a room\n1: join a room\n2: end program\n"))

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var room *Room

	// now add option switch statement
	switch choice {
	case "0": // create new room.
		user.connection.Write([]byte("enter room name:\n"))
		// get name for room
		roomName, _ := reader.ReadString('\n')
		roomName = strings.TrimSpace(roomName)
		//create room
		room = newRoom(roomName)
		room.addUser(user)
		// add new server to roomserver map
		roomServer.addRoom(roomName, room)
		fmt.Printf("name: %s\nmessages: %v,\nusers: %v\nconn: %v\n", room.name, room.messages, room.users, user.connection)

		return room, nil

	case "1": // join existing room
		//get choice on room name
		user.connection.Write([]byte("enter room name to join...\n"))
		for {
			roomChoice, _ := reader.ReadString('\n')
			roomChoice = strings.TrimSpace(roomChoice)

			// check to see if room exists
			roomServer.mu.Lock()
			room, ok := roomServer.rooms[roomChoice]
			roomServer.mu.Unlock()
			if ok {
				room.addUser(user) // add user to this room
				return room, nil
			} else {
				user.connection.Write([]byte("room not found, try again\n"))
				continue
			}
		}
	case "2": // quit the program
		roomServer.removeUser(user)
		user.connection.Close()
		return room, fmt.Errorf("quit: user chose to quit")
	}

	fmt.Println("end of mainMen")

	return room, nil
}

func handleRoom(user *User, room *Room) {
	fmt.Printf("count room messages: %d\n", len(room.messages))

	// Create a separate goroutine to print out new messages as they are received
	go func() {
		for msg := range room.messages {
			// Broadcast the message to all users in the room
			room.mu.Lock() // Ensure thread-safe access to the users list
			for _, u := range room.users {
				_, err := u.connection.Write([]byte(msg + "\n"))
				if err != nil {
					fmt.Printf("Error sending message to user %s: %v\n", u.name, err)
				}
			}
			room.mu.Unlock()
		}
	}()

	// Now, listen for user's input and send it to the room's messages channel
	reader := bufio.NewReader(user.connection)
	for {
		// Read message from user input
		userMsg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading user message:", err)
			break
		}

		// Clean up the message
		userMsg = strings.TrimSpace(userMsg)
		if userMsg == "" {
			continue // Skip if the user sends an empty message
		}

		// Send the message to the room's messages channel
		fullMsg := fmt.Sprintf("%s: %s", user.name, userMsg)
		fmt.Println("adding fullmsg: " + fullMsg)
		room.messages <- fullMsg
	}
}
