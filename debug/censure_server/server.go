package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Caminho para o html
	caminhoArquivo := "site_palavores/palavrao.html"

	// A função ServeFile lê o arquivo do disco e já envia com o Content-Type correto (text/html)
	http.ServeFile(w, r, caminhoArquivo)
}

func main() {
	porta := ":8080"
	
	// Configura a rota principal ("/")
	http.HandleFunc("/", handler)

	fmt.Printf("Servidor HTTP puro rodando.\nAcesse: http://localhost%s\n", porta)
	fmt.Println("Lendo o HTML do diretório: site_palavores/palavrao.html")
	
	// Inicia o servidor HTTP
	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}