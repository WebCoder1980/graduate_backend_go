package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"graduate_backend_image_microservice_go/internal/minio"
	"log"
	"strings"
)

const consumerGroup = "group0"

func Consumer() {
	reader := GetKafkaReader()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		msgStr := string(msg.Value)
		offset := strings.Index(msgStr, EndFileName)
		file := msg.Value[(offset + len(EndFileName)):]
		filename := msgStr[:offset]

		minio.Upsert(filename, file)
	}
}

func GetKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{Host},
		Topic:   TopicName,
		GroupID: consumerGroup,
	})
}
