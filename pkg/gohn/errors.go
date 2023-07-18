package gohn

import (
	"fmt"
	"net/http"
)

type InvalidItemError struct {
	Message string
}

func (e InvalidItemError) Error() string {
	return fmt.Sprintf("invalid item: %v", e.Message)
}

type ResponseError struct {
	Response *http.Response
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("Error %d for %v", r.Response.StatusCode, r.Response.Request.URL)
}
