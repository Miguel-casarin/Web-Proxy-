// recebe todas as rquisições
package handler

import (
"net/http" // servidor HTTP 
"net/url" // manipulação e parsing de URLs
"strings"

"web-proxy/config"
"web-proxy/logger"
)

func proxy(conf *config.Config, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {

			http.Error(w, "CONNECT nao implementado aqui", 501)
			return
	}
}