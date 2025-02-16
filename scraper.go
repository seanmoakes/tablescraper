package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeTables(url string) ([]*Table, error) {
	// Initialize Colly collector
	c := colly.NewCollector()

	// Slice to store all tables
	var tables []*Table

	// Set up callbacks to handle scraping
	c.OnHTML("table", func(e *colly.HTMLElement) {
		table := NewTable()
		// Use goquery to parse the table
		e.DOM.Find("tr").Each(func(i int, row *goquery.Selection) {
			var rowData []string
			row.Find("th, td").Each(func(j int, cell *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cell.Text()))
			})
			table.AddRow(rowData)
		})
		tables = append(tables, table)
	})

	// Start scraping
	err := c.Visit(url)
	if err != nil {
		return []*Table{}, fmt.Errorf("failed to visit URL: %v", err)
	}

	return tables, nil
}
