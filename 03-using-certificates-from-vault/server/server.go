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
	// PathToCAFile is the path in the filesystem from where to read the root
	// certificate authority's certificate
	PathToCAFile = "/secret/chain.pem"
	// PathToCertFile is the path in the filesystem from where to read the
	// certificate file for 'time-server'
	PathToCertFile = "/secret/crt.pem"
	// PathToPrivateKeyFile is the path in the filesystem from where to read the
	// private key file for 'time-server'
	PathToPrivateKeyFile = "/secret/key.pem"
)

var (
	port int
)

// getTime returns a string containing the current time in RFC3339 format
func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.IntVar(&port, "port", 443, "the port on which to listen")
	fs.Parse(os.Args[1:])

	// read the certificate authority bundle from PathToCAFile
	caBundle, err := ioutil.ReadFile(PathToCAFile)
	if err != nil {
		log.Fatalf("failed to read ca bundle: %v", err)
	}

	// create a new x509 CertPool where to store the CA bundle
	// certificates in this bundle will be trusted as root CAs
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(caBundle); !ok {
		log.Fatalf("failed to read ca bundle: %v", err)
	}
	// create a new tls.Config struct containing the desired TLS configuration:
	// - require and verify client certificates
	// - treat the certificates signed by CAs in the x509 CertPool as trusted
	// - require TLSv1.2
	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  roots,
		MinVersion: tls.VersionTLS12,
	}

	// create a new http server that will listen on the specified port and use
	// the TLS configuration created above
	srv := http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", port),
		TLSConfig: tlsConfig,
	}

	// add a handler to GET /
	// this is only executed after successful verification of the client's
	// certificate
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// grab the 'common name' from the client's certificate
		name := r.TLS.PeerCertificates[0].Subject.CommonName
		log.Printf("got time request from '%s'", name)
		// respond with the time
		w.Write([]byte(getTime()))
	})

	// listen to and serve requests until an error occurs using the contents of
	// PathToCertFile and PathToPrivateKeyFile to establish TLS on the server
	// side
	log.Fatal(srv.ListenAndServeTLS(PathToCertFile, PathToPrivateKeyFile))
}
