package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	JSON(w, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result map[string]string
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestJSON_WithStatusCreated(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]int{"id": 123}

	JSON(w, http.StatusCreated, data)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("something went wrong")

	Error(w, http.StatusBadRequest, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "something went wrong")
}

func TestError_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("database connection failed")

	Error(w, http.StatusInternalServerError, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection failed")
}
