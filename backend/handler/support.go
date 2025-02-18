package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

func createSession() string {
	sessionID := uuid.Must(uuid.NewV4()).String()
	SessionStore[sessionID] = make(map[string]interface{})
	return sessionID
}

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().UTC().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func getSessionData(sessionID string) (map[string]interface{}, error) {
	sessionData, exists := SessionStore[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return sessionData, nil
}

func setSessionData(sessionID string, key string, value interface{}) {
	SessionStore[sessionID][key] = value
}

func EnableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func deleteSession(userId int) {
	if len(SessionStore) == 0 {
		return
	}
	for k := range SessionStore {
		sessionData, _ := getSessionData(k)
		if len(sessionData) == 0 {
			continue
		}
		id := sessionData["userId"].(int)
		if id == userId {
			delete(SessionStore, k)
			return
		}
	}
}
