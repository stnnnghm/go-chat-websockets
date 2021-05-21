package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// upgrader is an object with methods for taking a normal HTTP conn
// and upgrading it to a WebSocket
var upgrader = websocket.Upgrader{}

// Message holds our messages
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"Message"`
}

func main() {
	// Create simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Listen for incoming chat messages
	go handleMessages()

	// Start the server and log any errors
	log.Println("initializing server :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection
	defer ws.Close()

	// Register new client
	clients[ws] = true

	for {
		var msg Message

		// Read new JSON message and map it to a Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Send message to broadcast chan
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab next msg from broadcast chan
		msg := <-broadcast

		// Send msg to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				// Log the error, close and remove the client from the clients map
				log.Printf("error: %v", err)
				delete(clients, client)
				client.Close()
			}
		}
	}
}
