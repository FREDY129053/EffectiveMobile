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
	ServiceName *string    `json:"service_name,omitempty" swaggertype:"string" format:"nullable"`
	Price       *uint      `json:"price,omitempty" swaggertype:"string" format:"nullable"`
	UserID      *uuid.UUID `json:"user_id,omitempty" swaggertype:"string" format:"nullable"`
	StartDate   *string    `json:"start_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
	EndDate     *string    `json:"end_date,omitempty" swaggertype:"string" format:"nullable" validate:"omitempty,mm_yyyy_date"`
}

type Pagination struct {
	PageNumber int  `json:"page_number" validate:"required,numeric,gt=0"`
	Size       int  `json:"size" validate:"required,numeric,gt=0"`
	TotalPages int  `json:"total_pages" validate:"required,numeric,gt=0"`
	HasNext    bool `json:"has_next" validate:"required"`
	HasPrev    bool `json:"has_prev" validate:"required"`
}

type PaginationResponse struct {
	Subscriptions []FullSubInfo `json:"subscriptions" validate:"required"`
	Pagination    Pagination    `json:"pagination" validate:"required"`
}
