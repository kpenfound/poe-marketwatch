package main

import (
	"fmt"
)

// TODO:
// - Store items
// - Calculate currency exchange

func main() {
	fmt.Println("Getting market data...")

	apiRes := apiGet("525912291-543063257-513377713-586182961-557264457")

	var items []stashItem
	categories := []int{3, 5, 6}
	league := "Standard"

	for _, s := range apiRes.Stashes {
		newItems := findItemsByCategoryAndLeague(s, categories, league)
		items = append(items, newItems...)
	}

	fmt.Printf("%+v\n", items)
	fmt.Println(apiRes.NextChangeId)

	processItems(items)
}

func processItems(items []stashItem) {
	for _, i := range items {
		switch i.FrameType {
		case 3:
			u := getUnique(i)
			saveUnique(u)
		case 5:
			c := getCurrency(i)
			saveCurrency(c)
		case 6:
			d := getDivination(i)
			saveDivination(d)
		}
	}
}
