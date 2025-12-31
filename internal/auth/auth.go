package auth

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Mode          string
	Header        string
	AllowedValues map[string]struct{}
	SharedHeader  string
	SharedSecret  string
	Enabled       bool
}

// Read env variables
func LoadConfig() *Config {
	cfg := &Config{
		Mode:         os.Getenv("AUTH_MODE"),
		Header:       os.Getenv("AUTH_HEADER"),
		SharedHeader: os.Getenv("AUTH_SHARED_HEADER"),
		SharedSecret: os.Getenv("AUTH_SHARED_KEY"),
	}

	allowed := os.Getenv("AUTH_ALLOWED")
	cfg.AllowedValues = make(map[string]struct{})
	for _, v := range strings.Split(allowed, ",") {
		v = strings.TrimSpace(strings.ToLower(v))
		if v != "" {
			cfg.AllowedValues[v] = struct{}{}
		}
	}

	// Decide if auth is enabled
	cfg.Enabled = cfg.Mode != ""

	if !cfg.Enabled {
		log.Println("auth disabled (AUTH_MODE not set)")
	}

	if cfg.Mode == "shared" && cfg.SharedSecret == "" {
		log.Println("warning: shared auth mode enabled but EXPLORER_AUTH_KEY not set")
	}

	return cfg
}

// Enforce auth
func Middleware(cfg *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch cfg.Mode {
			case "header":
				val := strings.TrimSpace(strings.ToLower(r.Header.Get(cfg.Header)))
				if val == "" {
					log.Printf("auth failed: missing header %s from %s", cfg.Header, r.RemoteAddr)
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}

				if len(cfg.AllowedValues) > 0 {
					if _, ok := cfg.AllowedValues[val]; !ok {
						log.Printf("auth forbidden: %s not in allowed list (remote=%s)", val, r.RemoteAddr)
						http.Error(w, "forbidden", http.StatusForbidden)
						return
					}
				}

			case "shared":
				key := r.Header.Get(cfg.SharedHeader)
				if key != cfg.SharedSecret {
					log.Printf("auth failed: invalid shared key from %s", r.RemoteAddr)
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}

			default:
				log.Printf("auth misconfigured: unknown mode %q", cfg.Mode)
				http.Error(w, "auth misconfigured", http.StatusInternalServerError)
				return
			}

			// passed auth
			next.ServeHTTP(w, r)
		})
	}
}
