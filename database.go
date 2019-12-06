package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
  host     = "localhost"
  port     = 5432
	user     = "kyle"
	password = "test"
  dbname   = "standard"
)


func init_db() *sql.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func saveUnique(item uniqueItem, db *sql.DB) {
	q := `
INSERT INTO uniques
(id, name, corrupted, original_price, original_price_currency)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE
SET name=$2, corrupted=$3, original_price=$4, original_price_currency=$5, created_at=now();`
	_, err := db.Exec(q, item.Id, item.Name, item.Corrupted, item.OriginalPrice.Price, item.OriginalPrice.Currency)
	if err != nil {
		log.Fatal(err)
	}
}

func saveCurrency(item currencyItem, db *sql.DB) {
	q := `
INSERT INTO currency
(id, type, original_price, original_price_currency, original_quantity)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE
SET type=$2, original_price=$3, original_price_currency=$4, original_quantity=$5, created_at=now();`
	_, err := db.Exec(q, item.Id, item.Type, item.OriginalPrice.Price, item.OriginalPrice.Currency, item.OriginalQuantity)
	if err != nil {
		log.Fatal(err)
	}
}

func saveDivination(item divinationCardItem, db *sql.DB) {
	q := `
INSERT INTO divination_cards
(id, name, mods, max_stack_size, original_price, original_quantity, original_price_currency)
VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE
SET name=$2, mods=$3, max_stack_size=$4, original_price=$5, original_quantity=$6, original_price_currency=$7, created_at=now();`
	_, err := db.Exec(q, item.Id, item.Name, item.Mods, item.MaxStackSize, item.OriginalPrice.Price, item.OriginalQuantity, item.OriginalPrice.Currency)
	if err != nil {
		log.Fatal(err)
	}
}

func saveApiPage(id string, db *sql.DB) {
	q := `
INSERT INTO api_pages
(id) VALUES ($1);`
	_, err := db.Exec(q, id)
	if err != nil {
		log.Fatal(err)
	}
}

func getLatestApiPage(db *sql.DB) string {
	q := `
SELECT id FROM api_pages ORDER BY created_at DESC LIMIT 1;`
	row := db.QueryRow(q)
	var id string
	if err := row.Scan(&id); err != nil {
		log.Fatal(err)
	}

	return id
}
