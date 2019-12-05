package main

import (
	"fmt"
	"database/sql"
)

// TODO:
// - Store items
// - Calculate currency exchange

func main() {
	fmt.Println("Getting market data...")

	apiRes := apiGet("525912291-543063257-513377713-586182961-557264457")
	db := init_db()

	var items []stashItem
	categories := []int{3, 5, 6}
	league := "Standard"

	for _, s := range apiRes.Stashes {
		newItems := findItemsByCategoryAndLeague(s, categories, league)
		items = append(items, newItems...)
	}

	fmt.Printf("%+v\n", items)

	processItems(items, db)

	db.Close()

	fmt.Println(apiRes.NextChangeId)
}

func processItems(items []stashItem, db *sql.DB) {
	for _, i := range items {
		switch i.FrameType {
		case 3:
			u := getUniqueFromStashItem(i)
			saveUnique(u, db)
		case 5:
			c := getCurrencyFromStashItem(i)
			saveCurrency(c, db)
		case 6:
			d := getDivinationFromStashItem(i)
			saveDivination(d, db)
		}
	}
}
