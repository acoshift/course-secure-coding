package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/noprg", noprg)
	http.HandleFunc("/prg", prg)
	http.HandleFunc("/prg/success", prgSuccess)
	http.ListenAndServe(":8080", nil)
}

func noprg(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("checkout success!")
		// language=HTML
		w.Write([]byte(`
			<!doctype html>
			<h1>Checkout Success !</h1>
		`))
		return
	}

	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<h1>Shopping Cart</h1>
		<form method="POST">
			<button>Checkout</button>
		</form>
	`))
}

func prg(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("checkout success!")
		http.Redirect(w, r, "/prg/success", http.StatusFound)
		return
	}

	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<h1>Shopping Cart</h1>
		<form method="POST">
			<button>Checkout</button>
		</form>
	`))
}

func prgSuccess(w http.ResponseWriter, r *http.Request) {
	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<h1>Checkout Success !</h1>
	`))
}
