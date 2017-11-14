package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PathToCAFile         = "/secret/chain.pem"
	PathToCertFile       = "/secret/crt.pem"
	PathToPrivateKeyFile = "/secret/key.pem"
)

var (
	port int
)

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 443, "the port on which to listen")
	fs.Parse(os.Args[1:])

	caBundle, err := ioutil.ReadFile(PathToCAFile)
	if err != nil {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(caBundle); !ok {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  roots,
		MinVersion: tls.VersionTLS12,
	}

	srv := http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", port),
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.TLS.PeerCertificates[0].Subject.CommonName
		log.Printf("got time request from '%s'", name)
		w.Write([]byte(getTime()))
	})

	log.Fatal(srv.ListenAndServeTLS(PathToCertFile, PathToPrivateKeyFile))
}
