package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikul1999-pixel/go-explorer/internal/auth"
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

	// Load auth config
	authCfg := auth.LoadConfig()

	// Root mux
	mux := http.NewServeMux()

	// Public route
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// FS mux (protected)
	fsMux := http.NewServeMux()
	fsMux.HandleFunc("/fs/list", handlers.ListHandler(rootDir))
	fsMux.HandleFunc("/fs/download", handlers.DownloadHandler(rootDir))
	fsMux.HandleFunc("/fs/upload", handlers.UploadHandler(rootDir))
	fsMux.HandleFunc("/fs/mkdir", handlers.MkdirHandler(rootDir))
	fsMux.HandleFunc("/fs/rename", handlers.RenameHandler(rootDir))
	fsMux.HandleFunc("/fs/delete", handlers.DeleteHandler(rootDir))

	// Wrap FS routes with auth
	protectedFS := auth.Middleware(authCfg)(fsMux)
	mux.Handle("/fs/", protectedFS)

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
