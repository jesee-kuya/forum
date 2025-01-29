package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jesee-kuya/forum/backend/handler"
)

var TestCase2 = []struct {
	name         string
	method       string
	endPoint     string
	expectedBody string
	code         int
}{
	{"Test1", http.MethodGet, "/upload", "Invalid reques method", http.StatusMethodNotAllowed},
}

func TestUploadMedia(t *testing.T) {
	for _, tc := range TestCase2 {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endPoint, nil)
			w := httptest.NewRecorder()

			handler.UploadMedia(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.code {
				t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
			}

			expectedBody := "Invalid request method"
			if !strings.Contains(w.Body.String(), expectedBody) {
				t.Errorf("Expected response body to contain %q, got %q", expectedBody, w.Body.String())
			}
		})
	}
}
