package kafkaproducer

import (
	"graduate_backend_image_processor_microservice/internal/model"
)

type MockProducer struct{}

func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

func (p *MockProducer) Write(imageStatus model.ImageStatus) error {
	return nil
}
