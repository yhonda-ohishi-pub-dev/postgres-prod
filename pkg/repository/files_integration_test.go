//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_Files_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewFileRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Setup: Create test organization
	uniqueSlug := fmt.Sprintf("test-files-%s", uuid.New().String()[:8])
	testOrg, err := orgRepo.Create(ctx, "Test Files Org", uniqueSlug)
	if err != nil {
		t.Fatalf("Setup: failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		blob := "test blob data"
		created := time.Now().Format(time.RFC3339)
		f, err := repo.Create(ctx, testOrg.ID, "test-file.txt", created, "text", &blob)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if f.UUID == "" {
			t.Error("Create: UUID should not be empty")
		}
		if f.OrganizationID != testOrg.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", f.OrganizationID, testOrg.ID)
		}
		if f.Filename != "test-file.txt" {
			t.Errorf("Create: Filename = %s, want test-file.txt", f.Filename)
		}
		fmt.Printf("✓ Create: UUID=%s, Filename=%s, Type=%s\n", f.UUID, f.Filename, f.Type)

		// 2. GetByUUID
		t.Run("GetByUUID", func(t *testing.T) {
			fetched, err := repo.GetByUUID(ctx, f.UUID)
			if err != nil {
				t.Fatalf("GetByUUID failed: %v", err)
			}
			if fetched.Filename != f.Filename {
				t.Errorf("GetByUUID: Filename = %s, want %s", fetched.Filename, f.Filename)
			}
			fmt.Printf("✓ GetByUUID: UUID=%s, Filename=%s\n", fetched.UUID, fetched.Filename)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			newBlob := "updated blob data"
			updated, err := repo.Update(ctx, f.UUID, "updated-file.txt", "pdf", &newBlob)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Filename != "updated-file.txt" {
				t.Errorf("Update: Filename = %s, want updated-file.txt", updated.Filename)
			}
			if updated.Type != "pdf" {
				t.Errorf("Update: Type = %s, want pdf", updated.Type)
			}
			fmt.Printf("✓ Update: UUID=%s, Filename=%s, Type=%s\n", updated.UUID, updated.Filename, updated.Type)
		})

		// 4. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			files, err := repo.ListByOrganization(ctx, testOrg.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(files) == 0 {
				t.Error("ListByOrganization: should return at least one file")
			}
			fmt.Printf("✓ ListByOrganization: returned %d files\n", len(files))
		})

		// 5. List
		t.Run("List", func(t *testing.T) {
			files, err := repo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(files) == 0 {
				t.Error("List: should return at least one file")
			}
			fmt.Printf("✓ List: returned %d files\n", len(files))
		})

		// 6. Delete
		t.Run("Delete", func(t *testing.T) {
			deletedAt := time.Now().Format(time.RFC3339)
			err := repo.Delete(ctx, f.UUID, deletedAt)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify soft delete
			_, err = repo.GetByUUID(ctx, f.UUID)
			if err != ErrFileNotFound {
				t.Errorf("Delete: GetByUUID should return ErrFileNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: file soft deleted\n")
		})
	})
}
