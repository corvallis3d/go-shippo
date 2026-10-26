package errors

import "fmt"

type APIError struct {
	Status       int
	ResponseBody []byte
	// Err preserves a response-body read failure, including cancellation/timeout.
	Err error
}

func (e *APIError) Unwrap() error { return e.Err }

func (e *APIError) Error() string {
	if len(e.ResponseBody) > 0 {
		return fmt.Sprintf("Error returned from API: [%d] %s", e.Status, e.ResponseBody)
	} else {
		return fmt.Sprintf("Error returned from API: [%d]", e.Status)
	}
}
