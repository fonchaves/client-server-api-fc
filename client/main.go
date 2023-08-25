package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

	//O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
	// Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

	// O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
