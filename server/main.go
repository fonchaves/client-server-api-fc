package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeData struct {
	Bid string `json:"bid"`
}

type Exchange struct {
	USDBRL ExchangeData `json:"USDBRL"`
}

type ExchangeDb struct {
	ID  int `gorm:"primaryKey"`
	Bid string
	gorm.Model
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
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
	error = saveToDatabase(exchange)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if error == context.DeadlineExceeded {
		println("Timeout")
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	// Write response
	body, error := json.Marshal(exchange.USDBRL)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getExchange() (*Exchange, error) {
	// add timeout after 200 miliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Create request
	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	// req, error := http.NewRequestWithContext(ctx, "GET", "http://localhost:3001/USD-BRL", nil)
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

func saveToDatabase(exchange *Exchange) error {
	// add timeout after 10 miliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Open database
	db, err := gorm.Open(sqlite.Open("exchangeDb.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	// db.AutoMigrate(&ExchangeDb{})

	// Create
	db.Create(&ExchangeDb{
		Bid: exchange.USDBRL.Bid,
	})

	return ctx.Err()
}
