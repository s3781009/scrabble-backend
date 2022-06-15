package network

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"scrabble-backend/game"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// SetupRoutes routes handlers
func SetupRoutes() {
	var games []game.Game
	//sets up socket connection to allow user to enter a game code and verify the game code
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		wsHandler(ws, &games)
	})

	//sets up a web socket connection upgrading the http route
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "GET" {
			//create an instance of the newGame and send json to client
			newGame := game.NewGame()
			games = append(games, newGame)
			jsonGame, err := json.Marshal(newGame)
			if err != nil {
				log.Println(err)
			}
			_, err = w.Write(jsonGame)
			if err != nil {
				log.Println(err)
			}
		}
	})
}
