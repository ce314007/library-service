package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandler_GivenGetRequest_ThenReturnCorrectMessage(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ping)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status code not as expected")
	assert.Equal(t, `{"message":"pong"}`, strings.TrimSpace(rr.Body.String()), "response body not as expected")
}
