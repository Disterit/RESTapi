package main

import (
	"fmt"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"strings"
)

func main() {
	const op = "http://localhost:8081/api/v1/wallet/{walletId}/send"

	walletID := strings.Split(op, "/")[6]

	fmt.Println(walletID)
}
