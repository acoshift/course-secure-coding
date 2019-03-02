package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	go func() {
		// internal service
		http.ListenAndServe(":9000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("secret data"))
		}))
	}()

	http.HandleFunc("/proxy/", proxy)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	target, _ := url.Parse("http://" + strings.TrimPrefix(r.RequestURI, "/proxy/"))
	if target == nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	nr := *r
	nr.Host = target.Host
	nr.URL = target
	resp, err := http.DefaultTransport.RoundTrip(&nr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	io.Copy(ioutil.Discard, resp.Body)
}

func index(w http.ResponseWriter, r *http.Request) {
	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<div id="app">
			<input v-model="name">
			<p>Hello, <span v-text="name"></span>.</p>
		</div>
		<script src="/proxy/cdnjs.cloudflare.com/ajax/libs/vue/2.6.7/vue.min.js"></script>
		<script>
			new Vue({
				data () {
					return {
						name: 'World'
					}
				}
			}).$mount('#app')
		</script>
	`))
}
