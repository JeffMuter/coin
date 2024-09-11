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
}
