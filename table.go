package main

import "fmt"

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
