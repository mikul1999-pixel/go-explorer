package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve(t *testing.T) {
	root := t.TempDir()

	// Create subdir for case
	err := os.Mkdir(filepath.Join(root, "valid"), 0755)
	if err != nil {
		t.Fatalf("failed to setup test dir: %v", err)
	}

	tests := []struct {
		name      string
		reqPath   string
		shouldErr bool
	}{
		{"root slash", "/", false},
		{"root empty", "", false},
		{"simple dir", "valid", false},
		{"nested dir", "valid/../valid", false},

		// Traversal attempts
		{"dot dot", "..", true},
		{"dot dot slash", "../", true},
		{"nested escape", "valid/../../", true},
		{"absolute escape", "/../", true},
		{"deep escape", "../../etc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Resolve(root, tt.reqPath)

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got path: %s", p)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Ensure result is inside root
			rel, err := filepath.Rel(root, p)
			if err != nil {
				t.Fatalf("rel error: %v", err)
			}

			if rel == ".." || filepath.HasPrefix(rel, "..") {
				t.Fatalf("resolved path escaped root: %s", p)
			}
		})
	}
}
