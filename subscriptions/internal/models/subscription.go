package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceName string     `json:"service_name" gorm:"size:150;not null;index:idx_subscriptions_service_name"`
	Price       uint       `json:"price" gorm:"not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index:idx_subscriptions_user_id"`
	StartDate   time.Time  `json:"start_date" gorm:"not null;type:date"`
	EndDate     *time.Time `json:"end_date,omitempty" gorm:"type:date"`
}
