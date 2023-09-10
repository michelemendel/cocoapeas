package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/michelemendel/cocopeas/entities"
	"github.com/michelemendel/cocopeas/utils"
	"golang.org/x/exp/slog"
)

func init() {
	utils.InitEnv()
}

//
// Consume messages from Kafka, which can be sent to various clients.
//

func main() {
	slog.Info("[CONSUMER]:start...")
	topics := []string{"mastodon", "other"}

	mp := makeMessagePresenter()
	mp.kafkaConsumer.SubscribeTopics(topics, nil)

	for mp.isRunning {
		mp.ConsumeMessage()
	}
	mp.kafkaConsumer.Close()
}

type MessagePresenter struct {
	kafkaConsumer *kafka.Consumer
	isRunning     bool
}

func makeMessagePresenter() *MessagePresenter {
	return &MessagePresenter{
		kafkaConsumer: makeConsumer(),
		isRunning:     true,
	}
}

// Here we read messages from Kafka and print it to the console.
// In a real-world scenario, we would do something more useful with the message, like sending it to a web socket to deliver to clients.
func (mp *MessagePresenter) ConsumeMessage() {
	kafkaMsg, err := mp.kafkaConsumer.ReadMessage(time.Second)
	if err == nil {
		var msg entities.Message
		json.Unmarshal(kafkaMsg.Value, &msg)
		slog.Info("[CONSUMER]:MessageReceived", "topic", kafkaMsg.TopicPartition, "msg", msg)
	} else if !err.(kafka.Error).IsTimeout() {
		slog.Info("[CONSUMER]", "err", err, "msg", kafkaMsg)
	}
}

func makeConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	})
	if err != nil {
		slog.Error("[CONSUMER]:Failed to create Kafka consumer. Shutting down.", "err", err)
		panic(err)
	}
	return c
}
