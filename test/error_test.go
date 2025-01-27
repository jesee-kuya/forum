package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleErrors(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	HandleErrors(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected the statusCode %v but got the statuscode %v", http.StatusOK, resp.StatusCode)
	}
}
