package kafkaproducer

import (
	"graduate_backend_task_microservice/internal/model"
)

type MockProducer struct{}

func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

func (p *MockProducer) Write(imageInfo *model.ImageRequest) error {
	return nil
}
