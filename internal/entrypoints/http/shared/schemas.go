package shared

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Details    string `json:"details"`
}

func (e *ErrorResponse) Error() string {
	return e.Details
}
