package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db, _ = sql.Open("postgres", "postgres://localhost/postgres?sslmode=disable")

func main() {
	db.Exec(`
		create table if not exists test_sqli (
			key varchar primary key,
			value varchar not null default ''
		);
		insert into test_sqli (key, value) values
			('key1', 'value1'),
			('key2', 'value2')
		on conflict (key) do nothing;

		create table if not exists test_sqli2 (
			key varchar primary key,
			value varchar not null default ''
		);
		insert into test_sqli2 (key, value) values
			('smtp_user', 'root'),
			('smtp_pass', 'toor')
		on conflict (key) do nothing;
	`)

	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")

	query := fmt.Sprintf(`
		select value
		from test_sqli
		where key = '%s'
	`, key)

	var value string
	err := db.QueryRow(query).Scan(&value)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, value)
}
