package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Tile struct {
	Char  rune `json:"char"`
	Value int  `json:"value"`
}
type Player struct {
	Name string `json:"name"`
	Hand []Tile `json:"hand"`
}
type Game struct {
	Players []Player `json:"players"`
	Board   []Tile   `json:"board"`
	Id      int      `json:"id"`
}

//routes handlers
func setupRoutes() {

	//should allow players to play the game modifying the tiles in the tile bag
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello!")
		if err != nil {
			fmt.Println(err)
		}
	})

	//sets up socket connection to allow user to enter a game code and verify the game code
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		log.Println("succesfully joined game ...")
		joinGame(ws)
	})

	//sets up a web socket connection upgrading the http route
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "GET" {
			//create an instance of the game and send json to client
			game := Game{Players: nil, Board: nil, Id: newGameId()}
			jsonGame, err := json.Marshal(game)
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

func newGameId() int {
	// set seed
	rand.Seed(time.Now().UnixNano())
	// generate random number and print on console
	gameCode := rand.Intn(10000000)
	fmt.Println(gameCode)
	return gameCode
}

//join game using message from client
func joinGame(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// when user joins read in the game code and verify it against the current game codes
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
		}
	}
}

//reads the http request to sends back the game code /
func newWsConnection(conn *websocket.Conn) {
	for {
		gameId := 0
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		gameId++
		log.Println(string(p))

		// send the client the generated game code
		if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(newGameId()))); err != nil {
			log.Println(err)
		}
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {

	setupRoutes()
	fmt.Println("Starting server at port at :" + os.Getenv("PORT"))
	srv := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}
	cors.

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
