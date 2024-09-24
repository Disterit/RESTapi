package api

import (
	"RESTapi/internal/storage/postgres"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestCreateWalletHandler(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := postgres.NewStorage(db)

	mock.ExpectQuery("insert into wallet \\(balance\\) values \\(100\\) returning id, balance").
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 100.0))

	r := chi.NewRouter()

	r.Post("/wallet", CreateWalletHandler(storage))

	req, err := http.NewRequest("POST", "/wallet", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}

	var wallet postgres.Wallet
	err = json.Unmarshal(rr.Body.Bytes(), &wallet)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if wallet.ID != 1 || wallet.Balance != 100.0 {
		t.Errorf("unexpected wallet data: got %+v, want {ID: 1, Balance: 100.0}", wallet)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSendMoneyHandler(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := postgres.NewStorage(db)

	transfer := Transfer{
		Amount: 50.0,
		ToID:   2,
	}

	transferData, err := json.Marshal(transfer)
	if err != nil {
		t.Fatalf("could not marshal transfer data: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE wallet SET balance = balance - \$1 WHERE id = \$2 AND balance >= \$1`).
		WithArgs(transfer.Amount, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE wallet SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(transfer.Amount, transfer.ToID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`insert into history \(fromID, toID, time, amount\)`).
		WithArgs(1, transfer.ToID, sqlmock.AnyArg(), transfer.Amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := chi.NewRouter()
	r.Post("/check/v1/wallet/{id}/send", SendMoneyHandler(storage))

	// Create a request
	req, err := http.NewRequest("POST", "/check/v1/wallet/1/send", bytes.NewBuffer(transferData))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}
}

func TestHistoryWalletHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := postgres.NewStorage(db)

	// Задать фиксированное время
	fixedTime := time.Date(2024, 9, 24, 14, 37, 29, 0, time.UTC)

	walletHistory := []postgres.WalletHistory{
		{From: 1, To: 2, Time: fixedTime, Amount: 50.0},
		{From: 1, To: 3, Time: fixedTime, Amount: 30.0},
	}

	// Мокируем ожидание вызова функции HistoryWallet
	mock.ExpectQuery(`SELECT fromid, toid, time, amount FROM history WHERE fromid = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"fromid", "toid", "time", "amount"}).
			AddRow(walletHistory[0].From, walletHistory[0].To, walletHistory[0].Time, walletHistory[0].Amount).
			AddRow(walletHistory[1].From, walletHistory[1].To, walletHistory[1].Time, walletHistory[1].Amount))

	r := chi.NewRouter()
	r.Get("/check/v1/wallet/{id}/history", HistoryWalletHandler(storage))

	req, err := http.NewRequest("GET", "/check/v1/wallet/1/history", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}

	expectedResponse, err := json.Marshal(walletHistory)
	if err != nil {
		t.Fatalf("could not marshal expected response: %v", err)
	}

	if rr.Body.String() != string(expectedResponse) {
		t.Errorf("expected response body %s, got %s", expectedResponse, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInfoWalletHandler(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := postgres.NewStorage(db)

	walletID := 1
	expectedWallet := postgres.Wallet{
		ID:      walletID,
		Balance: 100.0,
	}

	mock.ExpectQuery(`select \* from wallet where id = \$1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).
			AddRow(expectedWallet.ID, expectedWallet.Balance))

	r := chi.NewRouter()
	r.Get("/check/v1/wallet/{id}", InfoWalletHandler(storage))

	req, err := http.NewRequest("GET", "/check/v1/wallet/"+strconv.Itoa(walletID), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse, err := json.Marshal(expectedWallet)
	if err != nil {
		t.Fatalf("could not marshal expected response: %v", err)
	}

	assert.JSONEq(t, string(expectedResponse), rr.Body.String())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
