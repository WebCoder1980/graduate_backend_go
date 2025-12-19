package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"graduate_backend_image_microservice_go/internal/minio"
	"log"
	"strings"
)

func Consumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "image-topic",
		GroupID: "group0",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Ошибка при получении:", err)
		}

		msgStr := string(msg.Value)
		offset := strings.Index(msgStr, EndFileName)
		file := msg.Value[(offset + len(EndFileName)):]
		filename := msgStr[:offset]

		minio.Upsert(filename, file)
	}
}
