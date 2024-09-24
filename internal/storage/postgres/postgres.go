package postgres

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type WalletHistory struct {
	Time   time.Time `json:"time"`
	From   int       `json:"from"`
	To     int       `json:"to"`
	Amount float32   `json:"amount"`
}

type Wallet struct {
	ID      int     `json:"id"`
	Balance float32 `json:"balance"`
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateWallet() (Wallet, error) {

	const op = "storage.postgres.CreateWallet()"

	var wallet Wallet

	query := `insert into wallet (balance) values (100) returning id, balance`
	err := s.db.QueryRow(query).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		log.Printf("%s %v", op, err)
		return Wallet{}, err
	}
	return wallet, nil
}

func (s *Storage) TransferMoney(amount float32, walletID, walletToID int) (int, error) {

	const op = "storage.postgres.TransferMoney()"

	tx, err := s.db.Begin()
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

func (s *Storage) GetWalletInfo(walletID int) (Wallet, error) {
	const op = "storage.postgres.GetWalletInfo()"

	query := `select * from wallet where id = $1`
	row := s.db.QueryRow(query, walletID)

	var wallet Wallet

	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		log.Printf("%s %v", op, err)
		return wallet, err
	}

	return wallet, nil
}

func (s *Storage) HistoryWallet(id int) ([]WalletHistory, int) {
	const op = "storage.postgres.HistoryWallet()"

	query := `SELECT fromid, toid, time, amount FROM history WHERE fromid = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Printf("%s %v", op, err)
		return nil, http.StatusNotFound
	}

	var walletHistory []WalletHistory

	for rows.Next() {
		var wallet WalletHistory
		err = rows.Scan(&wallet.From,
			&wallet.To,
			&wallet.Time,
			&wallet.Amount)
		if err != nil {
			log.Printf("%s %v", op, err)
			return nil, http.StatusNotFound
		}

		walletHistory = append(walletHistory, wallet)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("%s %v", op, err)
		return nil, http.StatusNotFound
	}

	return walletHistory, http.StatusOK
}

func (s *Storage) AddHistoryWallet(date time.Time, fromID, toID int, amount float32) error {
	const op = "storage.postgres.AddHistoryWallet()"

	query := `insert into history (fromID, toID, time, amount) values ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, fromID, toID, date, amount)
	if err != nil {
		log.Printf("%s %v", op, err)
		return err
	}

	return nil
}

func (s *Storage) CreateTable() {
	const op = "storage.postgres.CreateTable()"

	createWalletTable := `
    CREATE TABLE IF NOT EXISTS wallet (
        id SERIAL PRIMARY KEY,
        balance NUMERIC(10, 2) DEFAULT 100
    );`

	createHistoryTable := `
    CREATE TABLE IF NOT EXISTS history (
        id SERIAL PRIMARY KEY,
        fromid INTEGER NOT NULL,
        toid INTEGER NOT NULL,
        time TIMESTAMPTZ NOT NULL,
        amount INTEGER NOT NULL,
        foreign key (fromID) references wallet (id), 
		foreign key (toID) references wallet (id)
    );`

	_, err := s.db.Exec(createWalletTable)
	if err != nil {
		log.Printf("%s %v", op, err)
	}

	_, err = s.db.Exec(createHistoryTable)
	if err != nil {
		log.Printf("%s %v", op, err)
	}

	return
}
