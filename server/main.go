package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type ExchangeData struct {
	Bid string `json:"bid"`
}

type Exchange struct {
	USDBRL ExchangeData `json:"USDBRL"`
}

func main() {

	// O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// O server.go deverá consumir a API contendo o câmbio de Dólar e Real no
	// endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida
	// deverá retornar no formato JSON o resultado para o cliente.

	// Usando o package "context", o server.go deverá registrar no banco de dados
	// SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API
	// de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir
	// persistir os dados no banco deverá ser de 10ms.

	// Check route
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get exchange
	exchange, error := getExchange()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if error == context.DeadlineExceeded {
		println("Timeout")
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	// Save in database
	// TODO

	// Write response
	body, error := json.Marshal(exchange.USDBRL)
	if error != nil {
		panic(error)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func getExchange() (*Exchange, error) {
	// add timeout after 200 miliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	req, error := http.NewRequestWithContext(ctx, "GET", "http://localhost:3001/USD-BRL", nil)
	if error != nil {
		return nil, error
	}

	// Do the request
	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	// Read all contents
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	println(string(body))

	// Unmarshal
	var exchange Exchange
	error = json.Unmarshal(body, &exchange)
	if error != nil {
		return nil, error
	}

	return &exchange, nil
}
