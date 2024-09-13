package api

import (
	"RESTapi/internal/storage/postgres"
	"encoding/json"
	"log"
	"net/http"
)

func CreateWallet(w http.ResponseWriter, req *http.Request) {
	const op = "Internal.api.createWallet()"
	wallet, err := postgres.CreateWallet()
	if err != nil {
		log.Printf("%s %v", op, err)
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(wallet)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func SendMoney(w http.ResponseWriter, req *http.Request) {

}

func HistoryWallet(w http.ResponseWriter, req *http.Request) {

}

func InfoWallet(w http.ResponseWriter, req *http.Request) {

}
