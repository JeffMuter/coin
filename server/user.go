package server

import (
	"bufio"
	"net"
	"strings"
)

type User struct {
	id         string
	name       string
	connection net.Conn
}

func newUserFromConn(conn net.Conn, name string) *User {
	return &User{
		id:         conn.RemoteAddr().String(),
		name:       name,
		connection: conn,
	}
}

func getUserName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	choice, _ := reader.ReadString('\n')
	cleanChoice := strings.TrimSpace(choice)
	return cleanChoice
}

func makeUserName(conn net.Conn, userMap map[string]*User) (string, error) {

	var name string
loop:
	for { // loop continues until a chose name is unique
		name = getUserName(conn)
		for _, thisUser := range userMap {
			if thisUser.name == name {
				conn.Write([]byte("this user name is currently taken. Your name must be unique, try again...\n"))
			} else {
				break loop
			}
		}
	}
	return name, nil
}
