package api

import (
	"RESTapi/internal/storage/postgres"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Transfer struct {
	ToID   int     `json:"ID"`
	Amount float32 `json:"amount"`
}

func CreateWalletHandler(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "Internal.api.createWallet()"

		wallet, err := storage.CreateWallet()
		if err != nil {
			log.Printf("%s %v", op, err)
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(wallet)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func SendMoneyHandler(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Internal.api.sendMoney()"

		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var transfer Transfer

		err = json.Unmarshal(jsonData, &transfer)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		walletID, err := strconv.Atoi(strings.Split(r.URL.String(), "/")[4])
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status, err := storage.TransferMoney(transfer.Amount, walletID, transfer.ToID)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(status)
			return
		}

		err = storage.AddHistoryWallet(time.Now(), walletID, transfer.ToID, transfer.Amount)
		if err != nil {
			log.Printf("%s %v", op, err)
			return
		}

		w.WriteHeader(status)
	}
}

func HistoryWalletHandler(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Internal.api.historyWallet()"

		walletID, err := strconv.Atoi(strings.Split(r.URL.String(), "/")[4])
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		walletHistory, status := storage.HistoryWallet(walletID)
		if status != http.StatusOK {
			log.Printf("%s %v", op, status)
			w.WriteHeader(status)
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(walletHistory)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func InfoWalletHandler(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Internal.api.InfoWallet()"

		walletID, err := strconv.Atoi(strings.Split(r.URL.String(), "/")[4])
		if err != nil {
			log.Printf("%s %v", op, err)
			return
		}

		wallet, err := storage.GetWalletInfo(walletID)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(wallet)
		if err != nil {
			log.Printf("%s %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
