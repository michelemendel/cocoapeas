package main

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/michelemendel/cocopeas/utils"
	"golang.org/x/exp/slog"
)

func init() {
	utils.InitEnv()
}

// Read from Mastodon and produce to Kafka
type MessageHandler struct {
	kafkaProducer *kafka.Producer
	// todo: Add variable to something that reads messages from e.g. Mastodon
}

func makeMessageHandler() *MessageHandler {
	return &MessageHandler{
		kafkaProducer: makeProducer(),
	}
}

func main() {
	slog.Debug("[PRODUCER]:Init producer")

	mh := makeMessageHandler()
	defer mh.kafkaProducer.Close()
	topic := os.Getenv("TOPIC_MASTODON")
	// topic := "TOPIC_MASTODON"

	for i := 0; i < 3; i++ {
		now := time.Now().UnixMilli()
		mh.ProduceMessage(topic, fmt.Sprintf("PING_%d", now))
		time.Sleep(1 * time.Second)
	}
	mh.kafkaProducer.Flush(15 * 1000)
}

func (mh *MessageHandler) ProduceMessage(topic string, message string) {
	err := mh.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)
	if err != nil {
		slog.Error("[PRODUCER]:Failed to produce message", "err", err)
	}
}

func makeProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS")})
	if err != nil {
		panic(err)
	}
	go deliveryCheck(p)
	return p
}

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
