package main

import (
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/michelemendel/cocopeas/utils"
	"golang.org/x/exp/slog"
)

func init() {
	utils.InitEnv()
}

type ClientPresenter func(msg string)

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

func main() {
	slog.Info("[CONSUMER]:Init")
	topic := os.Getenv("TOPIC_MASTODON")
	topics := []string{topic}

	mp := makeMessagePresenter()
	mp.kafkaConsumer.SubscribeTopics(topics, nil)

	for mp.isRunning {
		mp.ConsumeMessage(PresentMessage)
	}
	mp.kafkaConsumer.Close()
}

func PresentMessage(msg string) {
	slog.Info("[CONSUMER]:PresentMessage", "msg", msg)
}

func (mp *MessagePresenter) ConsumeMessage(presenter ClientPresenter) {
	msg, err := mp.kafkaConsumer.ReadMessage(time.Second)
	if err == nil {
		slog.Info("[CONSUMER]:MessageReceived", "topic", msg.TopicPartition, "msg", string(msg.Value))
		presenter(string(msg.Value))
	} else if !err.(kafka.Error).IsTimeout() {
		slog.Info("[CONSUMER]", "err", err, "msg", msg)
	}
}

func makeConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		slog.Error("[CONSUMER]:Failed to create Kafka consumer. Shutting down.", "err", err)
		panic(err)
	}
	return c
}
