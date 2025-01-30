package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesee-kuya/forum/backend/controllers"
	"github.com/jesee-kuya/forum/backend/models"
)

func TestGetAllPosts(t *testing.T) {
	t.Run("Success - Returns Posts", func(t *testing.T) {
		// Mock successful database call
		// repositories.GetPosts = mockGetPostsSuccess

		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		w := httptest.NewRecorder()

		controllers.GetAllPosts(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var posts []models.Post
		if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(posts) != 1 || posts[0].Body != "Test Post" {
			t.Errorf("Unexpected response body: %+v", posts)
		}
	})

	t.Run("Failure - Database Error", func(t *testing.T) {
		// Mock database error
		// repositories.GetPosts = mockGetPostsFailure

		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		w := httptest.NewRecorder()

		controllers.GetAllPosts(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}

		expectedBody := "Internal Server Error"
		if w.Body.String() != expectedBody {
			t.Errorf("Expected response body %q, got %q", expectedBody, w.Body.String())
		}
	})
}
