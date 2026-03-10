package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standard error response structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, statusCode int, message string) error {
	return WriteJSON(w, statusCode, ErrorResponse{Error: message})
}

// WriteSuccess writes a JSON success response with 200 OK status.
func WriteSuccess(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusOK, data)
}

// WriteCreated writes a JSON success response with 201 Created status.
func WriteCreated(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusCreated, data)
}
