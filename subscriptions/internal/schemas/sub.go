package schemas

import "github.com/google/uuid"

type CreateSub struct {
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable"`
}

type FullSubInfo struct {
	ID          uint      `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable"`
}

type FullUpdateSub struct {
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable"`
}

type PatchUpdateSub struct {
	ServiceName string    `json:"service_name,omitempty" swaggertype:"string" format:"nullable"`
	Price       uint      `json:"price,omitempty" swaggertype:"string" format:"nullable"`
	UserID      uuid.UUID `json:"user_id,omitempty" swaggertype:"string" format:"nullable"`
	StartDate   string    `json:"start_date,omitempty" swaggertype:"string" format:"nullable"`
	EndDate     *string   `json:"end_date,omitempty" swaggertype:"string" format:"nullable"`
}
