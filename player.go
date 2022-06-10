package main

type player struct {
	name  string
	score uint
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
