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
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := r.DB.Create(&newRecord).Error; err != nil {
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

func (r *SubscriptionRepository) GetSubsSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) *uint {
	var total_sum uint

	// SQL-запрос из PgAdmin
	// select user_id, service_name, sum(price) as total_sum
	// from subscriptions
	// where 
	// service_name = '' and 
	// user_id = '' and
	// ('01-2025' <= start_date and start_date <= '07-2025'
	// or '01-2025' <= end_date and end_date <= '07-2025')
	// group by user_id, service_name

	seq := r.DB.Table("subscriptions").Select("SUM(price) AS total_sum")

	seq = seq.Where(`
		(? <= start_date AND start_date <= ? 
		OR ? <= end_date AND end_date <= ?)`,
		startDate, endDate, startDate, endDate,
	)

	if userID != nil {
		seq = seq.Where("user_id = ?", userID)
	}
	if serviceName != nil {
		seq = seq.Where("service_name = ?", serviceName)
	}

	if err := seq.Scan(&total_sum).Error; err != nil {
		return nil
	}

	return &total_sum
}
