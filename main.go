package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/gocolly/colly/v2"
)

type Stock struct {
	Name              string `json:"Name"`
	Symbol            string `json:"Symbol"`
	MarketCap         string `json:"MarketCap"`
	Price             string `json:"Price"`
	CirculatingSupply string `json:"CirculatingSupply"`
	Volume24h         string `json:"Volume24h"`
	Change1h          string `json:"Change1h"`
	Change24h         string `json:"Change24h"`
	Change7d          string `json:"Change7d"`
}

func main() {
	fName := "cryptocoinmarketcap.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate collector with custom user agents
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36"),
		colly.AllowedDomains("coinmarketcap.com"),
		colly.MaxBodySize(0),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var stocks []Stock

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			stocks = append(stocks, Stock{
				Name:              e.ChildText(".cmc-table__column-name"),
				Symbol:            e.ChildText(".cmc-table__cell--sort-by__symbol"),
				MarketCap:         e.ChildText(".cmc-table__cell--sort-by__market-cap"),
				Price:             e.ChildText(".cmc-table__cell--sort-by__price"),
				CirculatingSupply: e.ChildText(".cmc-table__cell--sort-by__circulating-supply"),
				Volume24h:         e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
				Change1h:          e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
				Change24h:         e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
				Change7d:          e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
			})
		}()
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	wg.Wait()

	// Encode data to JSON and write to file
	jsonData, err := json.MarshalIndent(stocks, "", "    ")
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	if _, err := file.Write(jsonData); err != nil {
		log.Fatalf("Error writing JSON data to file: %v", err)
	}

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
