package main

import (
	"net/http"
	"os"
)

func main() {
	http.Handle("/dir1/", http.StripPrefix("/dir1", http.FileServer(http.Dir("./public"))))
	http.Handle("/dir2/", http.StripPrefix("/dir2", http.FileServer(&fs{http.Dir("./public")})))
	http.ListenAndServe(":8080", nil)
}

type fs struct {
	http.FileSystem
}

func (fs *fs) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
