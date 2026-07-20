package utils

import (
	"os"
	"path/filepath"
)

// FindTemplatesDir locates the HTML templates directory.
func FindTemplatesDir() string {
	candidates := []string{
		"templates",
		filepath.Join("..", "templates"),
		filepath.Join("..", "..", "templates"),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "templates"))
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "templates"
}
