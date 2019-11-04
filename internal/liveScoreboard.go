package wtfd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hash/crc32"
	"log"
	"net/http"
	"time"
)

type tableData struct {
	Names  []string `json:"name"`
	Points []int    `json:"points"`
}

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
	updateScoreboard()

	for {
		select {
		case text, _ := <-client:
			writer, _ := ws.NextWriter(websocket.TextMessage)
			writer.Write([]byte(text))
			writer.Close()
		}
	}

}

func generateTableData() (tableData, error) {
	allu, err := wtfdDB.AllUsersSortedByPoints()
	if err != nil {
		return tableData{}, err
	}
	var name []string
	var points []int

	for _, u := range allu {
		name = append(name, u.DisplayName)
		points = append(points, u.Points)
	}
	return tableData{Names: name, Points: points}, nil
}

func updateScoreboard() error {
	type chartDataPoint struct {
		T     string `json:"t"`
		Label string `json:"tooltipLabel"`
		Y     int    `json:"y"`
	}
	type chartData struct {
		Label  string           `json:"label"`
		Data   []chartDataPoint `json:"data"`
		Color  string           `json:"backgroundColor"`
		Pcolor string           `json:"borderColor"`
	}
	type leaderboardData struct {
		TableData tableData   `json:"table"`
		ChartData []chartData `json:"chart"`
	}

	log.Printf("Scoreboard Update\n")
	users, err := wtfdDB.AllUsersSortedByPoints()
	var datas []chartData
	for _, u := range users {
		solves := wtfdDB.GetSolvesWithTime(u.Name)
		data := make([]chartDataPoint, len(solves)+1)
		sum := 0
		data[0] = chartDataPoint{T: u.Created.Format(time.RubyDate), Y: sum, Label: "User Created"}
		for i, s := range solves {
			chall, err := challs.Find(s.ChallengeName)
			if err != nil {
				log.Printf("Scoreboard Update Error: %v, %v\n", err, s.ChallengeName)
				return err
			}
			sum += chall.Points
			data[i+1] = chartDataPoint{T: s.Created.Format(time.RubyDate), Y: sum, Label: s.ChallengeName}
		}
		a := fmt.Sprintf("#%X", crc32.ChecksumIEEE([]byte(u.DisplayName)))[0:7]
		datas = append(datas, chartData{Pcolor: a, Color: a, Label: u.DisplayName, Data: data})

	}

	td, err := generateTableData()
	if err != nil {
		log.Printf("Scoreboard Update Error: %v\n", err)
		return err
	}
	ld := leaderboardData{TableData: td, ChartData: datas}
	jsona, err := json.Marshal(&ld)
	if err != nil {
		log.Printf("Scoreboard Update Error: %v\n", err)
		return err
	}
	messageChan <- string(jsona)

	return nil
}
