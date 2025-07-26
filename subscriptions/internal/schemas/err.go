package schemas

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (a *AppError) Error() string {
	return fmt.Sprintf("%d: %s (%v)", a.Code, a.Message, a.Err)
}
