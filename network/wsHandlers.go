package network

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"scrabble-backend/game"
	"scrabble-backend/utils"
	"strconv"
)

func join(player game.Player, games *[]game.Game, conn *websocket.Conn, messageType int) {
	foundGameCode := false
	var currentGame *game.Game
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
		currentPlayer := game.Player{
			Name:       player.Name,
			Hand:       replacementTiles,
			Id:         conn.RemoteAddr().String(),
			Connection: conn,
			GamCode:    player.GamCode,
			Score:      0,
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
func wsHandler(conn *websocket.Conn, games *[]game.Game) {
	//tile bag should not be sent to the players to prevent cheating, initial loading of tile bag
	for {
		//waits for message from client to execute the loop
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var player game.Player
		err = json.Unmarshal(msg, &player)
		if err != nil {
			log.Println(err)
		}
		switch player.Action {
		case "join":
			join(player, games, conn, messageType)
		case "place":
			place(player, games, conn, messageType)
		case "replace":
			//todo
		case "pass":
			//todo
		case "reconnect":
			reconnect(player, games, conn, messageType)
		}
	}
}
func place(player game.Player, games *[]game.Game, conn *websocket.Conn, messageType int) {
	currentGame := getGame(player, games)
	jsonPlayer, err := json.Marshal(player)
	if err != nil {
		return
	}
	for _, p := range currentGame.Players {

		err := p.Connection.WriteMessage(messageType, jsonPlayer)
		if err != nil {
			return
		}
	}
}
func getGame(player game.Player, games *[]game.Game) *game.Game {
	var currentGame *game.Game
	for _, g := range *games {
		if strconv.Itoa(g.Id) == player.GamCode {
			currentGame = &g
		}
	}
	return currentGame
}
func resetConnection(currentGame *game.Game, player game.Player, conn *websocket.Conn) {
	for _, p := range currentGame.Players {
		if p.Name == player.Name {
			p.Connection = conn
		}
	}
}
func reconnect(player game.Player, games *[]game.Game, conn *websocket.Conn, messageType int) {
	currentGame := getGame(player, games)
	resetConnection(currentGame, player, conn)

}

func draw(numTiles int, tileBag *[]game.Tile) []game.Tile {
	removedTiles := utils.Remove(tileBag, numTiles)
	return removedTiles
}
