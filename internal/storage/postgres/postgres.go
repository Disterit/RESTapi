package postgres

import (
	"RESTapi/internal/storage"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Wallet struct {
	ID      int     `json:"id"`
	Balance float32 `json:"balance"`
}

func CreateWallet() (Wallet, error) {

	const op = "storage.postgres.CreateWallet()"

	db := storage.Connection()
	defer db.Close()

	var wallet Wallet

	query := `insert into wallet (balance) values (100) returning id, balance`
	err := db.QueryRow(query).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		log.Printf("%s %v", op, err)
		return Wallet{}, err
	}
	return wallet, nil
}

func TransferMoney(amount float32, walletID, walletToID int) (int, error) {

	const op = "storage.postgres.TransferMoney()"

	db := storage.Connection()
	defer db.Close()

	tx, err := db.Begin() // Начало транзакции
	if err != nil {
		log.Printf("%s %v", op, err)
		return http.StatusBadRequest, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result, err := tx.Exec(`UPDATE wallet SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, walletID)
	if err != nil {
		log.Printf("%s %v", op, err)
		return http.StatusBadRequest, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Printf("%s not enough balance or no rows affected", op)
		return http.StatusBadRequest, err
	}

	_, err = tx.Exec(`UPDATE wallet SET balance = balance + $1 WHERE id = $2`, amount, walletToID)
	if err != nil {
		log.Printf("%s %v", op, err)
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
