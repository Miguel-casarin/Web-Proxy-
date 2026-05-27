package handler

import (
	"html/template"
	"net/http"
)

// BlockData e passado para o template HTML.
type BlockData struct {
	Domain string
}

// pagina bloqueada retorna o template
func BlockHandler(w http.ResponseWriter, r *http.Request, domain string) {

	template, err := template.ParseFiles("templates/blocked.html")

	if err != nil {
		// Fallback: mensagem simples se o template nao for encontrado
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("<h1>Acesso bloqueado</h1>"))
		return
	}

	w.WriteHeader(http.StatusForbidden)
	template.Execute(w, BlockData{Domain: domain})
}
