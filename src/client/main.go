package main

import (
	"fmt"
	"log"
	"net/http"
	"pinguino/src/backend"

	"github.com/gorilla/websocket"
)

type ClientHub struct {
	// player Player
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fs := http.FileServer(http.Dir("../client"))
// }

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// helpful log statement to show connections
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	go handleConnection(ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func handleConnection(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println("main received", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func setupRoutes() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	// framework, player := backend.Setup()
	backend.Setup()
	setupRoutes()
	println("Launching pinguino client... visit http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
