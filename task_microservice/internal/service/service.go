package service

import (
	"context"
	"errors"
	"graduate_backend_task_microservice/internal/constant"
	"graduate_backend_task_microservice/internal/kafkaproducer"
	"graduate_backend_task_microservice/internal/minio"
	"graduate_backend_task_microservice/internal/model"
	"graduate_backend_task_microservice/internal/postgresql"
	"graduate_backend_task_microservice/internal/security"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
)

type Service struct {
	ctx           context.Context
	kafkaProducer *kafkaproducer.Producer
	minioClient   *minio.Client
	postgresql    *postgresql.PostgreSQL
	security      *security.Security
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

	secutityObj := security.NewSecurity()

	return &Service{
		ctx:           ctx,
		kafkaProducer: kafka,
		minioClient:   minioClient,
		postgresql:    psql,
		security:      secutityObj,
	}, nil
}

func (s *Service) TaskGetByUserUuid(token *string) ([]model.TaskInfo, error) {
	tokenSub, err := s.security.GetSubFromToken(token)
	if err != nil {
		return []model.TaskInfo{}, err
	}

	images, err := s.postgresql.TaskGetByUserUuid(tokenSub)
	if err != nil {
		return []model.TaskInfo{}, err
	}

	return images, nil
}

func (s *Service) GetImagesByTaskId(taskId int64, token *string) (model.TaskResponse, error) {
	taskInfo, err := s.postgresql.TaskGetById(taskId)
	if err != nil {
		return model.TaskResponse{}, err
	}

	tokenSub, err := s.security.GetSubFromToken(token)
	if err != nil {
		return model.TaskResponse{}, err
	}

	if tokenSub != taskInfo.UserUuid {
		return model.TaskResponse{}, errors.New("access denied for non owner")
	}

	images, err := s.postgresql.ImageGetByTaskId(taskId)
	if err != nil {
		return model.TaskResponse{}, err
	}

	var commonStatusId int64 = constant.StatusSuccessful

	isWork, isFailed := false, false

	for _, val := range images {
		switch val.StatusId {
		case constant.StatusInWork:
			isWork = true
		case constant.StatusFailed:
			isFailed = true
			break
		default:
		}
	}

	if isWork {
		commonStatusId = constant.StatusInWork
	}
	if isFailed {
		commonStatusId = constant.StatusFailed
	}

	return model.TaskResponse{
		TaskInfo:       taskInfo,
		CommonStatusId: commonStatusId,
		Images:         images,
		CreatedDT:      taskInfo.CreatedDT,
	}, nil
}

func (s *Service) Post(files *multipart.Form, width *int, height *int, targetFormat *string, quality *float64, tokenString *string) (int64, error) {
	if files == nil {
		return -1, errors.New("файл отсутствует")
	}

	userUuid, err := s.security.GetSubFromToken(tokenString)
	if err != nil {
		return -1, err
	}

	taskId, err := s.postgresql.TaskCreate(width, height, targetFormat, quality, userUuid)
	if err != nil {
		return -1, err
	}

	for i, v2 := range files.File["file"] {
		imageInfo := model.ImageInfo{
			TaskId:   taskId,
			Position: i + 1,
			StatusId: constant.StatusInWork,
		}

		formatSeparator := strings.LastIndex(v2.Filename, ".")
		imageInfo.Filename = v2.Filename[:formatSeparator]
		imageInfo.Format = strings.ToLower(v2.Filename[formatSeparator+1:])

		imageId, err := s.postgresql.ImageCreate(imageInfo)
		if err != nil {
			return -1, err
		}
		imageInfo.Id = imageId

		val, err := v2.Open()
		if err != nil {
			return -1, err
		}

		fileBytes, err := io.ReadAll(val)
		if err != nil {
			return -1, err
		}

		minioFilename := strconv.FormatInt(imageInfo.TaskId, 10) + "_" + strconv.Itoa(imageInfo.Position) + "." + imageInfo.Format

		s.minioClient.Upsert(fileBytes, minioFilename)

		imageRequest := model.ImageRequest{
			ImageInfo:    imageInfo,
			Width:        width,
			Height:       height,
			TargetFormat: targetFormat,
			Quality:      quality,
			UserUuid:     userUuid,
		}

		err = s.kafkaProducer.Write(&imageRequest)
		if err != nil {
			return -1, err
		}
	}
	return taskId, nil
}

func (s *Service) TaskUpdateStatus(imageStatus model.ImageStatus) error {
	err := s.postgresql.ImageUpdateStatus(imageStatus)

	if err != nil {
		return err
	}

	return nil
}
