package clients

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/michelemendel/cocopeas/entities"
)

//
// This is a mock client that sends messages to topics.
// There may be a need for transforming the messages before sending them to the channel, so we'll have an ETL pipeline.
//

type Client struct {
	Name        string
	Ch          chan entities.Message
	MockMessage string
	MockSenders []string
}

func MakeClient(name, mockMsg string, mockSenders []string) Client {
	return Client{
		Name:        name,
		Ch:          make(chan entities.Message, 2),
		MockMessage: mockMsg,
		MockSenders: mockSenders,
	}
}

func (cl Client) MessageLoop() {
	for {
		val := fmt.Sprintf("%s_%d", cl.MockMessage, time.Now().UnixMilli())
		sender := cl.MockSenders[rand.Intn(len(cl.MockSenders))]
		cl.Ch <- entities.MakeMessage(cl.Name, val, sender)
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second) //To simulate a client that sends messages at different rates.
	}
}

func StartMockClientLoops(client Client) {
	client.MessageLoop()
}
