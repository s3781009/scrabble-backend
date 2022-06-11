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
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Tile struct {
	Char  string `json:"char"`
	Value int    `json:"value"`
}
type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Hand []Tile `json:"hand"`
}
type Game struct {
	Players []Player `json:"players"`
	Board   []Tile   `json:"board"`
	Id      int      `json:"id"`
}

func GetIPAndUserAgent(r *http.Request) (ip string, user_agent string) {
	ip = r.RemoteAddr
	user_agent = r.UserAgent()

	return ip, user_agent

}

//routes handlers
func setupRoutes() {
	var games []Game
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
		wsHandler(ws, &games)
	})

	//sets up a web socket connection upgrading the http route
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "GET" {
			clientIP, _ := GetIPAndUserAgent(r)
			//create an instance of the game and send json to client
			fmt.Println(clientIP)
			newPlayer := Player{
				Id:   clientIP,
				Name: "",
				Hand: nil,
			}
			players := []Player{newPlayer}
			var board []Tile
			game := Game{Players: players, Board: board, Id: newGameId()}
			games = append(games, game)
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

func loadTiles() []Tile {
	file, err := os.ReadFile("ScrabbleTiles.txt")
	if err != nil {
		log.Println("could not load tiles")
	}
	s := strings.Split(string(file), "\n")
	var tiles []Tile
	for _, v := range s {
		char := string(v[0])
		value, err := strconv.Atoi(v[2 : len(v)-1])
		if err != nil {
			log.Println(err)
		}
		tiles = append(tiles, Tile{Char: char, Value: value})
	}
	for _, v := range tiles {
		fmt.Println(v)
	}
	return tiles
}

//join game using message from client
func wsHandler(conn *websocket.Conn, games *[]Game) {
	//tile bag should not be sent to the players to prevent cheating, initial loading of tile bag
	tiles := loadTiles()
	for {
		//waits for message from client to execute the loop
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		gameCode, err := strconv.Atoi(string(msg))
		if err != nil {
			log.Println(err)
			return
		}
		foundGameCode := false
		for i := 0; i < len(*games); i++ {
			if (*games)[i].Id == gameCode {
				foundGameCode = true
			}
		}
		if foundGameCode {
			replacementTiles := draw(7, &tiles)
			jsonPlayer, err := json.Marshal(
				Player{
					Name: "DEFAULT1",
					Hand: replacementTiles,
				},
			)
			err = conn.WriteMessage(messageType, jsonPlayer)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			err = conn.WriteMessage(messageType, []byte("no game code found"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func draw(numTiles int, tileBag *[]Tile) []Tile {
	removedTiles := remove(tileBag, numTiles)
	return removedTiles
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func remove[T any](slice *[]T, s int) []T {
	var removedTiles = (*slice)[:s]
	*slice = (*slice)[s:]
	return removedTiles
}

func main() {

	setupRoutes()
	fmt.Println("Starting server at port at :" + os.Getenv("PORT"))
	srv := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
