package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-shuffle"
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
	Connection *websocket.Conn //contains the web socket connection to allow for multi casting
	Id         string          `json:"id"`
	Action     string          `json:"action"`
	Name       string          `json:"name"`
	Hand       []Tile          `json:"hand"`
	GamCode    string          `json:"gameCode"`
}

type Game struct {
	Players []Player `json:"players"`
	Board   []Tile   `json:"board"`
	TileBag []Tile   `json:"tileBag"`
	Id      int      `json:"id"`
}

func GetIPAndUserAgent(r *http.Request) (ip string) {
	ip = strings.Split(r.RemoteAddr, ":")[0]
	return ip

}

//routes handlers
func setupRoutes() {
	var games []Game
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
			//create an instance of the game and send json to client
			var players []Player
			var board []Tile
			game := Game{Players: players, Board: board, Id: newGameId(), TileBag: loadTiles()}
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
	shuffle.Slice(tiles)
	return tiles
}

func join(player Player, games *[]Game, conn *websocket.Conn, messageType int) {
	foundGameCode := false
	var currentGame *Game
	for i := 0; i < len(*games); i++ {
		if strconv.Itoa((*games)[i].Id) == player.GamCode {
			currentGame = &(*games)[i]
			foundGameCode = true
		}
	}

	if foundGameCode {
		//same player cannot join twice
		for _, player := range currentGame.Players {
			if conn.RemoteAddr().String() == player.Id {
				return
			}
		}

		replacementTiles := draw(7, &currentGame.TileBag)
		currentPlayer := Player{
			Name:       player.Name,
			Hand:       replacementTiles,
			Id:         conn.RemoteAddr().String(),
			Connection: conn,
			GamCode:    player.GamCode,
		}
		currentGame.Players = append(currentGame.Players, currentPlayer)
		fmt.Printf("%#v", currentGame.Players)
		jsonGame, _ := json.Marshal(currentGame)
		if len(currentGame.Players) == 2 {
			for _, player := range currentGame.Players {
				err := player.Connection.WriteMessage(messageType, jsonGame)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}

	} else {
		err := conn.WriteMessage(messageType, []byte("not a valid game code"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

//join game using message from client
//both intiar and joiner will use this function to connect
func wsHandler(conn *websocket.Conn, games *[]Game) {
	//tile bag should not be sent to the players to prevent cheating, initial loading of tile bag
	for {
		//waits for message from client to execute the loop
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var player Player
		err = json.Unmarshal(msg, &player)
		if err != nil {
			log.Println(err)
		}
		switch player.Action {
		case "join":
			join(player, games, conn, messageType)
		case "place":
			//todo
		case "replace":
			//todo
		case "pass":
			//todo
		}
	}
}
func place() {

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
