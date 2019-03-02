package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/messages", messages)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func messages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// language=JSON
	w.Write([]byte(`
		[
			{ "content": "Hello", "type": "text" },
			{ "content": "\" onerror=\"localStorage.getItem('token') && alert('sending token ' + localStorage.getItem('token') + ' to hacker server')", "type": "image" }
		]
	`))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// language=JSON
	w.Write([]byte(`
		{
			"token": "tk1234"
		}
	`))
}
