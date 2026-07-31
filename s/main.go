package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	fmt.Println("Client connected!")

	for {
		// Read message from the client
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected or error: %v", err)
			break
		}
		
		// Print the received message to the console
		fmt.Printf("Received: %s\n", message)

		// Echo it back to the client (optional but helpful for testing)
		err = ws.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func main() {
	// Configure websocket route
	http.HandleFunc("/", handleConnections)

	port := ":8080"
	fmt.Printf("Server running on port %s, waiting for a connection...\n", port)
	
	// Start the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
