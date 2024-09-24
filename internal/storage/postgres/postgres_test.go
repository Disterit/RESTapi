package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	// Создаем мок для sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	// Эмулируем возвращаемый результат для запроса
	rows := sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 100)
	mock.ExpectQuery("insert into wallet").WillReturnRows(rows)

	// Создаем storage с моком базы данных
	storage := NewStorage(db)

	wallet, err := storage.CreateWallet()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Проверяем результат
	if wallet.ID != 1 || wallet.Balance != 100 {
		t.Errorf("unexpected wallet data: got ID %d, Balance %d", wallet.ID, wallet.Balance)
	}

	// Убеждаемся, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
