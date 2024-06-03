package handler

import (
	"gojek/library/database"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_GivenGETRequest_ThenReturnCorrectMessage(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	w := httptest.NewRecorder()

	db, _ := database.NewPostgresDB()
	bookHandler := NewBookHandler(db)
	bookHandler.HealthCheckHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "status code")

	expectedBody := `{"status":"OK"}`
	assert.JSONEq(t, expectedBody, w.Body.String(), "response body")
}
