package schemas

type APIError struct {
	Error string `json:"error"`
}

type CreateReturn struct {
	ID uint `json:"id"`
}

type MessageReturn struct {
	Message string `json:"message"`
}

type SumReturn struct {
	TotalSum uint `json:"total_sum"`
}