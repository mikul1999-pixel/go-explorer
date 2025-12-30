package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
	"github.com/mikul1999-pixel/go-explorer/internal/model"
)

type ListResponse struct {
	Path    string        `json:"path"`
	Entries []model.Entry `json:"entries"`
}

func ListHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Query().Get("path")

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

		if !info.IsDir() {
			http.Error(w, "not a directory", http.StatusBadRequest)
			return
		}

		entries, err := os.ReadDir(absPath)
		if err != nil {
			http.Error(w, "failed to read directory", http.StatusInternalServerError)
			return
		}

		resp := ListResponse{
			Path:    filepath.Clean("/" + reqPath),
			Entries: make([]model.Entry, 0, len(entries)),
		}

		for _, e := range entries {
			fi, err := e.Info()
			if err != nil {
				continue // skip unreadable entries
			}

			entry := model.Entry{
				Name:     e.Name(),
				Modified: fi.ModTime(),
			}

			if fi.IsDir() {
				entry.Type = "dir"
				entry.Size = 0
			} else {
				entry.Type = "file"
				entry.Size = fi.Size()
			}

			resp.Entries = append(resp.Entries, entry)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
