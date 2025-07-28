package repository

import (
	"subscriptions/rest-service/internal/models"
	"subscriptions/rest-service/pkg/logger"
	"time"

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
		logger.PrintLog(err.Error(), "error")
		return nil, err
	}

	return records, nil
}

func (r *SubscriptionRepository) GetRecord(id uint) (*models.Subscription, error) {
	var record models.Subscription

	if err := r.DB.Take(&record, id).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

func (r *SubscriptionRepository) CreateRecord(serviceName, startDate string, price uint, userID uuid.UUID, endDate *string) (*uint, error) {
	newRecord := models.Subscription{
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := r.DB.Create(&newRecord).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil, err
	}

	return &newRecord.ID, nil
}

func (r *SubscriptionRepository) FullUpdateRecord(
	id, price uint,
	serviceName, startDate string,
	userID uuid.UUID,
	endDate *string,
) error {
	var toUpdateRecord models.Subscription

	if err := r.DB.Take(&toUpdateRecord).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return gorm.ErrRecordNotFound
	}

	toUpdateRecord.ServiceName = serviceName
	toUpdateRecord.Price = price
	toUpdateRecord.UserID = userID
	toUpdateRecord.StartDate = startDate
	toUpdateRecord.EndDate = endDate

	if err := r.DB.Save(&toUpdateRecord).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return err
	}

	return nil
}

func (r *SubscriptionRepository) UpdateRecord(id uint, fields map[string]any) error {
	var company models.Subscription

	if err := r.DB.Take(&company, id).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return gorm.ErrRecordNotFound
	}

	if err := r.DB.Model(&company).Updates(fields).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return err
	}

	return nil
}

func (r *SubscriptionRepository) DeleteRecord(id uint) error {
	if err := r.DB.Take(&models.Subscription{}, id).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return gorm.ErrRecordNotFound
	}

	return r.DB.Delete(&models.Subscription{}, id).Error
}

func (r *SubscriptionRepository) GetSubsSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) *uint {
	var subsInfo []models.Subscription
	var total_sum uint
	seq := r.DB.Table("subscriptions").Select("*")

	seq = seq.Where(`
		(
			(
					TO_DATE(?, 'MM-YYYY') <= TO_DATE(start_date, 'MM-YYYY') 
				AND 
					TO_DATE(start_date, 'MM-YYYY') <= TO_DATE(?, 'MM-YYYY') 
			)	
			OR 
			(
					TO_DATE(?, 'MM-YYYY') <= TO_DATE(end_date, 'MM-YYYY') 
				AND 
					TO_DATE(end_date, 'MM-YYYY') <= TO_DATE(?, 'MM-YYYY') 
			)
		)`,
		startDate, endDate, startDate, endDate,
	)

	if userID != nil {
		seq = seq.Where("user_id = ?", userID)
	}
	if serviceName != nil {
		seq = seq.Where("service_name ILIKE ?", serviceName)
	}

	if err := seq.Scan(&subsInfo).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil
	}

	// Подсчет суммы стоимости подписок
	startDateDate, _ := time.Parse("01-2006", startDate)
	endDateDate, _ := time.Parse("01-2006", endDate)
	var nullDate time.Time
	
	for _, sub := range subsInfo {
		subStartDate, _ := time.Parse("01-2006", sub.StartDate)

		var subEndDate time.Time
		if sub.EndDate != nil {
			subEndDate, _ = time.Parse("01-2006", *sub.EndDate)
		}

		if subStartDate.Before(startDateDate) {
			subStartDate = startDateDate
		}

		if endDateDate.Before(subEndDate) || subEndDate.Equal(nullDate) {
			subEndDate = endDateDate
		}

		datesDiffMonths := int64(subEndDate.Sub(subStartDate).Hours()/24/30) + 1

		total_sum += uint(datesDiffMonths * int64(sub.Price))
	}

	return &total_sum 
}
