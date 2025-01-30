package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		errval     string
		statusCode int
		expected   string
	}{
		{
			name:       "Internal Server Error",
			errval:     "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			expected:   "500",
		},
		{
			name:       "Not Found",
			errval:     "Not Found",
			statusCode: http.StatusNotFound,
			expected:   "404",
		},
		{
			name:       "Bad Request",
			errval:     "Bad Request",
			statusCode: http.StatusBadRequest,
			expected:   "400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler
			_, err := http.NewRequest("GET", "/error", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the ErrorHandler function directly
			ErrorHandler(rr, tt.errval, tt.statusCode)

			// Check the status code
			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.statusCode)
			}

			// Check the response body
			expectedCode := tt.expected
			if rr.Body.String() != "" {
				t.Logf("Response body: %v", rr.Body.String())
			} else {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expectedCode)
			}
		})
	}
}
