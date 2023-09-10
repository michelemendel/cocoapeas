package entities

import (
	"time"

	"github.com/michelemendel/cocopeas/utils"
)

type Message struct {
	TraceId   string
	Client    string
	Value     string
	Sender    string
	UpdatedAt time.Time
}

func MakeMessage(client, value, sender string) Message {
	return Message{
		TraceId:   string(utils.GenerateUUID()),
		Client:    client,
		Value:     value,
		Sender:    sender,
		UpdatedAt: time.Now(),
	}
}
