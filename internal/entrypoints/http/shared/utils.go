package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// WriteJSON serializes the given pointer to struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func WriteJSON(statusCode int, v any, w http.ResponseWriter) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(jsonData)
}

// DecodeJSON validate request body and decode result in the value pointed to by v.
func DecodeJSON(r *http.Request, v any) error {
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		if contentType != "application/json" {
			msg := "content-Type header is not application/json"
			response := ErrorResponse{
				StatusCode: http.StatusUnsupportedMediaType,
				Details:    msg,
			}
			return &response
		}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			fallthrough

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "request body contains badly-formed JSON"
			response := &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Details:    msg,
			}
			return response

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprint(
				"request body contains an invalid value",
			)
			response := &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Details:    msg,
			}
			return response

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			response := &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Details:    msg,
			}
			return response

		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			response := &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Details:    msg,
			}
			return response

		default:
			response := &ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Details:    "invalid payload",
			}
			return response
		}
	}

	return nil
}
