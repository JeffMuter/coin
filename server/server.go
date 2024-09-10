package server

import (
	"fmt"
	"net"
	"sync"
)

type TableServer struct {
	tables map[string]*Table
	mu     sync.Mutex
}

func NewTableServer() *TableServer {
	return &TableServer{
		tables: make(map[string]*Table),
	}
}

func (tableServer *TableServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error listening to tcp at address (addr: %s): %w,", addr, err)
	}
	defer listener.Close()
	fmt.Println("server listening on: ", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %w\n", err)
			continue
		}
		go tableServer.handleNewConnection(conn)
	}
}

func (tableServer *TableServer) handleNewConnection(conn net.Conn) {
	// implement future logic here.
	fmt.Println("something has started")
}
