package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikul1999-pixel/go-explorer/internal/handlers"
)

func main() {
	rootDir := os.Getenv("ROOT_DIR")
	if rootDir == "" {
		log.Fatal("ROOT_DIR must be set")
	}

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":3030"
	}

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// API routes
	mux.HandleFunc("/fs/list", handlers.ListHandler(rootDir))
	mux.HandleFunc("/fs/download", handlers.DownloadHandler(rootDir))
	mux.HandleFunc("/fs/upload", handlers.UploadHandler(rootDir))
	mux.HandleFunc("/fs/mkdir", handlers.MkdirHandler(rootDir))
	mux.HandleFunc("/fs/rename", notImplemented)
	mux.HandleFunc("/fs/delete", notImplemented)

	server := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // allow large downloads
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting explorer API on %s (root=%s)", addr, rootDir)
	log.Fatal(server.ListenAndServe())
}


// placeholder response until implementation
func notImplemented(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
