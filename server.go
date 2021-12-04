package main

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	broadcast 		chan []byte
	registerRoom 	chan *Room
	deleteRoom		chan *Room
	joinRoom		chan *JoinRoom
	rooms     		map[string]*Room
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		broadcast: 		make(chan []byte),
		registerRoom:	make(chan *Room),
		deleteRoom:		make(chan *Room),
		joinRoom:		make(chan *JoinRoom),
		rooms:     		make(map[string]*Room),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {
		case room := <-server.registerRoom:
			fmt.Println("add room to server")
			fmt.Println("room : ", room)
			fmt.Println("room name:",room.Name)
			_, ok := server.rooms[room.Name]
			if ok {
				fmt.Println("room already being used")
			}else{
				server.rooms[room.Name]= room
				size1 := room.Size
				slotUsed := len(room.Players)
				fmt.Println("size: ",size1)
				fmt.Println("slotused: ",slotUsed)
			}
			fmt.Println("room in server:")
			for _, value:= range server.rooms{
				
				fmt.Println(value.Name)
			}
			fmt.Println("==================")

		case joinRoom := <-server.joinRoom:
			room, ok := server.rooms[joinRoom.name]
			if ok {
				size := room.Size
				slotUsed := len(room.Players)
				fmt.Println("size:",size)
				fmt.Println("slotUsed:",slotUsed)
				if size >= slotUsed + 1 {
					fmt.Println("get in here")
					room.Players[slotUsed+1] = joinRoom.conn
					
					fmt.Println("server52: ",room.Players)


					mess := map[string]interface{}{
						"grid": room.Grid,
						"theme": room.Theme,
						"size": size,
						"filled": strconv.Itoa(slotUsed+1),
					}
	
					msg := &Message{
						Action: "player-added",
						Message: mess,
						Sender: slotUsed+1,
					}

					messa := msg.encode()

					for _, conn := range room.Players {
						err := conn.WriteMessage(websocket.TextMessage, messa)
						if err != nil {
							fmt.Println("error:",err)
						}
					}

				}else{
					fmt.Println("not enough room")
				}

			}else{
				fmt.Println("room doesn't exist")
			}


		}

	}
}

// func (server *WsServer) createRoom(name string, grid string,theme string, size int) *Room {
// 	room := NewRoom(name, grid, theme, size)
// 	go room.RunRoom()
// 	server.rooms[room] = true

// 	return room
// }
