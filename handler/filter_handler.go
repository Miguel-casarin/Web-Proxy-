package handler

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

func FilterHandler(w http.ResponseWriter, r *http.Request, targetURL string, words map[string]string) {

	req, err := http.NewRequest(r.Method, targetURL, r.Body)

	// ps nil e a mesma coisa que null, e err e variavel de convernção para guardar erros
	// aqui eu basicamente verifico se houve um erro
	// a gente não tem algo como o try/except
	if err != nil {
		http.Error(w, "Erro ao criar requisicao", 500)
		return
	}

	for key, values := range r.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		http.Error(w, "Erro ao acessar o destino", 502)
		return
	}

	defer response.Body.Close()

	for key, values := range response.Header {
		if strings.EqualFold(key, "content-length") {
			continue
		}
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	contentType := response.Header.Get("Content-Type")

	// verifica se e um html
	if !strings.Contains(contentType, "text/html") {
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)
		return
	}

	body, err := io.ReadAll(response.Body) // carrega o html na memoria
	if err != nil {
		http.Error(w, "Erro ao ler resposta", 500)
		return
	}

	html := string(body)

	// substitui as palavras
	for word, replacement := range words {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		html = re.ReplaceAllString(html, replacement)
	}

	w.WriteHeader(response.StatusCode)
	w.Write([]byte(html))
}
