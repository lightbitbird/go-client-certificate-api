package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}
	var cert tls.Certificate
	if len(os.Args) == 2 && os.Args[1] == "pfx" {
		cert = *getCertificateFromPkcs12()
	} else {
		cert = *getCertificate()
	}

	// load CA cert
	caCert, err := ioutil.ReadFile(os.Getenv("CA_CERT_PEM_PATH"))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:            caCertPool,
		//InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// https client request
	url := "https://localhost:8443/cert"
	j := []byte(`{"id": "987654", "name": "golang"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: transport}

	// read response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	log.Println(string(contents))
}

func getCertificate() *tls.Certificate  {
	// load client cert
	cert, err := tls.LoadX509KeyPair(os.Getenv("CLIENT_CERT_PEM_PATH"), os.Getenv("CLIENT_CERT_KEY_PATH"))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	return &cert
}

func getCertificateFromPkcs12() *tls.Certificate {
	data, err := ioutil.ReadFile(os.Getenv("CLIENT_CERT_PKCS12_PATH"))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	blocks, err := pkcs12.ToPEM(data, os.Getenv("CLIENT_CERT_PWD"))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	var pemData []byte
	log.Printf("decoded %v blocks\n", len(blocks))
	for i, b := range blocks {
		log.Printf("block %v type %s\n", i, b.Type)
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	return &cert
}