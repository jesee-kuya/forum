package openauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const RedirectBaseURL = "http://localhost:9000"

// generateStateCookie generates a random state and sets it as a cookie.
func generateStateCookie(w http.ResponseWriter) string {
	// Generate a random string
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating random state: %v", err)
		return ""
	}

	state := base64.URLEncoding.EncodeToString(b)

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}

// validateState checks if the state parameter in the URL matches the state parameter in the cookie.
func validateState(r *http.Request) error {
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		return err
	}

	// Check if the two states match
	if cookie.Value != state {
		log.Printf("State mismatch. Cookie: %s, State: %s", cookie.Value, state)
		return errors.New("invalid state")
	}
	return nil
}

/*
getGoogleUser takes a token and makes a request to the Google UserInfo API to retrieve the user's information. Returns a GoogleUser struct if the request is successful, or an error if the request fails.
*/
func getGoogleUser(token string) (*GoogleUser, error) {
	req, _ := http.NewRequest("GET", GoogleUserInfo, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response as a GoogleUser
	var user GoogleUser
	json.NewDecoder(resp.Body).Decode(&user)

	// Return the user if the request was successful
	return &user, nil
}
