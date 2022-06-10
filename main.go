package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//routes handlers
func setupRoutes() {

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello!")
		if err != nil {
			fmt.Println(err)
		}
	})
	//sets up a wb conenction upgrading hte
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		log.Println("succesfully connected ...")
		newWsConnection(ws)
	})
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
		if err := conn.WriteMessage(messageType, p); err != nil {

			log.Println(err)
		}
	}

}

func main() {

	setupRoutes()
	fmt.Println("Starting server at port at " + os.Getenv("PORT"))
	srv := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
