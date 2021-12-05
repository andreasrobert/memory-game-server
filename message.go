package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Action  string  `json:"action"`
	Message interface{}  `json:"message"`
	Sender int `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
