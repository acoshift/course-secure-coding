package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store")

	if r.URL.Path == "/upload.txt" {
		w.Header().Set("Content-Type", "text/plain")
		http.ServeFile(w, r, "./upload.txt")
		return
	}

	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<script src="/upload.txt"></script>
		<h1>Welcome</h1>
	`))
}
