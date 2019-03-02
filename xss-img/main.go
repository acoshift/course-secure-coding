package main

import (
	"net/http"
	"text/template"
)

func main() {
	// curl "http://localhost:8080/?image=j%26%23X41vascript%3Aalert%28%27hello+%3AD%27%29" -XPOST
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

var imageList []string

// language=HTML
var tmpl = template.Must(template.New("").Parse(`
<!doctype html>
<h1>Image List</h1>
{{range .}}
<a href="{{.}}"><img src="{{.}}"></a>
{{end}}
`))

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		imageURL := r.FormValue("image")
		imageList = append(imageList, imageURL)
		w.Write([]byte("OK"))
		return
	}

	tmpl.Execute(w, imageList)
}
