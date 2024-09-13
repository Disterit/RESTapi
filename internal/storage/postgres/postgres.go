package postgres

import (
	"RESTapi/internal/storage"
	"fmt"

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
		fmt.Errorf("%s %w", op, err)
		return Wallet{}, err
	}
	return wallet, nil
}

// id balance
// id time fromID toID amount
