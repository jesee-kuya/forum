package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jesee-kuya/forum/backend/util"
)

var (
	requests = make(map[string]int)
	mu       sync.Mutex
)

func RateLimiter(next http.HandlerFunc, limit int, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ip := r.RemoteAddr
		requests[ip]++

		go func() {
			time.Sleep(duration)
			mu.Lock()
			defer mu.Unlock()
			requests[ip]--
		}()

		if requests[ip] > limit {
			log.Println("Too many requests from", ip)
			util.ErrorHandler(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
