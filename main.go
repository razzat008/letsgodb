package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	println("  Type [\\e;] to exit.")
	println("  CREATE DATABASE dbname; to create a new database.")
	println("  USE dbname; to switch to a database.")
	println("  DROP DATABASE dbname; to drop a database.")
}

func printHelpall() {
	println("Every command MUST end with a semicolon(;).")
	println("  Type 'helpall;' to see list of all commands.")
	println("  Type '\\e;' to exit.")
	println("  -> `CREATE DATABASE dbname`")
	println("  -> `USE dbname`")
	println("  -> `DROP DATABASE dbname`")
	println("  -> `CREATE TABLE tablename (column1 datatype, column2 datatype)`")
	println("  -> `DROP TABLE tablename`")
	println("  -> `INSERT INTO tablename (column1, column2) VALUES (value1, value2)`")
}

// ExecuteStatement handles parsed statements and interacts with the catalog and row storage.
func ExecuteStatement(stmt par.Statement, currentDB *string, cat **catalog.Catalog) error {
	switch s := stmt.(type) {
	case *par.CreateDatabaseStatement:
		// CREATE DATABASE dbname;
		dbDir := filepath.Join("data", s.DatabaseName)
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
		// Create an empty catalog.db if not exists
		catalogPath := filepath.Join(dbDir, "catalog.db")
		if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
			f, ferr := os.Create(catalogPath)
			if ferr != nil {
				return fmt.Errorf("failed to create catalog.db: %w", ferr)
			}
			f.Close()
		}
		fmt.Printf("Database '%s' created.\n", s.DatabaseName)
	case *par.UseDatabaseStatement:
		// USE dbname;
		dbDir := filepath.Join("data", s.DatabaseName)
		catalogPath := filepath.Join(dbDir, "catalog.db")
		if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
			return fmt.Errorf("database '%s' does not exist. Use CREATE DATABASE first.", s.DatabaseName)
		}
		newCat, err := catalog.NewCatalog(catalogPath)
		if err != nil {
			return fmt.Errorf("failed to load catalog for database '%s': %w", s.DatabaseName, err)
		}
		*cat = newCat
		*currentDB = s.DatabaseName
		fmt.Printf("Switched to database '%s'.\n", s.DatabaseName)
	case *par.CreateTableStatement:
		// CREATE TABLE
		if *currentDB == "" {
			return fmt.Errorf("no database selected. Use CREATE DATABASE and USE first.")
		}
		err := (*cat).AddTable(s.TableName, s.Columns)
		if err != nil {
			return fmt.Errorf("CREATE TABLE failed: %w", err)
		}
		fmt.Println("Table created:", s.TableName)
	case *par.ListTablesStatement:
		if *currentDB == "" {
			return fmt.Errorf("no database selected. Use CREATE DATABASE and USE first.")
		}
		tables := (*cat).ListTables()
		if tables == nil {
			fmt.Println("Tables:Empty Database")
		} else {
			fmt.Println("Tables:")
			for _, t := range tables {
				fmt.Println(" -", t.Name, " : ", t.Columns)
			}
		}
	case *par.DropStatement:
		// DROP TABLE
		if s.Table != "" {
			if *currentDB == "" {
				return fmt.Errorf("no database selected. Use CREATE DATABASE and USE first.")
			}
			// Remove table from catalog
			err := (*cat).DropTable(s.Table)
			if err != nil {
				return fmt.Errorf("DROP TABLE failed: %w", err)
			}
			// Delete the table's data file
			tablePath := filepath.Join("data", *currentDB, s.Table+".db")
			if err := os.Remove(tablePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete table file: %w", err)
			}
			fmt.Printf("Table '%s' dropped.\n", s.Table)
			return nil
		}

		// DROP DATABASE
		if s.Database != "" {
			if *currentDB == s.Database {
				return fmt.Errorf("cannot drop the currently selected database ('%s'). Switch to another database first.", s.Database)
			}
			dbDir := filepath.Join("data", s.Database)
			if _, err := os.Stat(dbDir); os.IsNotExist(err) {
				return fmt.Errorf("database '%s' does not exist", s.Database)
			}
			// Remove the entire database directory and its contents
			if err := os.RemoveAll(dbDir); err != nil {
				return fmt.Errorf("failed to drop database '%s': %w", s.Database, err)
			}
			fmt.Printf("Database '%s' dropped.\n", s.Database)
			return nil
		}

		return fmt.Errorf("DROP: no table or database specified")
	case *par.SelectStatement:
		// SELECT
		if *currentDB == "" {
			return fmt.Errorf("no database selected. Use CREATE DATABASE and USE first.")
		}
		schema := (*cat).GetTable(s.Table)
		if schema == nil {
			return fmt.Errorf("table %q does not exist", s.Table)
		}
		// Validate columns (skip if SELECT *)
		if !(len(s.Columns) == 1 && s.Columns[0] == "*") {
			if !db.ColumnsExist(schema.Columns, s.Columns) {
				return fmt.Errorf("one or more selected columns do not exist in table %q", s.Table)
			}
		}
		tablePath := filepath.Join("data", *currentDB, s.Table+".db")
		pager := storage.NewPager(tablePath)
		rows := db.ReadAllRows(pager)
		// Print header
		fmt.Println(schema.Columns)
		for _, row := range rows {
			if s.Where != nil && !db.EvalWhere(s.Where, schema.Columns, row) { // Skip rows that don't match the WHERE condition
				continue
			}
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
		if *currentDB == "" {
			return fmt.Errorf("no database selected. Use CREATE DATABASE and USE first.")
		}
		schema := (*cat).GetTable(s.Table)
		if schema == nil {
			return fmt.Errorf("table %q does not exist", s.Table)
		}
		// Validate columns
		if !db.ColumnsMatch(schema.Columns, s.Columns) {
			return fmt.Errorf("column mismatch: expected %v, got %v", schema.Columns, s.Columns)
		}
		// Flatten [][]string to []string for storage
		var flatValues []string
		for _, v := range s.Values {
			flatValues = append(flatValues, v...)
		}
		tablePath := filepath.Join("data", *currentDB, s.Table+".db")
		pager := storage.NewPager(tablePath)
		_, err := db.InsertRow(pager, flatValues)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
		fmt.Println("Row inserted!")

	case *par.ShowDatabasesStatement:
		entries, err := os.ReadDir("data")
		if err != nil {
			return fmt.Errorf("failed to list databases: %w", err)
		}
		if len(entries) == 0 {
			fmt.Println("Databaes: No database found.\n Write CREATE DATABASE <database_name> to create one.")
		} else {
			fmt.Println("Databases:")
			for _, entry := range entries {
				if entry.IsDir() {
					fmt.Println(" -", entry.Name())
				}
			}
		}
	default:
		return fmt.Errorf("unsupported statement type")
	}
	return nil
}

// main entry point of the program
func main() {
	// Track the current database and catalog
	var currentDB string
	var cat *catalog.Catalog

	// Ensure data directory exists
	_ = os.MkdirAll("data", 0755)

	lineBuffer := repl.InitLineBuffer()
	printHelp()
	for {
		repl.PrintDB(currentDB)
		lineBuffer.UserInput()
		input := string(lineBuffer.Buffer)
		if input == "help;" {
			printHelp()
			lineBuffer.Reset()
			continue
		} else if input == "helpall;" {
			printHelpall()
			lineBuffer.Reset()
			continue
		} else if input == "\\e;" {
			println("Exiting letsgodb...")
			println("Bye!!")
			os.Exit(0)
			break
		}
		tokens := tok.Tokenizer(lineBuffer)
		// fmt.Println(tokens) // print the obtained tokens
		stmt := par.ParseProgram(tokens)
		if stmt == nil {
			// Print a friendly message for parse errors
			fmt.Println("Error: Failed to parse statement. Please check your SQL syntax.")
			lineBuffer.Reset()
			continue
		}
		err := ExecuteStatement(stmt, &currentDB, &cat)
		if err != nil {
			fmt.Println("Error:", err)
		}
		lineBuffer.Reset()
	}
}
