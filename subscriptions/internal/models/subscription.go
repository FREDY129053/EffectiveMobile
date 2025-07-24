package models

import "github.com/google/uuid"


type Subscription struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceName string    `json:"service_name" gorm:"size:150;not null"`
	Price       uint      `json:"price" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	StartDate   string    `json:"start_date" gorm:"not null"`
	EndDate     *string   `json:"end_date,omitempty"`
}
