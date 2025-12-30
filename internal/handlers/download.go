package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
)

func DownloadHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		info, err := os.Stat(absPath)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if info.IsDir() {
			http.Error(w, "cannot download directory", http.StatusBadRequest)
			return
		}

		file, err := os.Open(absPath)
		if err != nil {
			http.Error(w, "failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Headers
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set(
			"Content-Disposition",
			`attachment; filename="`+filepath.Base(absPath)+`"`,
		)
		w.Header().Set("Content-Length", string(info.Size()))

		// Stream
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	}
}
