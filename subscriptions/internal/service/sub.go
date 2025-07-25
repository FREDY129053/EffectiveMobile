package service

import (
	"encoding/json"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/schemas"

	"github.com/google/uuid"
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
			ID: record.ID,
			ServiceName: record.ServiceName,
			Price: record.Price,
			UserID: record.UserID,
			StartDate: record.StartDate,
			EndDate: record.EndDate,
		}
	}

	return result, nil
}

func (s *SubscriptionService) GetSub(id uint) (*schemas.FullSubInfo, error) {
	record, err := s.repository.GetRecord(id)
	if err != nil {
		return nil, err
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

	return *res, err
}

func (s *SubscriptionService) FullUpdateSub(id uint, data schemas.FullUpdateSub) error {
	return s.repository.FullUpdateRecord(
		id,
		data.Price,
		data.ServiceName,
		data.StartDate,
		data.UserID,
		data.EndDate,
	)
}

func (s *SubscriptionService) PatchUpdateSub(id uint, data schemas.PatchUpdateSub) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var updateFields map[string]any
	if err = json.Unmarshal(jsonBytes, &updateFields); err != nil {
		return err
	}

	return s.repository.UpdateRecord(id, updateFields)
}

func (s *SubscriptionService) DeleteSub(id uint) error {
	return s.repository.DeleteRecord(id)
}

func (s *SubscriptionService) GetSubSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) *uint {
	return s.repository.GetSubsSum(userID, serviceName, startDate, endDate)
}