/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package scraper

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
)

var (
	Verbose    bool
	outputFile string
	format     string
)

// scrapeCmd represents the scrape command
var ScrapeCmd = &cobra.Command{
	Use:   "scrape [URL]",
	Short: "scrape the given URL for tables",
	Long: `Scrape the given URL for tables and output them. Default output 
  is to the terminal and default format is as text tables. Both output and
  format can be defined through the use of flags.`,

	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Scrape(args[0])
	},
}

func init() {
	ScrapeCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Explain what is happenning at every step of the scraping process.")
	ScrapeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write to the given filepath")
	ScrapeCmd.Flags().StringVarP(&format, "format", "f", "table", "Output tables in the given format, valid options are table, csv, json, markdown and html")
}

func Scrape(url string) error {
	tables, err := ScrapeTables(url)
	if err != nil {
		return fmt.Errorf("failed to visit URL: %v", err)
	}

	var writer io.Writer = os.Stdout // Default to stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer file.Close()
		writer = file
	}

	if format == "json" {
		PrintAllTablesJSON(writer, tables, url)
	} else {
		// Print all tables
		for i, table := range tables {
			switch format {
			case "table":
				fmt.Fprintf(writer, "Table %d:\n", i+1)
				table.Print(writer)
			case "markdown":
				fmt.Fprintf(writer, "## Table %d\n", i+1)
				table.PrintMarkdown(writer)
			case "csv":
				table.PrintCSV(writer)
			case "json":
				table.PrintJSON(writer)
			case "html":
				table.PrintHTML(writer)
			default:
				log.Fatalf("Unsupported format: %s", format)
			}
			fmt.Fprintln(writer) // Add a blank line between tables
		}
	}

	if outputFile != "" {
		fmt.Printf("Output written to: %s\n", outputFile)
	}

	return nil
}

func ScrapeTables(url string) ([]*Table, error) {
	// Enable Chromedp logging and set a custom User-Agent
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Run in non-headless mode for debugging
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-logging", true), // Enable logging
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for the operation
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Run Chromedp to scrape the page
	var tableNodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery), // Wait for the body to be visible
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Page loaded, looking for tables...")
			return nil
		}),
		chromedp.Nodes("table", &tableNodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("Found %d tables", len(tableNodes))
			return nil
		}),
	)
	if err != nil {
		log.Fatalf("Failed to scrape page: %v", err)
	}

	// Slice to store all tables
	var tables []*Table

	// Process each table
	for i, tableNode := range tableNodes {

		// Get the outer HTML of the table
		var tableHTML string
		err := chromedp.Run(ctx,
			chromedp.OuterHTML(tableNode.FullXPath(), &tableHTML),
		)
		if err != nil {
			log.Printf("Failed to get HTML for table %d: %v", i+1, err)
			continue
		}

		// Parse the table HTML using goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(tableHTML))
		if err != nil {
			log.Printf("Failed to parse table %d: %v", i+1, err)
			continue
		}

		table := NewTable()
		doc.Find("tr").Each(func(j int, row *goquery.Selection) {
			var rowData []string
			row.Find("th, td").Each(func(j int, cell *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cell.Text()))
			})
			table.AddRow(rowData)
		})
		tables = append(tables, table)
	}

	return tables, nil
}
