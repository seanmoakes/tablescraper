package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func run() int {
	flag.Parse()
	args := flag.Args()
	url := args[0]

	// Initialize Colly collector
	c := colly.NewCollector()

	// Set up callbacks to handle scraping
	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Use goquery to parse the table
		table := e.DOM
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			var rowData []string
			row.Find("th, td").Each(func(j int, cell *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cell.Text()))
			})
			fmt.Println(strings.Join(rowData, " | "))
		})
		fmt.Println("---") // Separator between tables
	})

	// Start scraping
	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Failed to visit URL: %v", err)
	}

	return 0
}

func main() {
	output := run()

	os.Exit(output)
}
