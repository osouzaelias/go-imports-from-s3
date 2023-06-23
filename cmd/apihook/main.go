package main

import (
	"encoding/json"
	"fmt"
	"go-import-from-s3/internal/webhook"
	"log"
	"net/http"
	"time"
)

func main() {
	// Define a rota para a API
	http.HandleFunc("/endpoint", handleRequest)

	// Inicia o servidor na porta 8080
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Verifica o método da requisição
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o JSON do corpo da requisição para a struct Payload
	var payload webhook.Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	// Simula algum processamento
	time.Sleep(2 * time.Second)

	// Cria uma resposta de sucesso
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Requisição bem-sucedida!",
	}

	// Converte a resposta para JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao converter resposta para JSON", http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho da resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Escreve a resposta
	w.Write(responseJSON)
}
