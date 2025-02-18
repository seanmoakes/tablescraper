package scraper

import (
	"fmt"
)

func Scrape(url string) error {
	tables, err := ScrapeTables(url)
	if err != nil {
		return fmt.Errorf("failed to visit URL: %v", err)
	}

	// Print all tables
	for i, table := range tables {
		fmt.Printf("Table %d:\n", i+1)
		table.Print()
		fmt.Println("---") // Separator between tables
	}

	return nil
}
