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
	mu    sync.Mutex
}

func NewRoomServer() *RoomServer {
	return &RoomServer{
		rooms: make(map[string]*Room),
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
			fmt.Printf("error accepting connection: %w\n", err)
			continue
		}
		go roomServer.handleNewClientConnection(conn)
	}
}

func (roomServer *RoomServer) handleNewClientConnection(conn net.Conn) {
	defer conn.Close()
	// implement future logic here.
	reader := bufio.NewReader(conn)

	conn.Write([]byte("welcome to coin...\n"))

	for {
		err := mainMenu(conn, reader)
		if err != nil {
			fmt.Println("error in main menu: %w,", err)
			return
		}
	}

}

func mainMenu(conn net.Conn, reader *bufio.Reader) error {
	defer conn.Close()
	conn.Write([]byte("choose an option from the menu:\n"))
	conn.Write([]byte("0: List all rooms\n1: create a room\n2: join a room\n3: end program\n"))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	// now add option switch statement

	switch choice {
	case "0":
		roomList := listRooms
		if len(roomList) == 0 {
			conn.Write([]byte("currently no rooms active...\n"))
		} else {
			conn.Write([]byte("Available rooms:\n"))
			for _, room := range roomList {
				conn.Write([]byte(room.RoomId + "\n"))
			}
		}
	case "1":
		conn.Write([]byte("enter room name:\n"))
		// get name for room
		roomName, _ := reader.ReadString('\n')
		//create room
		newRoom := NewRoom(roomName)
		RoomServer.rooms = [roomName]newRoom

	//add room to server map
	case "2":
	case "3":

	}
	return nil
}
