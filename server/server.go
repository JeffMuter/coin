package server

import (
	"bufio"
	"fmt"
	"net"
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
	defer roomServer.mu.Unlock()

	roomServer.rooms[roomName] = room
}

func (roomServer *RoomServer) addUser(user *User) {
	roomServer.mu.Lock()
	defer roomServer.mu.Unlock()

	roomServer.users[user.id] = user
}

func (roomServer *RoomServer) handleNewUserConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("welcome to coin...\n"))

	//create new user
	user := newUserFromConn(conn)

	// add user to roomServer
	roomServer.addUser(user)

	for { // infinite loop until meinMenu returns non nil err
		err := mainMenu(user, roomServer)
		if err != nil {
			fmt.Println("error in main menu: %w,", err)
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

func mainMenu(user *User, roomServer *RoomServer) error {
	defer user.connection.Close()
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

	// now add option switch statement
	switch choice {
	case "0": // create new room.
		user.connection.Write([]byte("enter room name:\n"))
		// get name for room
		roomName, _ := reader.ReadString('\n')
		//create room
		newRoom := newRoom(roomName)
		// add new server to roomserver map
		roomServer.addRoom(roomName, newRoom)
	case "1": // join existing room
		//get choice on room name
		user.connection.Write([]byte("enter room name to join...\n"))
		for {
			roomChoice, _ := reader.ReadString('\n')
			roomChoice = strings.TrimSpace(roomChoice)

			// check to see if room exists
			room, ok := roomServer.rooms[roomChoice]
			if ok {
				// add user to this room
				room.addUser(user)
				// TODO: some logic needs here for what to do to send user to a room.
				// valid input
				break
			} else {
				user.connection.Write([]byte("room not found, try again\n"))
				continue
			}
		}
	case "2": // quit the program
		roomServer.removeUser(user)
		return fmt.Errorf("quit: user chose to quit")
	}
	return nil
}
