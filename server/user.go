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
