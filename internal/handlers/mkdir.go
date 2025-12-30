package handlers

import (
	"net/http"
	"os"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
)

func MkdirHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reqPath := r.URL.Query().Get("path")
		if reqPath == "" {
			http.Error(w, "path required", http.StatusBadRequest)
			return
		}

		absPath, err := fs.Resolve(root, reqPath)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		// Ensure it does not already exist
		if _, err := os.Stat(absPath); err == nil {
			http.Error(w, "directory already exists", http.StatusConflict)
			return
		}

		// Create directory
		err = os.Mkdir(absPath, 0755)
		if err != nil {
			http.Error(w, "failed to create directory", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("directory created\n"))
	}
}
