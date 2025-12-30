package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
)

func RenameHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if from == "" || to == "" {
			http.Error(w, "from and to are required", http.StatusBadRequest)
			return
		}

		if from == "/" || to == "/" {
			http.Error(w, "refusing to rename root", http.StatusBadRequest)
			return
		}

		absFrom, err := fs.Resolve(root, from)
		if err != nil {
			http.Error(w, "invalid source path", http.StatusBadRequest)
			return
		}

		absTo, err := fs.Resolve(root, to)
		if err != nil {
			http.Error(w, "invalid destination path", http.StatusBadRequest)
			return
		}

		absRoot, _ := filepath.Abs(root)
		if absFrom == absRoot || absTo == absRoot {
			http.Error(w, "refusing to rename root", http.StatusBadRequest)
			return
		}

		// Source must exist
		if _, err := os.Stat(absFrom); err != nil {
			http.Error(w, "source not found", http.StatusNotFound)
			return
		}

		// Destination must not exist
		if _, err := os.Stat(absTo); err == nil {
			http.Error(w, "destination already exists", http.StatusConflict)
			return
		}

		// Parent of destination must exist
		parent := filepath.Dir(absTo)
		if info, err := os.Stat(parent); err != nil || !info.IsDir() {
			http.Error(w, "destination parent does not exist", http.StatusBadRequest)
			return
		}

		if err := os.Rename(absFrom, absTo); err != nil {
			http.Error(w, "rename failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("renamed\n"))
	}
}
