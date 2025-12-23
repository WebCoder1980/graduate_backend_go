package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"graduate_backend_task_microservice/internal/constant"
	"log"
	"os"
)

const (
	TopicName = "task_request"
)

type Producer struct {
	ctx         context.Context
	kafkaWriter *kafka.Writer
}

func NewProducer(ctx context.Context) *Producer {
	kafkaWriter := kafka.Writer{
		Addr:       kafka.TCP(os.Getenv("kafka_address")),
		Topic:      TopicName,
		BatchBytes: constant.FileMaxSize,
	}

	return &Producer{
		ctx:         ctx,
		kafkaWriter: &kafkaWriter,
	}
}

func (p *Producer) Write(filename string) {
	err := p.kafkaWriter.WriteMessages(p.ctx, kafka.Message{
		Value: []byte(filename),
	})
	if err != nil {
		log.Panic(err)
	}
}
