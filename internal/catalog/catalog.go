package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// TableSchema represents the schema of a table (name and columns).
type TableSchema struct {
	Name       string   `json:"name"`
	Columns    []string `json:"columns"`
	PrimaryKey string   `Json:"primary_key"` // <- Add this
}

// Catalog manages table schemas and persists them to a catalog file.
type Catalog struct {
	filename string
	mu       sync.Mutex
	tables   map[string]*TableSchema
}

// NewCatalog creates a new Catalog instance and loads existing schemas from file.
func NewCatalog(filename string) (*Catalog, error) {
	c := &Catalog{
		filename: filename,
		tables:   make(map[string]*TableSchema),
	}
	if err := c.load(); err != nil {
		return nil, err
	}
	return c, nil
}

// load reads all table schemas from the catalog file.
func (c *Catalog) load() error {
	file, err := os.OpenFile(c.filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open catalog file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var schema TableSchema
		if err := json.Unmarshal(scanner.Bytes(), &schema); err != nil {
			return fmt.Errorf("failed to parse catalog entry: %w", err)
		}
		c.tables[schema.Name] = &schema
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading catalog file: %w", err)
	}
	return nil
}

// AddTable adds a new table schema to the catalog and persists it.
func (c *Catalog) AddTable(name string, columns []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.tables[name]; exists {
		return fmt.Errorf("table %q already exists", name)
	}
	schema := &TableSchema{
		Name:       name,
		Columns:    columns,
		PrimaryKey: columns[0],
	}
	// Append to file
	file, err := os.OpenFile(c.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open catalog file for writing: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write schema to catalog: %w", err)
	}

	c.tables[name] = schema
	return nil
}

// GetTable returns the schema for a given table name, or nil if not found.
func (c *Catalog) GetTable(name string) *TableSchema {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.tables[name]
}

// ListTables returns all table schemas.
func (c *Catalog) ListTables() []*TableSchema {
	c.mu.Lock()
	defer c.mu.Unlock()
	schemas := make([]*TableSchema, 0, len(c.tables))
	for _, schema := range c.tables {
		schemas = append(schemas, schema)
	}
	if len(schemas) == 0 {
		return nil
	}
	return schemas
}

// DropTable removes a table schema from the catalog and updates the catalog file.
func (c *Catalog) DropTable(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.tables[name]; !exists {
		return fmt.Errorf("table %q does not exist", name)
	}
	// Remove from in-memory map
	delete(c.tables, name)

	// Rewrite the catalog file with all remaining tables
	file, err := os.OpenFile(c.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open catalog file for rewriting: %w", err)
	}
	defer file.Close()

	for _, schema := range c.tables {
		data, err := json.Marshal(schema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema: %w", err)
		}
		if _, err := file.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("failed to write schema to catalog: %w", err)
		}
	}
	return nil
}
