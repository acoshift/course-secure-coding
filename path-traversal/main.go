package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// curl http://localhost:8080/../main.go --path-as-is
	http.ListenAndServe(":8080", http.HandlerFunc(serveFile))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.Path == "/" {
		path = "/index.html"
	}
	fn := filepath.Join("./public", path)
	fs, err := os.Open(fn)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer fs.Close()
	io.Copy(w, fs)
}
