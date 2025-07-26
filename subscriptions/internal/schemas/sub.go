package schemas

import "github.com/google/uuid"

type CreateSub struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       uint      `json:"price" validate:"required,numeric,gt=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
}

type FullSubInfo struct {
	ID          uint      `json:"id" validate:"required"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       uint      `json:"price" validate:"required,numeric,gt=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
}

type FullUpdateSub struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       uint      `json:"price" validate:"required,numeric,gt=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
}

type PatchUpdateSub struct {
	ServiceName string    `json:"service_name,omitempty" swaggertype:"string" format:"nullable"`
	Price       uint      `json:"price,omitempty" swaggertype:"string" format:"nullable"`
	UserID      uuid.UUID `json:"user_id,omitempty" swaggertype:"string" format:"nullable"`
	StartDate   string    `json:"start_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
}
