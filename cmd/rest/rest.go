package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tagptroll1/receipt-processor/receipt"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("POST /v1/receipts/process", receipt.ProcessReceipts)
	handler.HandleFunc("GET /v1/receipts/{id}/points", receipt.GetPoints)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
