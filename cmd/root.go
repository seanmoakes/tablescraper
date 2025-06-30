/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/seanmoakes/tablescraper/cmd/scraper"
	"github.com/spf13/cobra"
)

var (
	Verbose        bool
	OutputFilePath string
)

var rootCmd = &cobra.Command{
	Use:   "tablescraper",
	Short: "An app to scrape all existing tables from a webpage",
	Long: `Tablescraper is a CLI application to scrape table data from a website.
This app allows you to view the scraped data in the terminal, save scraped data to a file,
and gives you the option to parse the output into several formats such as markdown.`,

	// Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	return scraper.Scrape(args[0])
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scraper.ScrapeCmd)
}
