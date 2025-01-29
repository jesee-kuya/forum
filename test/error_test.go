package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesee-kuya/forum/backend/handler"
)

var TestCase1 = []struct {
	name       string
	errval     string
	statusCode int
}{
	{"Test1", "Not found", http.StatusNotFound},
	{"Test2", "Method not allowed", http.StatusMethodNotAllowed},
	{"Test2", "Method not allowed", http.StatusMethodNotAllowed},
}

func TestErrorHandler(t *testing.T) {
	for _, tc := range TestCase1 {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ErrorHandler(w, tc.errval, tc.statusCode)
			resp := w.Result()

			if resp.StatusCode != tc.statusCode {
				t.Errorf("Expected the statusCode %v but got the statuscode %v", http.StatusNotFound, resp.StatusCode)
			}
		})
	}
}
