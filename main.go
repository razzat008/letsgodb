package main

import (
	"fmt"
	"os"

	par "github.com/razzat008/letsgodb/internal/Parser"
	repl "github.com/razzat008/letsgodb/internal/REPl"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
	catalog "github.com/razzat008/letsgodb/internal/catalog"
	"github.com/razzat008/letsgodb/internal/db"
	"github.com/razzat008/letsgodb/internal/storage"
)

// to print help message
func printHelp() {
	println("letsgodb Help:")
	println("Every command MUST end with a semicolon(;).")
	println("  Type SQL commands to interact with the database.")
	println("  Type 'help;' to see this message.")
	println("  Type '\\e;' to exit.")
	println("  To learn more about letsgodb, visit https://github.com/razzat008/letsgodb")
}

// ExecuteStatement handles parsed statements and interacts with the catalog and row storage.
func ExecuteStatement(stmt par.Statement, cat *catalog.Catalog) error {
	switch s := stmt.(type) {
	case *par.CreateTableStatement:
		// CREATE TABLE
		err := cat.AddTable(s.TableName, s.Columns)
		if err != nil {
			return fmt.Errorf("CREATE TABLE failed: %w", err)
		}
		fmt.Println("Table created:", s.TableName)
	case *par.SelectStatement:
		// SELECT
		schema := cat.GetTable(s.Table)
		if schema == nil {
			return fmt.Errorf("table %q does not exist", s.Table)
		}
		pager := storage.NewPager(s.Table + ".db")
		rows := db.ReadAllRows(pager)
		// Print header
		fmt.Println(schema.Columns)
		for _, row := range rows {
			// SELECT *: print all columns
			if len(s.Columns) == 1 && s.Columns[0] == "*" {
				fmt.Println(row)
			} else {
				// Print only requested columns
				var selected []string
				for _, col := range s.Columns {
					for i, schemaCol := range schema.Columns {
						if col == schemaCol && i < len(row) {
							selected = append(selected, row[i])
						}
					}
				}
				fmt.Println(selected)
			}
		}
	case *par.InsertStatement:
		// INSERT
		schema := cat.GetTable(s.Table)
		if schema == nil {
			return fmt.Errorf("table %q does not exist", s.Table)
		}
		// Flatten [][]string to []string for storage
		var flatValues []string
		for _, v := range s.Values {
			flatValues = append(flatValues, v...)
		}
		pager := storage.NewPager(s.Table + ".db")
		_, err := db.InsertRow(pager, flatValues)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
		fmt.Println("Row inserted!")
	default:
		return fmt.Errorf("unsupported statement type")
	}
	return nil
}

// main entry point of the program
func main() {
	// Initialize the catalog (persisted in a file, e.g., "catalog.db")
	cat, err := catalog.NewCatalog("catalog.db")
	if err != nil {
		println("Failed to initialize catalog:", err.Error())
		os.Exit(1)
	}

	lineBuffer := repl.InitLineBuffer()
	printHelp()
	for {
		repl.PrintDB()
		lineBuffer.UserInput()
		input := string(lineBuffer.Buffer)
		if input == "help" {
			printHelp()
			lineBuffer.Reset()
			continue
		} else if input == "\\e" {
			println("Exiting letsgodb...")
			println("Bye!!")
			os.Exit(0)
			break
		}
		tokens := tok.Tokenizer(lineBuffer)
		fmt.Println(tokens)
		stmt := par.ParseProgram(tokens)
		if stmt == nil {
			// Print a friendly message for parse errors
			fmt.Println("Error: Failed to parse statement. Please check your SQL syntax.")
			lineBuffer.Reset()
			continue
		}
		err = ExecuteStatement(stmt, cat)
		if err != nil {
			fmt.Println("Error:", err)
		}
		lineBuffer.Reset()
	}
}
