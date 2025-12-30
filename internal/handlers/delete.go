package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
)

func DeleteHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reqPath := r.URL.Query().Get("path")
		if reqPath == "" || reqPath == "/" {
			http.Error(w, "refusing to delete root", http.StatusBadRequest)
			return
		}

		absPath, err := fs.Resolve(root, reqPath)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		// never allow deleting root
		absRoot, _ := filepath.Abs(root)
		if absPath == absRoot {
			http.Error(w, "refusing to delete root", http.StatusBadRequest)
			return
		}

		// Ensure it exists
		if _, err := os.Stat(absPath); err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		// Remove file or directory recursively
		if err := os.RemoveAll(absPath); err != nil {
			http.Error(w, "failed to delete", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted\n"))
	}
}
