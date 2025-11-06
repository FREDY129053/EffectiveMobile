/*
Package repository implements the data access layer for subscription records.
*/
package repository

import (
	"database/sql"
	"fmt"
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

func (r *SubscriptionRepository) GetRecords(offset, size int) ([]models.Subscription, *int, error) {
	var records []models.Subscription

	var total int64
	if err := r.DB.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil, nil, err
	}

	totalPages := int((total + int64(size) - 1) / int64(size))

	if err := r.DB.Limit(size).Offset(offset).Find(&records).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil, nil, err
	}

	return records, &totalPages, nil
}

func (r *SubscriptionRepository) GetRecord(id uint) (*models.Subscription, error) {
	var record models.Subscription

	if err := r.DB.Take(&record, id).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

func (r *SubscriptionRepository) CreateRecord(serviceName string, startDate time.Time, price uint, userID uuid.UUID, endDate *time.Time) (*uint, error) {
	var newID uint

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		newRecord := models.Subscription{
			ServiceName: serviceName,
			Price:       price,
			UserID:      userID,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		if err := tx.Create(&newRecord).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return err
		}

		newID = newRecord.ID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &newID, nil
}

func (r *SubscriptionRepository) FullUpdateRecord(
	id, price uint,
	serviceName string,
	startDate time.Time,
	userID uuid.UUID,
	endDate *time.Time,
) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var toUpdateRecord models.Subscription

		if err := r.DB.Take(&toUpdateRecord, id).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return gorm.ErrRecordNotFound
		}

		toUpdateRecord.ServiceName = serviceName
		toUpdateRecord.Price = price
		toUpdateRecord.UserID = userID
		toUpdateRecord.StartDate = startDate
		toUpdateRecord.EndDate = endDate

		if err := tx.Save(&toUpdateRecord).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return err
		}

		return nil
	})

	return err
}

func (r *SubscriptionRepository) UpdateRecord(id uint, fields map[string]any) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var record models.Subscription

		if err := r.DB.Take(&record, id).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return gorm.ErrRecordNotFound
		}

		if err := tx.Model(&record).Updates(fields).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return err
		}

		return nil
	})

	return err
}

func (r *SubscriptionRepository) DeleteRecord(id uint) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.DB.Take(&models.Subscription{}, id).Error; err != nil {
			logger.PrintLog(err.Error(), "error")
			return gorm.ErrRecordNotFound
		}

		if err := tx.Delete(&models.Subscription{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *SubscriptionRepository) GetSubsSum(userID *uuid.UUID, serviceName *string, startDate, endDate string) *uint {
	var totalSum sql.NullInt64

	args := []any{startDate, endDate}
	nextPlaceholder := 3
	whereClauses := ""

	if userID != nil {
		whereClauses += fmt.Sprintf(" AND user_id = $%d", nextPlaceholder)
		args = append(args, *userID)
		nextPlaceholder++
	}

	if serviceName != nil {
		whereClauses += fmt.Sprintf(" AND service_name ILIKE $%d", nextPlaceholder)
		args = append(args, *serviceName)
		nextPlaceholder++
	}

	rawSQL := `
		SELECT
			COALESCE(SUM(
				(
					EXTRACT(
						MONTH FROM AGE(real_sub_end, real_sub_start)
					) 
					+ 
					EXTRACT(
						YEAR FROM AGE(real_sub_end, real_sub_start)
					) * 12
				) * price
			)) AS total_sum
		FROM (
			SELECT
				service_name,
				user_id,
				price, 
				start_date, 
				end_date, 
			
				CASE WHEN (
					EXTRACT(
						MONTH FROM AGE(start_date, $1::date)
					) 
					+ 
					EXTRACT(
						YEAR FROM AGE(start_date, $1::date)) * 12
					) > 0 THEN start_date ELSE $1::date END AS real_sub_start,
			
				CASE WHEN (
					EXTRACT(
						MONTH FROM AGE(end_date, $2::date)
					) 
					+ 
					EXTRACT(
						YEAR FROM AGE(end_date, $2::date)) * 12
					) < 0 THEN end_date ELSE $2::date END AS real_sub_end
			
			FROM subscriptions
			WHERE
				($1::date, $2::date) OVERLAPS 
				(start_date::date, end_date::date)` + whereClauses + `);`

	if err := r.DB.Raw(rawSQL, args...).Scan(&totalSum).Error; err != nil {
		logger.PrintLog(err.Error(), "error")
		return nil
	}

	total := uint(totalSum.Int64)
	return &total
}
