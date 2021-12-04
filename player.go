package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Player struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	rooms    map[*Room]bool
	Slot     int
}

func newPlayer(conn *websocket.Conn, wsServer *WsServer, name string) *Player {
	return &Player{
		ID:       uuid.New(),
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}

}

func (player *Player) readPump() {
	defer func() {
		player.disconnect()
	}()

	player.conn.SetReadLimit(maxMessageSize)
	player.conn.SetReadDeadline(time.Now().Add(pongWait))
	player.conn.SetPongHandler(func(string) error { player.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from player
	for {
		_, jsonMessage, err := player.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		player.handleNewMessage(jsonMessage)
	}

}



func (player *Player) disconnect() {
	
}

type JoinRoom struct{
	name string
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

		room := NewRoom(name[0],conn, theme[0],grid[0],intSize)

		go room.RunRoom(wsServer)
		wsServer.registerRoom <- room

	

	} else {
		fmt.Println("joining")
		joinRoom := JoinRoom{
			name:name[0],
			conn : conn,
		}

		wsServer.joinRoom <- &joinRoom
		
		go func(){
			for{
				ww, jsonMessage, err := conn.ReadMessage()
				fmt.Println("One22one")
				fmt.Println("err:", err)				
				fmt.Println("ww:",ww)

				if err != nil {
					// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						fmt.Printf("unexpected close error: %v", err)
						for key, val := range wsServer.rooms[name[0]].Players{
							if val == conn {
								wsServer.rooms[name[0]].leave(wsServer, conn,key )
							}
						}
					// }
					break
				}
				wsServer.rooms[name[0]].handleNewMessage(jsonMessage)
			}
		}()

		}
	}


func (player *Player) handleNewMessage(jsonMessage []byte) {


}



func (player *Player) isInRoom(room *Room) bool {
	log.Println("going")

	if _, ok := player.rooms[room]; ok {
		return true
	}

	return false
}

func (player *Player) notifyRoomJoined(room *Room, sender *Player) {
	message := Message{
		Action: RoomJoinedAction,
		// Target: room,
	}
	player.send <- message.encode()
}

func (player *Player) GetName() string {
	return player.Name
}
