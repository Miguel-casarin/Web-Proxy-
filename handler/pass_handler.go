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

	if err != nil {
		http.Error(w, "Erro ao acessar o destino", 502)
		return
	}

	defer response.Body.Close()

	// Copia headers da resposta de volta ao cliente
	for key, values := range response.Header {
		for _, v := range values { // no go não existe while, aqui e um for quer tu não que pegar o indice
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(response.StatusCode)

	// copia a pagina acessada
	io.Copy(w, response.Body)
}
