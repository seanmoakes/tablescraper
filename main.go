package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

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
