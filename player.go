package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type JoinRoom struct{
	name string
	conn *websocket.Conn
}

type RegRoom struct{
	room *Room
	conn *websocket.Conn
}

// ConnectWs handles websocket requests from clients requests.
func ConnectWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	create := r.URL.Query()["create"]
	name := r.URL.Query()["name"]
	conn, err := upgrader.Upgrade(w, r, nil)

	if create[0] == "true" {
		theme := r.URL.Query()["theme"]
		grid := r.URL.Query()["grid"]
		size := r.URL.Query()["size"]
		intSize, _ := strconv.Atoi(size[0])
		if err != nil {
			log.Println(err)
			return
		}

		room := (NewRoom(name[0],conn, theme[0],grid[0],intSize))

		go room.RunRoom(wsServer)
		regRoom := RegRoom{
			room : room,
			conn : conn,
		}
		wsServer.registerRoom <- &regRoom
		

	

	} else {
		_, ok := wsServer.rooms[name[0]]
		
		if ok {
			fmt.Println("joining")
			joinRoom := JoinRoom{
				name:name[0],
				conn : conn,
			}
	
			wsServer.joinRoom <- &joinRoom

			

			// go func(){
			// 	for{
			// 		ww, jsonMessage, err := conn.ReadMessage()
			// 		fmt.Println("One22one")
			// 		fmt.Println("err:", err)				
			// 		fmt.Println("ww:",ww)
	
			// 		if err != nil {
			// 			fmt.Printf("unexpected close error: %v", err)
			// 			for key, val := range wsServer.rooms[name[0]].Players{
			// 				if val == conn {
			// 					wsServer.rooms[name[0]].leave(wsServer, conn,key )
			// 				}
			// 			}
			// 			break
			// 		}
			// 		wsServer.rooms[name[0]].handleNewMessage(jsonMessage)
			// 	}
			// }()
		}
		

		}
	}

