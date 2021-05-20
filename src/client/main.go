package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	pinguino "pinguino/src/backend"

	"github.com/gorilla/websocket"
)

type ClientHub struct {
	framework *pinguino.Framework
	player    *pinguino.Player
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fs := http.FileServer(http.Dir("../client"))
// }

func (ch *ClientHub) wsEndpoint(w http.ResponseWriter, r *http.Request) {
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
	go ch.handleConnection(ws)
}

type Messages struct {
	Control string `json:"control"`
	X       json.RawMessage
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func (ch *ClientHub) handleConnection(conn *websocket.Conn) {
	for {
		// read in a message
		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// // print out that message for clarity
		// fmt.Println("main received", string(p))

		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println(err)
		// 	return
		// }

		var m Messages
		err := conn.ReadJSON(&m)
		if err != nil {
			// handle error
			log.Println(err)
			continue
		}
		switch m.Control {
		case "Move":
			var move pinguino.Move
			if err := json.Unmarshal([]byte(m.X), &move); err != nil {
				// handle error
				log.Println(err)
				return
			}

			ch.player.ClientMovePlayer(move.X, move.Y)

		case "ChatMessage":
			var cm pinguino.ChatMessage
			if err := json.Unmarshal([]byte(m.X), &cm); err != nil {
				// handle error
				log.Println(err)
				return
			}

			ch.player.SendChatMessage(cm.Message)
		}

	}
}

func (ch *ClientHub) setupRoutes() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", ch.wsEndpoint)
}

func main() {
	framework, player := pinguino.Setup()

	port := os.Args[1]

	ch := ClientHub{framework, player}
	ch.setupRoutes()
	println("Launching pinguino client... visit http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
