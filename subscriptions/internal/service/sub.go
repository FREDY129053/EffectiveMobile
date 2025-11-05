package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/schemas"
	"subscriptions/rest-service/pkg/logger"
	"time"

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

func (s *SubscriptionService) GetAllSubs(pageNumber, pageSize int) (*schemas.PaginationResponse, error) {
	offset := (pageNumber - 1) * pageSize
	records, totalPages, err := s.repository.GetRecords(offset, pageSize)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
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

	paginationInfo := schemas.Pagination{
		PageNumber: pageNumber,
		Size:       pageSize,
		TotalPages: *totalPages,
		HasNext:    pageNumber < *totalPages,
		HasPrev:    pageNumber > 1 && *totalPages > 0,
	}

	response := schemas.PaginationResponse{
		Subscriptions: result,
		Pagination:    paginationInfo,
	}

	return &response, nil
}

func (s *SubscriptionService) GetSub(id uint) (*schemas.FullSubInfo, error) {
	record, err := s.repository.GetRecord(id)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
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

	logger.PrintLog(fmt.Sprintf("Get record with ID = %d", id))
	return (*schemas.FullSubInfo)(record), nil
}

func (s *SubscriptionService) CreateSub(data schemas.CreateSub) (uint, error) {
	startDate, err := time.Parse("01-2006", data.StartDate)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
		return 0, err
	}

	var endDate *time.Time

	if data.EndDate != nil {
		t, err := time.Parse("01-2006", *data.EndDate)
		if err != nil {
			logger.PrintLog(err.Error(), "error")
			return 0, nil
		}
		endDate = &t
	}

	res, err := s.repository.CreateRecord(
		data.ServiceName, startDate, data.Price, data.UserID, endDate,
	)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
		return 0, err
	}

	logger.PrintLog("Subscription record created")
	return *res, nil
}

func (s *SubscriptionService) FullUpdateSub(id uint, data schemas.FullUpdateSub) error {
	startDate, err := time.Parse("01-2006", data.StartDate)
	if err != nil {
		return &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid start date format",
			Err:     err,
		}
	}

	var endDate *time.Time
	if data.EndDate != nil {
		t, err := time.Parse("01-2006", *data.EndDate)
		if err != nil {
			return &schemas.AppError{
				Code:    http.StatusBadRequest,
				Message: "Invalid end date format",
				Err:     err,
			}
		}
		endDate = &t
	}

	err = s.repository.FullUpdateRecord(
		id,
		data.Price,
		data.ServiceName,
		startDate,
		data.UserID,
		endDate,
	)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
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

	logger.PrintLog("Subscription updated")
	return nil
}

func (s *SubscriptionService) PatchUpdateSub(id uint, data schemas.PatchUpdateSub) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
		return &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid update data",
			Err:     err,
		}
	}

	var updateFields map[string]any
	if err = json.Unmarshal(jsonBytes, &updateFields); err != nil {
		logger.PrintLog(err.Error(), "error")
		return &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Failed to parse update fields",
			Err:     err,
		}
	}

	err = s.repository.UpdateRecord(id, updateFields)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
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

	logger.PrintLog("Subscription updated")
	return nil
}

func (s *SubscriptionService) DeleteSub(id uint) error {
	err := s.repository.DeleteRecord(id)
	if err != nil {
		logger.PrintLog(err.Error(), "error")
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

	logger.PrintLog("Subscription deleted")
	return nil
}

func (s *SubscriptionService) GetSubSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) (uint, error) {
	if *serviceName == "" {
		serviceName = nil
	}

	startDateParsed, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid start date format",
			Err:     err,
		}
	}

	endDateParsed, err := time.Parse("01-2006", endDate)
	if err != nil {
		return 0, &schemas.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid end date format",
			Err:     err,
		}
	}

	y1, m1, d1 := startDateParsed.Date()
	y2, m2, d2 := endDateParsed.Date()
	startDateSQL := fmt.Sprintf("%04d-%02d-%02d", y1, m1, d1)
	endDateSQL := fmt.Sprintf("%04d-%02d-%02d", y2, m2, d2)

	totalSum := s.repository.GetSubsSum(userID, serviceName, startDateSQL, endDateSQL)

	if totalSum == nil {
		logger.PrintLog("error get sum with this params", "error")
		return 0, &schemas.AppError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Cannot calculate sum of subscriptions",
			Err:     errors.New("returned nil sum from repo"),
		}
	}

	logger.PrintLog("Get sum")
	return *totalSum, nil
}
