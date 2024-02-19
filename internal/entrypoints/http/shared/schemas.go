package shared

type ResultResponse struct {
	Ok bool `json:"ok"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode" example:"400"`
	Details    string `json:"details" example:"error message"`
}

func (e *ErrorResponse) Error() string {
	return e.Details
}
