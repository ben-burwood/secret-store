package httpx

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// SPA serves static files from distDir, falling back to index.html for any
// path that does not resolve to a real file. Methods other than GET/HEAD return 405.
func SPA(distDir string) http.Handler {
	indexPath := filepath.Join(distDir, "index.html")
	absBase, _ := filepath.Abs(distDir)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clean := path.Clean(r.URL.Path)
		rel := strings.TrimPrefix(clean, "/")
		if rel == "" || rel == "." {
			http.ServeFile(w, r, indexPath)
			return
		}
		full := filepath.Join(distDir, filepath.FromSlash(rel))
		absFull, err := filepath.Abs(full)
		if err != nil || !strings.HasPrefix(absFull, absBase) {
			http.NotFound(w, r)
			return
		}
		info, err := os.Stat(full)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, indexPath)
			return
		}
		http.ServeFile(w, r, full)
	})
}
