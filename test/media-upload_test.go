package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadMedia(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	UploadMedia(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected the statusCode %v but got the statuscode %v", http.StatusAccepted, resp.StatusCode)
	}
}
