package models

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceName string    `json:"service_name" gorm:"size:150;not null"`
	Price       uint      `json:"price" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	StartDate   string    `json:"start_date" gorm:"not null"`
	EndDate     *string   `json:"end_date,omitempty"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	// TODO: Check month
	var dateRegex = regexp.MustCompile(`^\d{2}-\d{4}$`)
	checkDate := func(date string) bool { return dateRegex.MatchString(date) }

	if !checkDate(s.StartDate) {
		return errors.New("дата начала подписки не подходит под формат 'mm-yyyy'")
	}
	if s.EndDate != nil && !checkDate(*s.EndDate) {
		return errors.New("дата конца подписки не подходит под формат 'mm-yyyy'")
	}

	return nil
}
