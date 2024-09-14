package main

import (
	"coin/server"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting Chatroom Server...")

	srv := server.NewRoomServer()
	err := srv.Start("localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	// for each room of the srv, we want to continually keep open the listen
}
