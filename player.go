package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"time"

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

func (player *Player) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		player.conn.Close()
	}()
	for {
		select {
		case message, ok := <-player.send:
			player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				player.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := player.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(player.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-player.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (player *Player) disconnect() {
	
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

}

func (player *Player) handleNewMessage(jsonMessage []byte) {


}

func (player *Player) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	player.joinRoom(roomName, nil)
}

func (player *Player) handleLeaveRoomMessage(message Message) {

}


func (player *Player) joinRoom(roomName string, sender *Player) {

}

func (player *Player) isInRoom(room *Room) bool {
	if _, ok := player.rooms[room]; ok {
		return true
	}

	return false
}

func (player *Player) notifyRoomJoined(room *Room, sender *Player) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
	}

	player.send <- message.encode()
}

func (player *Player) GetName() string {
	return player.Name
}
