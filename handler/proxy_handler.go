package handler

import (
	"net/http" // servidor HTTP
	"net/url"  // manipulação e parsing de URLs
	"strings"

	"web-proxy/config"
	"web-proxy/logger"
)

func Proxy(conf *config.Config, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			http.Error(w, "CONNECT nao implementado aqui", 501)
			return
		}

		destinationURL := strings.TrimPrefix(r.URL.Path, "/")

		// separa a url de destino em protocolo, host e path
		target, err := url.Parse(destinationURL)
		if err != nil || target.Host == "" { // se der erro ou faltar pedaços
			http.Error(w, "URL invalida", 400)
			return
		}

		domain := target.Hostname()

		// dominio bloqueado
		if conf.IsBlocked(domain) {
			log.Write(destinationURL, "bloqueado")
			BlockHandler(w, r, domain)
			return
		}

		// Censura
		if len(conf.Words) > 0 {
			log.Write(destinationURL, "filtrado")
			FilterHandler(w, r, destinationURL, conf.Words)
			return
		}

		// repasse
		log.Write(destinationURL, "permitido")
		PassHandler(w, r, destinationURL)
	}
}
