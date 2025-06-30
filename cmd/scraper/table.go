package scraper

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"strings"
	"unicode"
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

func isNumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}

func (t *Table) CalculateColumnWidths() []int {
	if len(t.Rows) == 0 {
		return nil
	}

	colWidths := make([]int, len(t.Rows[0]))
	for _, row := range t.Rows {
		for i, cell := range row {
			for _, line := range strings.Split(cell, "\n") {
				if len(line) > colWidths[i] {
					colWidths[i] = len(line)
				}
			}
		}
	}
	return colWidths
}

func printRow(row []string, colWidths []int, w io.Writer) {
	// Check for more than one line in the row
	isMultiLineRow := false
	for _, c := range row {
		if strings.Contains(c, "\n") {
			isMultiLineRow = true
		}
	}

	if !isMultiLineRow {
		for j, cell := range row {
			if isNumeric(cell) {
				fmt.Fprintf(w, "| %*s ", colWidths[j], cell)
			} else {
				fmt.Fprintf(w, "| %-*s ", colWidths[j], cell)
			}
		}
		// Right edge border
		fmt.Fprintln(w, "|")
	} else {
		printMultiLineRow(row, colWidths, w)
	}
}

func printMultiLineRow(row []string, colWidths []int, w io.Writer) {
	var secondRow []string
	for j, cell := range row {
		before, after, _ := strings.Cut(cell, "\n")
		secondRow = append(secondRow, strings.TrimSpace(after))

		if isNumeric(before) {
			fmt.Fprintf(w, "| %*s ", colWidths[j], strings.TrimSpace(before))
		} else {
			fmt.Fprintf(w, "| %-*s ", colWidths[j], strings.TrimSpace(before))
		}
	}
	// Right edge border
	fmt.Fprintln(w, "|")

	// print the additionnal row
	printRow(secondRow, colWidths, w)
}

func (t *Table) Print(w io.Writer) {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, "(empty table)")
		return
	}

	colWidths := t.CalculateColumnWidths()
	printBorder := func() {
		for _, width := range colWidths {
			fmt.Fprint(w, "+"+strings.Repeat("-", width+2))
		}
		fmt.Fprintln(w, "+")
	}

	// Top border
	printBorder()
	for i, row := range t.Rows {
		printRow(row, colWidths, w)

		// Header Separator
		if i == 0 {
			printBorder()
		}
	}
	// Bottom border
	printBorder()
}

func (t *Table) PrintMarkdown(w io.Writer) {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, "(empty table)")
		return
	}

	// Print header row
	for _, cell := range t.Rows[0] {
		fmt.Fprintf(w, "| %s ", cell)
	}
	fmt.Fprintln(w, "|")

	// Print separator row
	for range t.Rows[0] {
		fmt.Fprint(w, "| --- ")
	}
	fmt.Fprintln(w, "|")

	// Print data rows
	for _, row := range t.Rows[1:] {
		for _, cell := range row {
			fmt.Fprintf(w, "| %s ", cell)
		}
		fmt.Fprintln(w, "|")
	}
}

// PrintCSV prints the table in CSV format
func (t *Table) PrintCSV(w io.Writer) {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, "(empty table)")
		return
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	for _, row := range t.Rows {
		if err := csvWriter.Write(row); err != nil {
			fmt.Fprintf(w, "Error writing CSV: %v\n", err)
			return
		}
	}
}

func (t *Table) PrintJSON(w io.Writer) {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, "[]")
		return
	}

	// Convert the table to a slice of maps for JSON encoding
	var jsonData []map[string]string
	headers := t.Rows[0] // First row is the header
	for _, row := range t.Rows[1:] {
		rowMap := make(map[string]string)
		for i, cell := range row {
			rowMap[headers[i]] = cell
		}
		jsonData = append(jsonData, rowMap)
	}

	// Encode the data as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty-print with indentation
	if err := encoder.Encode(jsonData); err != nil {
		fmt.Fprintf(w, "Error encoding JSON: %v\n", err)
	}
}

func PrintAllTablesJSON(w io.Writer, tables []*Table, url string) {
	// Create a slice to hold all tables' JSON data
	var jsonData []map[string]interface{}

	// Convert each table to JSON format
	for i, table := range tables {
		tableData := map[string]interface{}{
			"name": fmt.Sprintf("Table %d", i+1),
			"url":  url,
			"rows": table.ToJSON(),
		}
		jsonData = append(jsonData, tableData)
	}

	// Encode the data as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty-print with indentation
	if err := encoder.Encode(jsonData); err != nil {
		fmt.Fprintf(w, "Error encoding JSON: %v\n", err)
	}
}

func (t *Table) ToJSON() []map[string]string {
	if len(t.Rows) == 0 {
		return nil
	}

	// Convert the table to a slice of maps for JSON encoding
	var jsonData []map[string]string
	headers := t.Rows[0] // First row is the header
	for _, row := range t.Rows[1:] {
		rowMap := make(map[string]string)
		for i, cell := range row {
			rowMap[headers[i]] = cell
		}
		jsonData = append(jsonData, rowMap)
	}

	return jsonData
}

func (t *Table) PrintHTML(w io.Writer) error {
	if len(t.Rows) == 0 {
		_, err := fmt.Fprintln(w, "<!-- Empty table -->")
		return err
	}

	// Write HTML header with basic CSS
	_, err := fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <style>
        table {
            border-collapse: collapse;
            width: 100%%;
            margin: 1em 0;
            font-family: sans-serif;
        }
        th, td {
            padding: 0.5em;
            border: 1px solid #ddd;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
    </style>
</head>
<body>
`)
	if err != nil {
		return err
	}

	// Start table with semantic markup
	if _, err := fmt.Fprintln(w, "<table>"); err != nil {
		return err
	}

	// Header row in <thead>
	if _, err := fmt.Fprintln(w, "    <thead>"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, "        <tr>"); err != nil {
		return err
	}
	for _, cell := range t.Rows[0] {
		if _, err := fmt.Fprintf(w, "\n            <th>%s</th>", html.EscapeString(cell)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "\n        </tr>"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "    </thead>"); err != nil {
		return err
	}

	// Data rows in <tbody>
	if _, err := fmt.Fprintln(w, "    <tbody>"); err != nil {
		return err
	}
	for _, row := range t.Rows[1:] {
		if _, err := fmt.Fprint(w, "        <tr>"); err != nil {
			return err
		}
		for _, cell := range row {
			if _, err := fmt.Fprintf(w, "\n            <td>%s</td>", html.EscapeString(cell)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "\n        </tr>"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "    </tbody>"); err != nil {
		return err
	}

	// Close table and document
	if _, err := fmt.Fprintln(w, "</table>"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "</body>\n</html>"); err != nil {
		return err
	}

	return nil
}
