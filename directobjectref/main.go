package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", users)

	http.ListenAndServe(":8080", mux)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Secret   bool   `json:"-"`
}

var userList = []*User{
	{1, "admin", "Admin", true},
	{2, "staff", "Staff", true},
	{5, "user-1", "User 1", false},
	{6, "user-2", "User 2", false},
}

func users(w http.ResponseWriter, r *http.Request) {
	rawID := r.FormValue("id")

	// list
	if rawID == "" {
		list := make([]*User, 0)
		for _, user := range userList {
			if !user.Secret {
				list = append(list, user)
			}
		}
		writeJSON(w, list)
		return
	}

	id, _ := strconv.ParseInt(rawID, 10, 64)
	if id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// broken access control
	for _, user := range userList {
		if user.ID == id {
			writeJSON(w, user)
			return
		}
	}

	http.NotFound(w, r)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}
