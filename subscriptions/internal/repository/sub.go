package repository

import (
	"subscriptions/rest-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewRepository(database *gorm.DB) SubscriptionRepository {
	return SubscriptionRepository{
		DB: database,
	}
}

func (r *SubscriptionRepository) GetAllRecords() ([]models.Subscription, error) {
	var records []models.Subscription
	if err := r.DB.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *SubscriptionRepository) GetRecord(id uint) (*models.Subscription, error) {
	var record models.Subscription

	if err := r.DB.Take(&record, id).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

func (r *SubscriptionRepository) CreateRecord(serviceName, startDate string, price uint, userID uuid.UUID, endDate *string) (*uint, error) {
	newRecord := models.Subscription{
		ServiceName: serviceName,
		Price: price,
		UserID: userID,
		StartDate: startDate,
		EndDate: endDate,
	}

	if err := r.DB.Create(&newRecord).Error; err != nil {
		return nil, err
	}

	return &newRecord.ID, nil
}

func (r *SubscriptionRepository) FullUpdateRecord(id, price uint, serviceName, startDate string, userID uuid.UUID, endDate *string) error {
	var toUpdateRecord models.Subscription

	if err := r.DB.Take(&toUpdateRecord).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	toUpdateRecord.ServiceName = serviceName
	toUpdateRecord.Price = price
	toUpdateRecord.UserID = userID
	toUpdateRecord.StartDate = startDate
	toUpdateRecord.EndDate = endDate
	
	if err := r.DB.Save(&toUpdateRecord).Error; err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) UpdateRecord(id uint, fields map[string]any) error {
	var company models.Subscription

	if err := r.DB.Take(&company, id).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.DB.Model(&company).Updates(fields).Error; err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) DeleteRecord(id uint) error {
	if err := r.DB.Take(&models.Subscription{}, id).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	return r.DB.Delete(&models.Subscription{}, id).Error
}