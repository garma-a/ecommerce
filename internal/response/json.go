package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// JSON writes a JSON response with the given status code and payload.
// This is a production-grade helper that ensures consistent JSON responses.
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		return err
	}

	return nil
}

// Error writes a JSON error response with the given status code and message.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]interface{}{
		"error": message,
	})
}

// ValidationError writes a JSON validation error response with field-level errors.
func ValidationError(w http.ResponseWriter, message string, fields map[string]string) {
	JSON(w, http.StatusBadRequest, map[string]interface{}{
		"error":  message,
		"fields": fields,
	})
}

// OK writes a 200 OK JSON response.
func OK(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusOK, data)
}

// Created writes a 201 Created JSON response.
func Created(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusCreated, data)
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// InternalServerError writes a 500 Internal Server Error response.
func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}

// BadRequest writes a 400 Bad Request response.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound writes a 404 Not Found response.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}
