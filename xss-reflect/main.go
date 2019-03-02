package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-XSS-Protection", "0")

	if r.URL.Path == "/report" {
		name := r.FormValue("name")
		fmt.Fprintln(w, name)
		fmt.Fprintln(w, "id,name")
		fmt.Fprintln(w, "1,hello")
		fmt.Fprintln(w, "2,world")
		return
	}

	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<h1>Welcome</h1>
		<a href="/report?name=users">Click to see Report</a>
	`))
}
