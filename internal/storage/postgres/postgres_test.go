package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"testing"
	"time"
)

func TestCreateWallet(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 100)
	mock.ExpectQuery("insert into wallet").WillReturnRows(rows)

	storage := NewStorage(db)

	wallet, err := storage.CreateWallet()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if wallet.ID != 1 || wallet.Balance != 100 {
		t.Errorf("unexpected wallet data: got ID %d, Balance %f", wallet.ID, wallet.Balance)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransferMoney(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := NewStorage(db)

	amount := float32(50)
	walletID := 1
	walletToID := 2

	mock.ExpectBegin()

	mock.ExpectExec("UPDATE wallet SET balance = balance -").
		WithArgs(amount, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE wallet SET balance = balance +").
		WithArgs(amount, walletToID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	status, err := storage.TransferMoney(amount, walletID, walletToID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if status != http.StatusOK {
		t.Errorf("unexpected status: got %v, want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetWalletInfo(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := NewStorage(db)

	walletID := 1
	expectedWallet := Wallet{
		ID:      walletID,
		Balance: 100.0,
	}

	rows := sqlmock.NewRows([]string{"id", "balance"}).AddRow(expectedWallet.ID, expectedWallet.Balance)
	mock.ExpectQuery("select \\* from wallet where id = \\$1").
		WithArgs(walletID).
		WillReturnRows(rows)

	wallet, err := storage.GetWalletInfo(walletID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if wallet.ID != expectedWallet.ID || wallet.Balance != expectedWallet.Balance {
		t.Errorf("unexpected wallet info: got %+v, want %+v", wallet, expectedWallet)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHistoryWallet(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := NewStorage(db)

	walletID := 1
	expectedHistory := []WalletHistory{
		{
			From:   walletID,
			To:     2,
			Time:   time.Now(),
			Amount: 100.0,
		},
		{
			From:   walletID,
			To:     3,
			Time:   time.Now(),
			Amount: 50.0,
		},
	}

	rows := sqlmock.NewRows([]string{"fromid", "toid", "time", "amount"}).
		AddRow(expectedHistory[0].From, expectedHistory[0].To, expectedHistory[0].Time, expectedHistory[0].Amount).
		AddRow(expectedHistory[1].From, expectedHistory[1].To, expectedHistory[1].Time, expectedHistory[1].Amount)

	mock.ExpectQuery("SELECT fromid, toid, time, amount FROM history WHERE fromid = \\$1").
		WithArgs(walletID).
		WillReturnRows(rows)

	history, status := storage.HistoryWallet(walletID)
	if status != http.StatusOK {
		t.Errorf("unexpected status: got %v, want %v", status, http.StatusOK)
	}

	if len(history) != len(expectedHistory) {
		t.Fatalf("unexpected length of history: got %d, want %d", len(history), len(expectedHistory))
	}
	for i := range history {
		if history[i].From != expectedHistory[i].From ||
			history[i].To != expectedHistory[i].To ||
			history[i].Amount != expectedHistory[i].Amount ||
			!history[i].Time.Equal(expectedHistory[i].Time) {
			t.Errorf("unexpected history entry at index %d: got %+v, want %+v", i, history[i], expectedHistory[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddHistoryWallet(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	storage := NewStorage(db)

	date := time.Now()
	fromID := 1
	toID := 2
	amount := float32(100.0)

	mock.ExpectExec("insert into history \\(fromID, toID, time, amount\\) values \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(fromID, toID, date, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.AddHistoryWallet(date, fromID, toID, amount)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
