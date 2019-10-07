package wtfd

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	serverChan  = make(chan chan string, 4)
	messageChan = make(chan string, 1)
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func leaderboardMessageServer(serverChan chan chan string) {
	var clients []chan string
	// And now we listen to new clients and new messages:
	for {
		select {
		case client, _ := <-serverChan:
			clients = append(clients, client)
		case msg, _ := <-messageChan:
			// Send the uptime to all connected clients:
			for _, c := range clients {
				c <- msg
			}
		}
	}
}

func leaderboardServer(serverChan chan chan string) {
	var clients []chan string
	for {
		select {
		case client, _ := <-serverChan:
			clients = append(clients, client)
		}
	}
}

func leaderboardWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	client := make(chan string, 1)
	serverChan <- client // i have no idea what this go magic is

	for {
		select {
		case text, _ := <-client:
			writer, _ := ws.NextWriter(websocket.TextMessage)
			writer.Write([]byte(text))
			writer.Close()
		}
	}

}

func updateScoreboard() error {
	log.Printf("Scoreboard Update\n")
	type userNamePoints struct {
		Name   []string `json:"name"`
		Points []int    `json:"points"`
	}
	var name []string
	var points []int
	allu, err := ormAllUsersSortedByPoints()
	if err != nil {
		log.Printf("Scoreboard Update Error: %v\n", err)
		return err
	}
	for _, u := range allu {
		name = append(name, u.DisplayName)
		points = append(points, u.Points)
	}

	json, err := json.Marshal(&userNamePoints{Name: name, Points: points})
	if err != nil {
		log.Printf("Scoreboard Update Error: %v\n", err)
		return err
	}
	messageChan <- string(json)
        log.Printf("Scoreboard Update String: %s\n", string(json))

	return nil
}
