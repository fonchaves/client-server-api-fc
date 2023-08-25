package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ExchangeData struct {
	Bid string `json:"bid"`
}

func main() {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	// Do the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		panic(err)
	}
	if err == context.DeadlineExceeded {
		println("Timeout")
		return
	}
	defer res.Body.Close()

	// Read all contents
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	// Parse JSON
	var data ExchangeData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	// Create to file
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}
	defer file.Close()

	// Write to file
	_, err = file.WriteString(fmt.Sprintf("Dolar: %s", data.Bid))
	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Dolar: ", data.Bid)

	io.Copy(os.Stdout, res.Body)
}
