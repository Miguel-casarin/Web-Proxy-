package tunnel

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"web-proxy/logger"
)

func Handle(w http.ResponseWriter, r *http.Request, log *logger.Logger) {

	// abre conexão TCP com o servidor de destino
	destConection, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, "Nao foi possivel conectar ao destino", http.StatusBadGateway)
		return
	}

	// fecha a conexão
	defer destConection.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "\r\n")

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("ERRO: servidor nao suporta Hijack")
		// salvo nos meus logs que o servidor que tentei acessar nao suporta hijacker
		log.Write(r.Host, "erro-hijack")
		http.Error(w, "Servidor nao suporta tunnel", http.StatusInternalServerError)
		return
	}

	// pega a conexão tcp do cliente
	clientConection, _, err := hijacker.Hijack()
	if err != nil {
		return
	}

	defer clientConection.Close()

	// buffer para as duas goroutines
	done := make(chan struct{}, 2)

	go func() {
		io.Copy(destConection, clientConection)
		// avisa que esse lado terminou
		done <- struct{}{}
	}()

	go func() {
		io.Copy(clientConection, destConection)
		// avisa que esse lado terminou
		done <- struct{}{}
	}()

	<-done
}
