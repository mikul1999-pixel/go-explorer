package fs

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrInvalidPath = errors.New("invalid path")

func Resolve(root, reqPath string) (string, error) {
	// Treat empty or "/" as root
	if reqPath == "" || reqPath == "/" {
		return filepath.Abs(root)
	}

	// Reject absolute paths
	if filepath.IsAbs(reqPath) {
		return "", ErrInvalidPath
	}

	// Clean user path
	clean := filepath.Clean(reqPath)

	// Disallow paths that resolve upward
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", ErrInvalidPath
	}

	// Join and clean again
	full := filepath.Join(root, clean)

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}

	// Enforce subdir containment
	rel, err := filepath.Rel(absRoot, absFull)
	if err != nil {
		return "", ErrInvalidPath
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrInvalidPath
	}

	return absFull, nil
}
