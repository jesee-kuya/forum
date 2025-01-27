package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleErrors(t *testing.T) {
	str := "Not found"
	code := 404
	w := httptest.NewRecorder()
	ErrorHandler(w, str, code)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected the statusCode %v but got the statuscode %v", http.StatusNotFound, resp.StatusCode)
	}
}
