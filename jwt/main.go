package main

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// openssl genrsa -out key.rsa 1024
// openssl rsa -in key.rsa -pubout > key.pub

var (
	privKey []byte
	pubKey  []byte
)

func loadPem(filename string) []byte {
	f, _ := ioutil.ReadFile(filename)
	block, _ := pem.Decode(f)
	return block.Bytes
}

func main() {
	// load keys
	privKey = loadPem("key.rsa")
	pubKey = loadPem("key.pub")

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", auth)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/user", user)
	mux.HandleFunc("/publickey", publickey)

	http.ListenAndServe(":8080", mux)
}

func auth(w http.ResponseWriter, r *http.Request) {
	user := r.PostFormValue("user")
	pass := r.PostFormValue("password")

	if user != "user" || pass != "1234" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := Token{
		Subject: user,
		Admin:   false,
	}
	tk, err := token.Sign("RS256", privKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "token: %s\n", tk)
}

func parseToken(r *http.Request) (*Token, error) {
	tk := r.Header.Get("Authorization")
	tk = strings.TrimPrefix(tk, "Bearer ")

	return ParseToken(tk, pubKey)
}

func admin(w http.ResponseWriter, r *http.Request) {
	token, err := parseToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if !token.Admin {
		http.Error(w, "You are not admin!", http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Welcome, admin!\n")
}

func user(w http.ResponseWriter, r *http.Request) {
	token, err := parseToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Hello, %s!\n", token.Subject)
}

func publickey(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "key.pub")
}
