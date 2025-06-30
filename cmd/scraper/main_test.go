package scraper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable()
	if table == nil {
		t.Error("Expected a new Table instance, got nil")
	}
	if len(table.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(table.Rows))
	}
}

func TestAddRow(t *testing.T) {
	table := NewTable()
	row := []string{"Name", "Age", "Occupation"}
	table.AddRow(row)

	if len(table.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.Rows))
	}
	if table.Rows[0][0] != "Name" || table.Rows[0][1] != "Age" || table.Rows[0][2] != "Occupation" {
		t.Error("Row data does not match expected values")
	}
}

func TestCalculateColumnWidths(t *testing.T) {
	table := NewTable()
	table.AddRow([]string{"Name", "Age", "Occupation"})
	table.AddRow([]string{"John Doe", "28", "Software Engineer"})
	table.AddRow([]string{"Jane Smith", "34", "Data Scientist"})

	widths := table.CalculateColumnWidths()
	expectedWidths := []int{10, 3, 17} // Based on the data above

	if len(widths) != len(expectedWidths) {
		t.Errorf("Expected %d columns, got %d", len(expectedWidths), len(widths))
	}
	for i, width := range widths {
		if width != expectedWidths[i] {
			t.Errorf("Expected width %d for column %d, got %d", expectedWidths[i], i, width)
		}
	}
}

func TestPrint(t *testing.T) {
	// This test is more about ensuring the function runs without errors.
	// Capturing stdout for exact output comparison is more complex and not always necessary.
	table := NewTable()
	table.AddRow([]string{"Name", "Age", "Occupation"})
	table.AddRow([]string{"John Doe", "28", "Software Engineer"})
	var writer io.Writer = os.Stdout // Default to stdout

	// Simply call Print and ensure it doesn't panic
	table.Print(writer)
}

func TestScrapeTables(t *testing.T) {
	// Create a test server with a mock HTML response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
			<body>
				<table>
					<tr><th>Name</th><th>Age</th></tr>
					<tr><td>John Doe</td><td>28</td></tr>
					<tr><td>Jane Smith</td><td>34</td></tr>
				</table>
			</body>
		</html>
		`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	// Scrape tables from the test server URL
	tables, err := ScrapeTables(ts.URL)
	if err != nil {
		t.Fatalf("ScrapeTables failed: %v", err)
	}

	// Verify the number of tables
	if len(tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(tables))
	}

	// Verify the table data
	table := tables[0]
	if len(table.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(table.Rows))
	}

	expectedRows := [][]string{
		{"Name", "Age"},
		{"John Doe", "28"},
		{"Jane Smith", "34"},
	}

	for i, row := range table.Rows {
		for j, cell := range row {
			if cell != expectedRows[i][j] {
				t.Errorf("Expected cell %s at row %d, column %d, got %s", expectedRows[i][j], i, j, cell)
			}
		}
	}
}

func TestScrapeTables_NoTables(t *testing.T) {
	// Create a test server with no tables in the response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
			<body>
				<p>No tables here</p>
			</body>
		</html>
		`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	// Scrape tables from the test server URL
	tables, err := ScrapeTables(ts.URL)
	if err != nil {
		t.Fatalf("ScrapeTables failed: %v", err)
	}

	// Verify no tables were found
	if len(tables) != 0 {
		t.Errorf("Expected 0 tables, got %d", len(tables))
	}
}
