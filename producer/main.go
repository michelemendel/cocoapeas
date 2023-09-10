package main

import (
	"encoding/json"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/michelemendel/cocopeas/clients"
	"github.com/michelemendel/cocopeas/entities"
	"github.com/michelemendel/cocopeas/utils"
	"golang.org/x/exp/slog"
)

//
// This starts the Kafka producer and mock clients.
//

func init() {
	utils.InitEnv()
}

func main() {
	slog.Debug("[PRODUCER]:start...")
	mh := makeMessageHandler()
	defer mh.kafkaProducer.Close()

	// Setup mock mockClients
	mockClients := []clients.Client{
		clients.MakeClient("mastodon", "Mastodon_message", []string{"abe", "bob", "carl"}),
		clients.MakeClient("other", "Other_Message", []string{"anna", "bea", "carol"}),
	}

	for _, client := range mockClients {
		go messageReceiver(mh, client)
		go clients.StartMockClientLoops(client)
	}

	// Just to keep the main thread alive. Not for production.
	<-make(chan struct{})
}

type MessageHandler struct {
	kafkaProducer *kafka.Producer
}

func makeMessageHandler() *MessageHandler {
	return &MessageHandler{
		kafkaProducer: makeProducer(),
	}
}

func messageReceiver(mh *MessageHandler, client clients.Client) {
	for msg := range client.Ch {
		mh.ProduceMessage(msg)
	}
}

// Send a message to the Kafka topic
func (mh *MessageHandler) ProduceMessage(message entities.Message) {
	msgAsBytes, err := json.Marshal(message)
	if err != nil {
		slog.Error("[PRODUCER]:Failed to marshal message", "err", err)
		return
	}

	err = mh.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &message.Client,
			Partition: kafka.PartitionAny,
		},
		Value: msgAsBytes,
	}, nil)
	if err != nil {
		slog.Error("[PRODUCER]:Failed to produce message", "err", err)
	}
}

func makeProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	})
	if err != nil {
		panic(err)
	}
	go deliveryCheck(p)
	return p
}

// Feedback on the delivery of messages
func deliveryCheck(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				slog.Error("[PRODUCER]:DeliveryFailed:", "err", ev.TopicPartition)
			} else {
				slog.Info("[PRODUCER]:DeliveryOk", "topic", ev.TopicPartition, "msg", string(ev.Value))
			}
		}
	}
}
