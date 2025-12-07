//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	dsn := "postgres://postgres:postgres@localhost:5432/myapp_postgres?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return pool
}

func TestIntegration_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// 1. Create (slug is auto-generated)
	t.Run("Create", func(t *testing.T) {
		org, err := repo.Create(ctx, "Test Organization")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if org.ID == "" {
			t.Error("Create: ID should not be empty")
		}
		if org.Name != "Test Organization" {
			t.Errorf("Create: Name = %s, want Test Organization", org.Name)
		}
		fmt.Printf("✓ Create: ID=%s, Name=%s, Slug=%s\n", org.ID, org.Name, org.Slug)

		// 2. GetByID
		t.Run("GetByID", func(t *testing.T) {
			fetched, err := repo.GetByID(ctx, org.ID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}
			if fetched.Name != org.Name {
				t.Errorf("GetByID: Name = %s, want %s", fetched.Name, org.Name)
			}
			fmt.Printf("✓ GetByID: ID=%s, Name=%s\n", fetched.ID, fetched.Name)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := repo.Update(ctx, org.ID, "Updated Organization", org.Slug)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Name != "Updated Organization" {
				t.Errorf("Update: Name = %s, want Updated Organization", updated.Name)
			}
			fmt.Printf("✓ Update: ID=%s, Name=%s, Slug=%s\n", updated.ID, updated.Name, updated.Slug)
		})

		// 4. List
		t.Run("List", func(t *testing.T) {
			orgs, err := repo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(orgs) == 0 {
				t.Error("List: should return at least one organization")
			}
			fmt.Printf("✓ List: returned %d organizations\n", len(orgs))
		})

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, org.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify soft delete
			_, err = repo.GetByID(ctx, org.ID)
			if err != ErrNotFound {
				t.Errorf("Delete: GetByID should return ErrNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: organization soft deleted\n")
		})
	})
}
