package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}

	http.HandleFunc("/cert", hello)
	address := ":8443"
	println("starting server on address" + address)

	caCert, err := ioutil.ReadFile(os.Getenv("CA_CERT_PEM_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS(os.Getenv("CLIENT_CERT_PEM_PATH"), os.Getenv("CLIENT_CERT_KEY_PATH"))
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Printf("%v\n", r.Header)
	fmt.Printf("%v\n", r.TLS)
	fmt.Printf("%v\n", r.Trailer)
	println(string(b))

	j := `{"status": "OK"}`
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(j))

	return
}
