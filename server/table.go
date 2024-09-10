package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Table struct {
	clients   []net.Conn
	messages  chan string
	mu        sync.Mutex
	tableID   string
	maxPeople int
}

// NewTable creates a new table instance.
func NewTable(tableID string) *Table {
	return &Table{
		clients:   []net.Conn{},
		messages:  make(chan string),
		tableID:   tableID,
		maxPeople: 4,
	}
}

// AddClient adds a client to the table.
func (r *Table) AddClient(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients = append(r.clients, conn)
}

// BroadcastMessages sends messages to all clients in the table.
func (r *Table) BroadcastMessages() {
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
func (table *Table) HandleClient(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			table.RemoveClient(conn)
			return fmt.Errorf("error reading message from client (msg: %s): %w,", msg, err)
		}
		if strings.HasPrefix(msg, "/") { // with a / prefix, it's assumed to be a command
			err = table.handleCommand(conn, msg)
			if err != nil {
				return fmt.Errorf("error handling detected commend (cmd: %s): %w,", msg, err)
			}
		} else {
			// broadcast the message, assume it's a normal message
			table.messages <- fmt.Sprintf("%s: %s\n", conn.RemoteAddr(), msg)
		}
	}
}

// RemoveClient removes a client from the table.
func (r *Table) RemoveClient(conn net.Conn) {
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

func (table *Table) handleCommand(conn net.Conn, msg string) error {
	msg = strings.TrimSpace(msg)

	switch msg {
	case: "/new":
		// handle make new room.
	newNewTable()
	}

	return nil
}
