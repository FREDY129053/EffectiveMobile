package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	repository repository.SubscriptionRepository
}

func NewService(repo repository.SubscriptionRepository) SubscriptionService {
	return SubscriptionService{
		repository: repo,
	}
}

func (s *SubscriptionService) GetAllSubs() ([]schemas.FullSubInfo, error) {
	records, err := s.repository.GetAllRecords()
	if err != nil {
		return nil, err
	}

	result := make([]schemas.FullSubInfo, len(records))

	for i, record := range records {
		result[i] = schemas.FullSubInfo{
			ID:          record.ID,
			ServiceName: record.ServiceName,
			Price:       record.Price,
			UserID:      record.UserID,
			StartDate:   record.StartDate,
			EndDate:     record.EndDate,
		}
	}

	return result, nil
}

func (s *SubscriptionService) GetSub(id uint) (*schemas.FullSubInfo, error) {
	record, err := s.repository.GetRecord(id)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, &schemas.AppError{
				Code:    http.StatusNotFound,
				Message: "Subscription not found!",
				Err:     err,
			}
		default:
			return nil, err
		}
	}

	return (*schemas.FullSubInfo)(record), nil
}

func (s *SubscriptionService) CreateSub(data schemas.CreateSub) (uint, error) {
	res, err := s.repository.CreateRecord(
		data.ServiceName, data.StartDate, data.Price, data.UserID, data.EndDate,
	)

	if err != nil {
		return 0, err
	}

	return *res, nil
}

func (s *SubscriptionService) FullUpdateSub(id uint, data schemas.FullUpdateSub) error {
	err := s.repository.FullUpdateRecord(
		id,
		data.Price,
		data.ServiceName,
		data.StartDate,
		data.UserID,
		data.EndDate,
	)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return &schemas.AppError{
				Code:    http.StatusNotFound,
				Message: "Subscription not found",
				Err:     err,
			}
		default:
			return &schemas.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update subscription",
				Err:     err,
			}
		}
	}

	return nil
}

func (s *SubscriptionService) PatchUpdateSub(id uint, data schemas.PatchUpdateSub) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid update data",
			Err:     err,
		}
	}

	var updateFields map[string]any
	if err = json.Unmarshal(jsonBytes, &updateFields); err != nil {
		return &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Failed to parse update fields",
			Err:     err,
		}
	}

	err = s.repository.UpdateRecord(id, updateFields)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return &schemas.AppError{
				Code:    http.StatusNotFound,
				Message: "Subscription not found",
				Err:     err,
			}
		default:
			return &schemas.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to patch update subscription",
				Err:     err,
			}
		}
	}

	return nil
}

func (s *SubscriptionService) DeleteSub(id uint) error {
	err := s.repository.DeleteRecord(id)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return &schemas.AppError{
				Code:    http.StatusNotFound,
				Message: "Subscription not found",
				Err:     err,
			}
		default:
			return &schemas.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete subscription",
				Err:     err,
			}
		}
	}

	return nil
}

func (s *SubscriptionService) GetSubSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) (uint, error) {
	if *serviceName == "" {
		serviceName = nil
	}

	totalSum := s.repository.GetSubsSum(userID, serviceName, startDate, endDate)

	if totalSum == nil {
		return 0, &schemas.AppError{
			Code: http.StatusUnprocessableEntity,
			Message: "Cannot calculate sum of subscriptions",
			Err: errors.New("returned nil sum from repo"),
		}
	}

	return *totalSum, nil
}
