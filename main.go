package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {

	//TLS config to only allow 1.2 or better
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	server := http.Server{
		Addr:         ":" + conf.SrvPort,
		Handler:      &myHandler{},
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 2000 * time.Millisecond,
		IdleTimeout:  30 * time.Second,
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = postHandlerWithToken

	//server.ListenAndServe()
	log.Println("Waiting for connections...")
	log.Fatal(server.ListenAndServeTLS(conf.CrtFile, conf.KeyFile))

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}

// stringContains checkes the srcString for any matches in the
// list, which is an array of strings.
func stringContains(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}
