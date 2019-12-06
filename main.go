package main

import (
	"fmt"
	"database/sql"
	"regexp"
	"sync"
	"os"
	"os/signal"
	"syscall"
)

// TODO:
// - Calculate currency
// - Analyze prices

func main() {
	sigs := make(chan os.Signal, 1)
	done := false
	var wg sync.WaitGroup

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done =true
	}()
	
	fmt.Println("Getting market data...")
	categories := []int{3, 5, 6}
	league := "Standard"
	//pageId := "525912708-543063663-513378145-586183293-557264789"
	db := init_db()
	re := init_currency()
	i := 1

	for !done {
		pageId := getLatestApiPage(db)
		apiRes := apiGet(pageId)
		fmt.Printf("Retrieved API page %v\n", i)
	
		wg.Add(1)
		go handleResponseAsync(league, categories, apiRes, db, re, &wg)
	
		saveApiPage(apiRes.NextChangeId, db)
		i++
	}

	fmt.Println("Waiting for workers to finish")
	wg.Wait()

	db.Close()
}

func handleResponseAsync(league string, categories []int, res *apiResponse, db *sql.DB, re *regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	var items []stashItem
	
	for _, s := range res.Stashes {
		newItems := findItemsByCategoryAndLeague(s, categories, league)
		items = append(items, newItems...)
	}

	processItems(items, db, re)
}

func processItems(items []stashItem, db *sql.DB, re *regexp.Regexp) {
	fmt.Printf("Processing %d items...\n", len(items))
	for _, i := range items {
		if isPrice(i.Note, re) {
		//	fmt.Printf("Saving %+v\n", i)
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
}
