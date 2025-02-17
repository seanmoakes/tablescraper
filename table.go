package main

import (
	"fmt"
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

func printRow(row []string, colWidths []int) {
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
				fmt.Printf("| %*s ", colWidths[j], cell)
			} else {
				fmt.Printf("| %-*s ", colWidths[j], cell)
			}
		}
		// Right edge border
		fmt.Println("|")
	} else {
		printMultiLineRow(row, colWidths)
	}
}

func printMultiLineRow(row []string, colWidths []int) {
	var secondRow []string
	for j, cell := range row {
		before, after, _ := strings.Cut(cell, "\n")
		secondRow = append(secondRow, strings.TrimSpace(after))

		if isNumeric(before) {
			fmt.Printf("| %*s ", colWidths[j], strings.TrimSpace(before))
		} else {
			fmt.Printf("| %-*s ", colWidths[j], strings.TrimSpace(before))
		}
	}
	// Right edge border
	fmt.Println("|")

	// print the additionnal row
	printRow(secondRow, colWidths)
}

func (t *Table) Print() {
	if len(t.Rows) == 0 {
		fmt.Println("(empty table)")
		return
	}

	colWidths := t.CalculateColumnWidths()
	printBorder := func() {
		for _, width := range colWidths {
			fmt.Print("+" + strings.Repeat("-", width+2))
		}
		fmt.Println("+")
	}

	// Top border
	printBorder()
	for i, row := range t.Rows {
		printRow(row, colWidths)

		// Header Separator
		if i == 0 {
			printBorder()
		}
	}
	// Bottom border
	printBorder()
}
