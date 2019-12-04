package main

import (
	"fmt"
)

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
}
