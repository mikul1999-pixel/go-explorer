package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikul1999-pixel/go-explorer/internal/fs"
)

const maxUploadSize = 100 << 20 // 100 MB (Cloudflare limit)

func UploadHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		targetDir := r.URL.Query().Get("path")

		absDir, err := fs.Resolve(root, targetDir)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		info, err := os.Stat(absDir)
		if err != nil || !info.IsDir() {
			http.Error(w, "target must be an existing directory", http.StatusBadRequest)
			return
		}

		// Parse multipart
		err = r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			http.Error(w, "invalid multipart form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file field required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := filepath.Base(header.Filename)
		finalPath := filepath.Join(absDir, filename)

		// Do not overwrite
		if _, err := os.Stat(finalPath); err == nil {
			http.Error(w, "file already exists", http.StatusConflict)
			return
		}

		tmpPath := finalPath + ".tmp"

		tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
		if err != nil {
			http.Error(w, "failed to create temp file", http.StatusInternalServerError)
			return
		}

		defer func() {
			tmpFile.Close()
			os.Remove(tmpPath) // cleanup on failure
		}()

		_, err = io.Copy(tmpFile, file)
		if err != nil {
			http.Error(w, "failed to write file", http.StatusInternalServerError)
			return
		}

		err = tmpFile.Close()
		if err != nil {
			http.Error(w, "failed to finalize file", http.StatusInternalServerError)
			return
		}

		err = os.Rename(tmpPath, finalPath)
		if err != nil {
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("upload successful\n"))
	}
}
