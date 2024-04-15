package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type CORSConfig struct {
	AllowedOrigins []string 
	AllowedMethods []string 
	AllowedHeaders []string  
	AllowCredentials bool    
}

func Logging() Middleware {
    return func(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            defer func() { log.Println(r.Method, r.URL.Path, time.Since(start)) }()

            f(w, r)
        }
    }
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}

func CORS(config CORSConfig) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			allowedOrigin := ""
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowedOrigin = o
					break
				}
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.Header().Set("Access-Control-Allow-Methods", stringJoin(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", stringJoin(config.AllowedHeaders, ", "))
			}

			if r.Method == http.MethodOptions {
				if allowedOrigin != "" {
					w.WriteHeader(http.StatusOK)
					return
				}
			} else {
				f(w, r)
			}
		}
	}
}

// Helper function to concatenate slice of strings with a delimiter
func stringJoin(items []string, delim string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for _, item := range items[1:] {
		result += delim + item
	}
	return result
}