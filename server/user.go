package server

import (
	"bufio"
	"net"
)

type User struct {
	id         string
	name       string
	connection net.Conn
}

func newUserFromConn(conn net.Conn) *User {
	return &User{
		id:         conn.RemoteAddr().String(),
		name:       "guest",
		connection: conn,
	}
}

func getUserName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	choice, _ := reader.ReadString('\n')
	return choice
}
