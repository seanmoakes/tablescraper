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

type Table struct {
	Rows [][]string
}

func NewTable() *Table {
	return &Table{
		Rows: make([][]string, 0),
	}
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) CalculateColumnWidths() []int {
	if len(t.Rows) == 0 {
		return nil
	}

	colWidths := make([]int, len(t.Rows[0]))
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	return colWidths
}

func (t *Table) Print() {
	if len(t.Rows) == 0 {
		fmt.Println("(empty table)")
		return
	}

	colWidths := t.CalculateColumnWidths()
	for _, row := range t.Rows {
		for i, cell := range row {
			// Print cell content with padding
			fmt.Printf(" %-*s |", colWidths[i], cell)
		}
		fmt.Println() // Move to the next line after printing a row
	}
}

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

func run() int {
	flag.Parse()
	args := flag.Args()
	url := args[0]

	tables, err := ScrapeTables(url)
	if err != nil {
		log.Fatalf("Failed to visit URL: %v", err)
	}

	// Print all tables
	for i, table := range tables {
		fmt.Printf("Table %d:\n", i+1)
		table.Print()
		fmt.Println("---") // Separator between tables
	}

	return 0
}

func main() {
	output := run()
	os.Exit(output)
}
