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
(id, name, corrupted, original_price, calculated_price)
VALUES ($1, $2, $3, $4, $5);`
	_, err := db.Exec(q, item.Id, item.Name, item.Corrupted, item.OriginalPrice, item.CalculatedPrice)
	if err != nil {
		log.Fatal(err)
	}
}

func saveCurrency(item currencyItem, db *sql.DB) {
	q := `
INSERT INTO currency
(id, type, original_price, calculated_price)
VALUES ($1, $2, $3, $4);`
	_, err := db.Exec(q, item.Id, item.Type, item.OriginalPrice, item.CalculatedPrice)
	if err != nil {
		log.Fatal(err)
	}
}

func saveDivination(item divinationCardItem, db *sql.DB) {
	q := `
INSERT INTO divination_cards
(id, name, mods, max_stack_size, original_price, original_quantity, calculated_price)
VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := db.Exec(q, item.Id, item.Name, item.Mods, item.MaxStackSize, item.OriginalPrice, item.OriginalQuantity, item.CalculatedPrice)
	if err != nil {
		log.Fatal(err)
	}
}
