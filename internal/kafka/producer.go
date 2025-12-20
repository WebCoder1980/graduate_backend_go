package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"graduate_backend_image_microservice_go/internal/constant"
	"io"
	"log"
	"mime/multipart"
	"os"
)

func Producer(file multipart.File, filename string) {
	ctx := context.Background()

	writer := GetKafkaWriter()
	defer writer.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	bytesResult := append([]byte(filename), EndFileName...)
	bytesResult = append(bytesResult, fileBytes...)

	err = writer.WriteMessages(ctx, kafka.Message{
		Value: bytesResult,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetKafkaWriter() kafka.Writer {
	return kafka.Writer{
		Addr:       kafka.TCP(os.Getenv("kafka_address")),
		Topic:      TopicName,
		BatchBytes: constant.FileMaxSize,
	}
}
