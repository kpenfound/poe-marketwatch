package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

// TODO:
// - Save out raw data to backup files
// - Track time series price data

func main() {
	sigs := make(chan os.Signal, 1)
	done := false
	var wg sync.WaitGroup

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println(sig)
		done = true
	}()

	log.Println("Getting market data...")
	categories := []int{3, 5, 6}
	league := os.Getenv("LEAGUE")

	db := init_db()
	re := init_currency()
	i := 1
	pageId := getLatestApiPage(db)

	for !done {
		apiRes := apiGet(pageId)
		log.Printf("Retrieved API page %v\n", i)

		// Async process all the data so we can keep requesting
		wg.Add(1)
		go handleResponseAsync(league, categories, apiRes, db, re, &wg)

		log.Printf("Last Change: %v\n", pageId)
		log.Printf("Next Change: %v\n", apiRes.NextChangeId)
		if pageId == apiRes.NextChangeId || apiRes.NextChangeId == "" {
			// Wait a minute if we're caught up
			log.Println("All caught up.... waiting")
			time.Sleep(time.Second * 60)
		} else {
			pageId = apiRes.NextChangeId
			saveApiPage(apiRes.NextChangeId, db)
		}

		if i%100 == 0 { // Wait a minute every 100 iterations
			time.Sleep(time.Second * 60)
		}
		i++
	}

	log.Println("Waiting for workers to finish")
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
	log.Printf("Processing %d items...\n", len(items))
	for _, i := range items {
		if isPrice(i.Note, re) {
			switch i.FrameType { // Easiest way to tell what kind of item it is
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
