package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"scrabble-backend/network"
	"time"
)

func remove[T any](slice *[]T, s int) []T {
	var removedTiles = (*slice)[:s]
	*slice = (*slice)[s:]
	return removedTiles
}

func main() {

	network.SetupRoutes()
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
