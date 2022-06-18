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

	//search fro the game code
	for i := 0; i < len(*games); i++ {
		if strconv.Itoa((*games)[i].Id) == player.GamCode {
			currentGame = &(*games)[i]
			foundGameCode = true
		}
	}

	if foundGameCode {
		fmt.Println("game code found")
		//same player cannot join twice
		for _, player := range currentGame.Players {
			if conn.RemoteAddr().String() == player.Id {
				fmt.Println("smae player")
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
		if len(currentGame.Players) == 1 {
			currentPlayer.Turn = true
		} else {
			currentPlayer.Turn = false
		}
		currentGame.Players = append(currentGame.Players, currentPlayer)
		jsonGame, _ := json.Marshal(currentGame)

		if len(currentGame.Players) == 2 {
			fmt.Println("2 players joined")
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
		//var sindPing = func() {
		//	for {
		//		//time.AfterFunc(3*time.Second, func() {
		//		//	for _, game := range *games {
		//		//		for _, player := range game.Players {
		//		//			if player.Connection != nil {
		//		//				msg := []byte("hello")
		//		//				player.Connection.WriteMessage(1, msg)
		//		//			}
		//		//		}
		//		//	}
		//		//})
		//	}
		//}
		//go sindPing()
		//waits for message from client to execute the loop
		messageType, msg, err := conn.ReadMessage()
		fmt.Println("message received")
		if err != nil {
			log.Println(err)
			return
		}
		var player game.Player
		err = json.Unmarshal(msg, &player)

		if err != nil {
			log.Println(err)
			return
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
		case "ping":
			sendPong(conn)

		}

	}
}
func sendPong(conn *websocket.Conn) {
	err := conn.WriteMessage(1, []byte("pong"))
	if err != nil {
		log.Println("error sendingg poing")
		return
	}
}

func place(message game.Player, games *[]game.Game, conn *websocket.Conn, messageType int) {

	currentGame := getGame(message, games)

	//reset player hand to 7 cards by drawing from the tile bag
	newTiles := draw(7-len(message.Hand), &currentGame.TileBag)
	for _, player := range currentGame.Players {
		//update the tile bag for both players
		player.Board = message.Board
		if player.Name == message.Name {

			player.Hand = append(message.Hand, newTiles...)
			fmt.Println(player.Hand)
			//set the player who has just placed the tile to have their turn off
			player.Turn = false

		} else {
			player.Turn = true
		}
		jsonHand, _ := json.Marshal(player.Hand)
		player.Connection.WriteMessage(messageType, jsonHand)
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

func reconnect(player game.Player, games *[]game.Game, conn *websocket.Conn, messageType int) {
	currentGame := getGame(player, games)
	for _, p := range currentGame.Players {
		if p.Name == player.Name {
			p.Connection = conn
		}
	}
	jsonGame, _ := json.Marshal(currentGame)
	err := conn.WriteMessage(messageType, jsonGame)
	if err != nil {
		log.Println(err)
		return
	}

}

func draw(numTiles int, tileBag *[]game.Tile) []game.Tile {
	removedTiles := utils.Remove(tileBag, numTiles)
	return removedTiles
}
