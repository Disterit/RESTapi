package api

import (
	"RESTapi/internal/storage/postgres"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Transfer struct {
	ToID   int     `json:"ID"`
	Amount float32 `json:"amount"`
}

func CreateWallet(w http.ResponseWriter, req *http.Request) {
	const op = "Internal.api.createWallet()"
	wallet, err := postgres.CreateWallet()
	if err != nil {
		log.Printf("%s %v", op, err)
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(wallet)
	if err != nil {
		log.Printf("%s %v", op, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func SendMoney(w http.ResponseWriter, req *http.Request) {
	const op = "Internal.api.sendMoney()"

	jsonData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("%s %v", op, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var transfer Transfer

	err = json.Unmarshal(jsonData, &transfer)
	if err != nil {
		log.Printf("%s %v", op, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	walletID, err := strconv.Atoi(strings.Split(req.URL.String(), "/")[4])
	if err != nil {
		log.Printf("%s %v", op, err)
		return
	}

	status, err := postgres.TransferMoney(transfer.Amount, walletID, transfer.ToID)
	if err != nil {
		log.Printf("%s %v", op, err)
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)

}

func HistoryWallet(w http.ResponseWriter, req *http.Request) {

}

func InfoWallet(w http.ResponseWriter, req *http.Request) {

}
