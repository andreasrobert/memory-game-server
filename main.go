package main

import (
	"flag"
	"log"
	"net/http"
	// "os"
)

var addr = flag.String("addr", ":4000", "http server address")

func main() {
	flag.Parse()
	// port := os.Getenv("PORT")
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ConnectWs(wsServer, w, r)
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
