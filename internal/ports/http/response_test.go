package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	// Act
	err := WriteJSON(w, http.StatusOK, data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

func TestWriteError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	// Act
	err := WriteError(w, http.StatusBadRequest, "invalid request")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request", response.Error)
}

func TestWriteSuccess(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	data := map[string]int{"count": 42}

	// Act
	err := WriteSuccess(w, data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]int
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 42, response["count"])
}

func TestWriteCreated(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	data := map[string]string{"id": "123"}

	// Act
	err := WriteCreated(w, data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "123", response["id"])
}

func TestWriteError_InternalServerError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	// Act
	err := WriteError(w, http.StatusInternalServerError, "database connection failed")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "database connection failed", response.Error)
}

func TestWriteError_NotFound(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	// Act
	err := WriteError(w, http.StatusNotFound, "resource not found")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "resource not found", response.Error)
}
