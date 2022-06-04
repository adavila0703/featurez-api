package api

import "fmt"

type apiError struct {
	message string
}

func (a *apiError) Error() string {
	return fmt.Sprintf("%s", a.message)
}

var (
	ErrInvalidParameter     = &apiError{message: "error: invalid parameters"}
	ErrInvalidJSON          = &apiError{message: "error: invalid JSON"}
	ErrFeatureAlreadyExists = &apiError{message: "error: Feature Flag already exists use UpdateFeature endpoint to update your FF"}
)
