package main

import (
	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Players    map[int]*Player `json:"players"`
	broadcast  chan *Message
	Grid	   string	 `json:"grid"`
	Theme	   string	 `json:"theme"`
	Size      int 		 `json:"size"`
}

// NewRoom creates a new Room
func NewRoom(name string, grid string,theme string, size int) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		Players:    make(map[int]*Player),
		broadcast:  make(chan *Message),
		Grid:    	grid,
		Theme:      theme,
		Size:		size,
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {

		}

	}
}

func(room *Room) takeSlot(player *Player){

}



func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}
