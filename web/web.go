package web

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func OpenServer(a <-chan time.Time) {

	//ReadHeaderTimeout can be used in place of ReadTimeout, but for audit purposes RT's preferred
	mainserver := myServer{
		serve: &http.Server{
			Addr:         "localhost:8080",
			ReadTimeout:  4 * time.Second,
			WriteTimeout: 8 * time.Second,
			IdleTimeout:  10 * time.Second,
			Handler:      nil,
			TLSConfig: &tls.Config{
				MinVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
			},
		},
	}

	mainserver.Routes(a)

	log.Fatal(mainserver.serve.ListenAndServeTLS("tls/cert.pem", "tls/key.pem"))

}
