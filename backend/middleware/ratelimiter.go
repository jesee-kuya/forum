package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jesee-kuya/forum/backend/util"
)

type client struct {
	timestamps []time.Time
	mu         sync.Mutex
}

var (
	clients     = make(map[string]*client)
	clientsLock sync.Mutex
)

func RateLimiter(next http.HandlerFunc, limit int, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			util.ErrorHandler(w, "Invalid IP address", http.StatusInternalServerError)
			return
		}

		clientsLock.Lock()
		c, ok := clients[ip]
		if !ok {
			c = &client{}
			clients[ip] = c
		}
		clientsLock.Unlock()

		c.mu.Lock()
		defer c.mu.Unlock()

		now := time.Now()
		// Filter out expired timestamps
		valid := []time.Time{}
		for _, t := range c.timestamps {
			if now.Sub(t) <= duration {
				valid = append(valid, t)
			}
		}
		c.timestamps = valid

		if len(c.timestamps) >= limit {
			util.ErrorHandler(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		c.timestamps = append(c.timestamps, now)
		next(w, r)
	}
}
