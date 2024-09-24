package main

import (
	"fmt"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"time"
)

func main() {
	a := time.Now()

	fmt.Println(a)

}

func printa(time time.Time) {
	fmt.Println(time)
}
