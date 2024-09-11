package server

import (
	"bufio"
	"net"
	"strconv"
	"time"
)

type User struct {
	id         string
	name       string
	connection net.Conn
}

func newUserFromConn(conn net.Conn) *User {
	userId := "user-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	userName := getUserName(conn)
	user := &User{
		id:         userId,
		name:       userName,
		connection: conn,
	}
	return user
}

func getUserName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	choice, _ := reader.ReadString('\n')
	return choice
}
