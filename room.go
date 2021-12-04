package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	// "time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	// "github.com/gorilla/websocket"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Players    map[int]*websocket.Conn `json:"players"`
	broadcast  chan *Message
	Grid	   string	 `json:"grid"`
	Theme	   string	 `json:"theme"`
	Size       int 		 `json:"size"`
}

// NewRoom creates a new Room
func NewRoom(name string,conn *websocket.Conn, theme string,  grid string,size int) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		Players:    map[int]*websocket.Conn{1:conn},
		broadcast:  make(chan *Message),
		Theme:      theme,
		Grid:    	grid,
		Size:		size,
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom(server *WsServer) {

	for key, conn := range room.Players{
		fmt.Println("something something 42")
		defer func() {
			room.leave(server, conn, key)
			
		}()
		
	}

	for {
		ww, jsonMessage, err := room.Players[1].ReadMessage()
		fmt.Println("seven 9999:")
		fmt.Println("err:", err)
		fmt.Println("ww:",ww)

		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("unexpected close error: %v", err)
				room.leave(server,room.Players[1], 1)
			// }
			break
		}
		room.handleNewMessage(jsonMessage)
	}
	
}

func (room *Room) leave(server *WsServer, conn *websocket.Conn, key int) {
	fmt.Println("GOODBYE GOOD BYE :",key, room.Name)
	conn.Close()
	delete(room.Players,key)
	fmt.Println(len(room.Players))
	if len(room.Players) <= 0 || key == 1 {
		delete(server.rooms, room.Name)
		fmt.Println("room deleted")
		
	}else{
		fmt.Println("Player deleted")
		mess := map[string]interface{}{
			"something": "or rather",
			"number": strconv.Itoa(len(room.Players)),
		}

		msg := &Message{
			Action: "player-deleted",
			Message: mess,
			Sender: key,
		}
		messa := msg.encode()
		for _, conn := range room.Players {
			err := conn.WriteMessage(websocket.TextMessage, messa)
			if err != nil {
				fmt.Println("error:",err)
			}
		}
	}

}

func (room *Room) listen(server *WsServer) {

		for {
			_, jsonMessage, err := room.Players[1].ReadMessage()
			fmt.Println("key81:")

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("unexpected close error: %v", err)
				}
				break
			}
			room.handleNewMessage(jsonMessage)
		}
	 

}

func (room *Room) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		fmt.Println("Error on unmarshal JSON message: ", err)
		return
	}
	fmt.Println(message.Action)
	fmt.Println("message:")
	fmt.Println(message)
	
	for _,conn:= range room.Players {
		err := conn.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				fmt.Println("error:",err)
			}
	}

	

}


func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}
