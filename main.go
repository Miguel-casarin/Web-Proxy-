package main

import (
	"fmt"
	"log"
	"net/http"

	"web-proxy/config"
	"web-proxy/handler"
	"web-proxy/logger"
	"web-proxy/tunnel"
)

const tunnelEnabled = true

func main() {

	conf := config.LoadFiles("config/blocked.json", "config/words.json")

	lg := logger.New("log.json")

	// Roteador principal
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Proxy(conf, lg))

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			if tunnelEnabled {
				tunnel.Handle(w, r, lg)
			} else {
				http.Error(w, "CONNECT nao suportado", 501)
			}
			return
		}
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Proxy rodando em http://localhost:5000")
	log.Fatal(server.ListenAndServe())
}
