package catalog

import (
	"os"
	"testing"
)

func TestCatalogBasicUsage(t *testing.T) {
	// Use a temporary file for testing
	testFile := "test_catalog.db"
	defer os.Remove(testFile)

	// Create a new catalog
	cat, err := NewCatalog(testFile)
	if err != nil {
		t.Fatalf("Failed to create catalog: %v", err)
	}

	// Add a table
	err = cat.AddTable("users", []string{"id", "name", "email"})
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	// Add another table
	err = cat.AddTable("posts", []string{"id", "user_id", "content"})
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	// Try to add duplicate table
	err = cat.AddTable("users", []string{"id", "name"})
	if err == nil {
		t.Errorf("Expected error when adding duplicate table, got nil")
	}

	// GetTable should return correct schema
	users := cat.GetTable("users")
	if users == nil {
		t.Fatalf("GetTable returned nil for 'users'")
	}
	if users.Name != "users" || len(users.Columns) != 3 || users.Columns[0] != "id" {
		t.Errorf("GetTable returned wrong schema for 'users': %+v", users)
	}

	posts := cat.GetTable("posts")
	if posts == nil {
		t.Fatalf("GetTable returned nil for 'posts'")
	}
	if posts.Name != "posts" || len(posts.Columns) != 3 || posts.Columns[2] != "content" {
		t.Errorf("GetTable returned wrong schema for 'posts': %+v", posts)
	}

	// ListTables should return both tables
	all := cat.ListTables()
	if len(all) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(all))
	}

	// Re-open catalog and check persistence
	cat2, err := NewCatalog(testFile)
	if err != nil {
		t.Fatalf("Failed to re-open catalog: %v", err)
	}
	users2 := cat2.GetTable("users")
	if users2 == nil || users2.Name != "users" {
		t.Errorf("Catalog did not persist 'users' table")
	}
	posts2 := cat2.GetTable("posts")
	if posts2 == nil || posts2.Name != "posts" {
		t.Errorf("Catalog did not persist 'posts' table")
	}
}
