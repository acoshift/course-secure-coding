package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/transfer", transfer)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	{
		s, _ := r.Cookie("sess")
		if s != nil && s.Value == "1234" {
			isLogin = true
		}
	}

	if !isLogin {
		// language=HTML
		w.Write([]byte(`
			<!doctype html>
			<h1>Welcome to the most secure bank!</h1>
			<a href="/login">Click here to Login</a>
		`))
		return
	}

	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<h1>Welcome to the most secure bank!</h1>
		<a href="/logout">Click here to Logout</a><br>
		<h3>Transfer</h3>
		<form method="POST" action="/transfer">
			<label>Amount</label>
			<input name="amount" type="number">
			<br>
			<button>Sent my money!</button>
		</form>
	`))
}

func login(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    "1234",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	amount := r.PostFormValue("amount")
	fmt.Println("transfer money:", amount)

	http.Redirect(w, r, "/", http.StatusFound)
}
