package game

import (
	"fmt"
	"github.com/shogo82148/go-shuffle"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	Players []Player `json:"players"`
	Board   []Tile   `json:"board"`
	TileBag []Tile   `json:"tileBag"`
	Id      int      `json:"id"`
}

func NewGame() Game {
	return Game{
		Players: nil,
		Board:   nil,
		TileBag: loadTiles(),
		Id:      newGameId(),
	}
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

	shuffle.Slice(tiles)
	return tiles
}
