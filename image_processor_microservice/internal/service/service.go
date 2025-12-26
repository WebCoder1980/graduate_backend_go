package service

import (
	"context"
	"errors"
	"graduate_backend_image_processor_microservice/internal/constant"
	"graduate_backend_image_processor_microservice/internal/kafkaproducer"
	"graduate_backend_image_processor_microservice/internal/minio"
	"graduate_backend_image_processor_microservice/internal/model"
	"graduate_backend_image_processor_microservice/internal/postgresql"
	"strconv"
)

type Service struct {
	ctx           context.Context
	postgresql    *postgresql.PostgreSQL
	minioClient   *minio.Client
	kafkaProducer *kafkaproducer.Producer
}

func NewService(ctx context.Context) (*Service, error) {
	psql, err := postgresql.NewPostgreSQL()
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	kafka, err := kafkaproducer.NewProducer(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		ctx:           ctx,
		postgresql:    psql,
		minioClient:   minioClient,
		kafkaProducer: kafka,
	}, nil
}

func (s *Service) ImageProcessor(imageInfo *model.ImageInfo) error {
	minioFilename := strconv.FormatInt(imageInfo.TaskId, 10) + "_" + strconv.Itoa(imageInfo.Position) + "." + imageInfo.Format

	source, err := s.minioClient.Get(minioFilename)
	if err != nil {
		err2 := s.ImageProcessorKafkaWrite(imageInfo, constant.StatusFailed)
		if err2 != nil {
			return errors.New(err.Error() + "; " + err2.Error())
		}
		return nil
	}

	// TODO processing

	imageId, err := s.postgresql.ImageCreate(*imageInfo)
	if err != nil {
		err2 := s.ImageProcessorKafkaWrite(imageInfo, constant.StatusFailed)
		if err2 != nil {
			return errors.New(err.Error() + "; " + err2.Error())
		}
		return nil
	}
	imageInfo.Id = imageId

	err = s.minioClient.Upsert(source, minioFilename)
	if err != nil {
		err2 := s.ImageProcessorKafkaWrite(imageInfo, constant.StatusFailed)
		if err2 != nil {
			return errors.New(err.Error() + "; " + err2.Error())
		}
		return nil
	}

	err = s.ImageProcessorKafkaWrite(imageInfo, constant.StatusSuccessful)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImageProcessorKafkaWrite(imageInfo *model.ImageInfo, imageStatusId int64) error {
	imageStatus := model.ImageStatus{
		TaskId:   imageInfo.TaskId,
		Position: imageInfo.Position,
		StatusId: imageStatusId,
	}

	err := s.kafkaProducer.Write(imageStatus)
	if err != nil {
		return err
	}

	return nil
}
