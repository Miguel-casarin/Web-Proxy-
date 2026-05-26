package handler

import (
	"io"
	"net/http"
)

func PassHandler(w http.ResponseWriter, r *http.Request, destinationURLURL string) {
	
	// Cria e guarda a requisição
	req, err := http.NewRequest(r.Method, destinationURLURL, r.Body)

	if err != nil {
		http.Error(w, "Erro ao criar requisicao", 500)
		return
	}

	// copia os headers do cliente 
	for key, values := range r.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	// executa a requisição
	client := &http.Client{}
	response, err := client.Do(req)

	defer response.Body.Close()

	// Copia headers da resposta de volta ao cliente
	for key, values := range response.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(response.StatusCode)

	// copia a pagina acessada
	io.Copy(w, resp.Body)
}	