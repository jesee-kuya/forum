package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var TestCase2 = []struct {
	name         string
	method       string
	endPoint     string
	expectedBody string
	code         int
}{
	{"Test1", http.MethodGet, "/upload", "Invalid request method", http.StatusMethodNotAllowed},
}

var TestCase3 = []struct {
	name        string
	content     []byte
	expectedExt string
	expectError bool
}{
	{"Valid JPEG", []byte("\xff\xd8\xff\xe0"), ".jpg", false},
	{"Valid PNG", []byte("\x89PNG\r\n\x1a\n"), ".png", false},
	{"Valid GIF", []byte("GIF87a"), ".gif", false},
	{"Valid WEBP", []byte("RIFF\x00\x00\x00\x00WEBPVP8 "), ".webp", false},
	{"Invalid Text File", []byte("This is a text file"), "", true},
	{"Empty File", []byte{}, "", true},
}

func TestUploadMedia(t *testing.T) {
	for _, tc := range TestCase2 {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endPoint, nil)
			w := httptest.NewRecorder()

			CreatePost(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.code {
				t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
			}

			if !strings.Contains(w.Body.String(), tc.expectedBody) {
				t.Errorf("Expected response body to contain %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func createTestFile(content []byte) (multipart.File, error) {
	return &mockMultipartFile{bytes.NewReader(content)}, nil
}

func TestValidateMimeType(t *testing.T) {
	for _, tc := range TestCase3 {
		t.Run(tc.name, func(t *testing.T) {
			file, err := createTestFile(tc.content)
			if err != nil {
				t.Fatalf("Error creating test file: %v", err)
			}

			ext, err := ValidateMimeType(file)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err)
			}

			if ext != tc.expectedExt {
				t.Errorf("Expected extension: %q, got: %q", tc.expectedExt, ext)
			}
		})
	}
}
