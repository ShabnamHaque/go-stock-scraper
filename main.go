package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {
	tickers := []string{
		"MSFT", "IBM", "AAPL", "MMM", "JPM", "GE",
	}
	stocks := []Stock{}
	//slice of (collection of) Stock called stocks made.

	c := colly.NewCollector()
	//initialises a new instance of colly. call any fucntion using c.

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting stock : ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}
		stock.company = e.ChildText("h1")
		fmt.Println("Company: ", stock.company)
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("Company: ", stock.price)
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		fmt.Println("Company: ", stock.change)

		stocks = append(stocks, stock)
	})

	c.Wait() //wait before u start writing onto csv file.

	for _, t := range tickers {

		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}
	fmt.Println(stocks)
	//print onto terminal.

	file, err := os.Create("stocks.csv")

	if err != nil {
		log.Fatal("Failed to create CSV file!", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	headers := []string{
		"company", "price", "change",
	}
	writer.Write(headers) 
	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}
	defer writer.Flush()
}

//base url  "https://finance.yahoo.com/quote/"+ticker+"/"
