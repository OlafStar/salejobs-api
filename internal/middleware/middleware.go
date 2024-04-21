package middleware

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type CORSConfig struct {
	AllowedOrigins []string 
	AllowedMethods []string 
	AllowedHeaders []string  
	AllowCredentials bool    
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gzw gzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.Writer.Write(data)
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

func RateLimit(r rate.Limit, b int) Middleware {
	limiter := rate.NewLimiter(r, b)
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			f(w, r)
		}
	}
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

//This middlewate will increese runtime but reduce bandwidth (tested on advertisements endpoint)
//Requests/sec: 110000 -> 8000
//Bandwidth: 5.69GB -> 3.97GB
func GzipMiddleware() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				f(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				http.Error(w, "Could not create gzip writer", http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			gzw := gzipResponseWriter{ResponseWriter: w, Writer: gz}

			f(gzw, r)

			gz.Flush()
		}
	}
}

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