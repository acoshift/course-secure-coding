package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func loadPem(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	block, _ := pem.Decode(b)
	return block.Bytes
}

func main() {
	caPriv, _ := x509.ParsePKCS1PrivateKey(loadPem("ca.key"))
	caCert, _ := x509.ParseCertificate(loadPem("ca.crt"))

	mu := sync.Mutex{}
	certs := make(map[string]*tls.Certificate)

	srv := http.Server{
		Addr:    ":9443",
		Handler: http.HandlerFunc(roundTripHTTPS),
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				mu.Lock()
				defer mu.Unlock()

				if cert := certs[info.ServerName]; cert != nil {
					return cert, nil
				}

				key, _ := rsa.GenerateKey(rand.Reader, 2048)
				serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
				now := time.Now()
				certBytes, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
					Subject: pkix.Name{
						CommonName: info.ServerName,
					},
					Issuer:       caCert.Subject,
					SerialNumber: serial,
					NotBefore:    now.UTC(),
					NotAfter:     now.Add(30 * 24 * time.Hour).UTC(),
					KeyUsage:     x509.KeyUsageDigitalSignature,
					DNSNames:     []string{info.ServerName},
					ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				}, caCert, &key.PublicKey, caPriv)
				if err != nil {
					return nil, err
				}

				cert := &tls.Certificate{
					Certificate: [][]byte{certBytes},
					PrivateKey:  key,
					Leaf:        caCert,
				}
				certs[info.ServerName] = cert
				return cert, nil
			},
		},
	}
	go srv.ListenAndServeTLS("", "")

	http.ListenAndServe(":9000", http.HandlerFunc(proxy))
}

var tr = &http.Transport{
	MaxConnsPerHost: 10000,
}

func proxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		upstream, err := net.Dial("tcp", "127.0.0.1:9443")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer upstream.Close()

		downstream, wr, err := w.(http.Hijacker).Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer downstream.Close()

		wr.WriteString("HTTP/1.1 200 OK\n\n")
		wr.Flush()

		go io.Copy(upstream, downstream)
		io.Copy(downstream, upstream)
		return
	}

	roundTrip(w, r)
}

func roundTripHTTPS(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "https"
	roundTrip(w, r)
}

func roundTrip(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = r.Host
	r.Header.Del("Accept-Encoding")
	resp, err := tr.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Println(resp.Proto, resp.Status)
	resp.Header.Write(os.Stdout)
	fmt.Println()
	io.Copy(w, io.TeeReader(resp.Body, os.Stdout))
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	fmt.Println()
}
