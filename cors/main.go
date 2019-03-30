package main

import (
	"net/http"
)

func main() {
	// start web page on another port
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// language=HTML
			w.Write([]byte(`
				<!doctype html>
				<button onclick="invokeApi('/with-cors')">Fetch API with CORS</button>
				<button onclick="invokeApi('/no-cors')">Fetch API without CORS</button>
				<div id=result></div>
				<script>
					function invokeApi (path) {
						const result = document.querySelector('#result')
						result.innerHTML = ''

						fetch('http://localhost:3333' + path, {
							method: 'POST',
							headers: new Headers({
								'Content-Type': 'application/json'
							}),
							body: JSON.stringify({})
						})
							.then((resp) => {
								result.innerHTML += 'X-Request-Id: ' + resp.headers.get('X-Request-Id') + '<br>'
								return resp.text()
							})
							.then((res) => {
								result.innerHTML += res
							})
							.catch((err) => {
								result.innerHTML += err
							})
					}
				</script>
			`))
		})

		http.ListenAndServe(":8080", mux)
	}()

	mux := http.NewServeMux()
	mux.Handle("/with-cors", cors(http.HandlerFunc(result)))
	mux.Handle("/no-cors", http.HandlerFunc(result))

	http.ListenAndServe(":3333", mux)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow only http://localhost:8080
		if r.Header.Get("Origin") != "http://localhost:8080" {
			http.Error(w, "Forbidden - Origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			// pre-flight request

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Max-Age", "7200")
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			w.WriteHeader(http.StatusOK) // or http.StatusNoContent
			return
		}

		// toggle Access-Control-Expose-Headers to see result in browser
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-Id")

		h.ServeHTTP(w, r)
	})
}

func result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Request-Id", "1234")
	w.Write([]byte(`{"name":"launcher-1234"}`))
}
