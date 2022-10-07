package gohn

import (
	"fmt"
	"net/http"
)

type ErrInvalidItem struct {
	Message string
}

func (e ErrInvalidItem) Error() string {
	return fmt.Sprintf("invalid item: %v", e.Message)
}

type ErrResponse struct {
	Response *http.Response
}

func (r *ErrResponse) Error() string {
	return fmt.Sprintf("Error %d for %v", r.Response.StatusCode, r.Response.Request.URL)
}
