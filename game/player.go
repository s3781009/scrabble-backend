package game

import "github.com/gorilla/websocket"

type Player struct {
	Connection *websocket.Conn //contains the web socket connection to allow for multi casting
	Id         string          `json:"id"`
	Action     string          `json:"action"`
	Name       string          `json:"name"`
	Hand       []Tile          `json:"hand"`
	GamCode    string          `json:"gameCode"`
	Score      int             `json:"score"`
}

//place Tiles on the board and remove from tile bag
func placeTile(removedTiles []byte, tiles TileBag) {
	for i, tile := range removedTiles {
		for _, tileinBag := range tiles.tiles {
			if tile == tileinBag {
				tiles.tiles = append(tiles.tiles[:i], tiles.tiles[i+1:]...)
				break
			}
		}

	}

}
