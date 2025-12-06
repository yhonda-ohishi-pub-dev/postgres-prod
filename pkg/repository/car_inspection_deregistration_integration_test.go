//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_CarInspectionDeregistration_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCarInspectionDeregistrationRepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	now := time.Now().Format(time.RFC3339)
	testData := &CarInspectionDeregistration{
		OrganizationID: orgID,
		Created:        now,
		Modified:       now,
	}

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, testData)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.UUID == "" {
			t.Error("Create: UUID should not be empty")
		}
		if record.OrganizationID != orgID {
			t.Errorf("Create: OrganizationID = %s, want %s", record.OrganizationID, orgID)
		}
		testData.UUID = record.UUID
		fmt.Printf("✓ Create: UUID=%s, OrgID=%s\n", record.UUID, record.OrganizationID)

		// 2. GetByUUID
		t.Run("GetByUUID", func(t *testing.T) {
			fetched, err := repo.GetByUUID(ctx, testData.UUID)
			if err != nil {
				t.Fatalf("GetByUUID failed: %v", err)
			}
			if fetched.UUID != testData.UUID {
				t.Errorf("GetByUUID: UUID = %s, want %s", fetched.UUID, testData.UUID)
			}
			fmt.Printf("✓ GetByUUID: UUID=%s, OrgID=%s\n", fetched.UUID, fetched.OrganizationID)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			testData.Modified = time.Now().Format(time.RFC3339)
			updated, err := repo.Update(ctx, testData)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.UUID != testData.UUID {
				t.Errorf("Update: UUID = %s, want %s", updated.UUID, testData.UUID)
			}
			fmt.Printf("✓ Update: UUID=%s\n", updated.UUID)
		})

		// 4. List
		t.Run("ListByOrganization", func(t *testing.T) {
			records, err := repo.ListByOrganization(ctx, orgID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(records) == 0 {
				t.Error("ListByOrganization: should return at least one record")
			}
			fmt.Printf("✓ ListByOrganization: returned %d records\n", len(records))
		})

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			deletedTime := time.Now().Format(time.RFC3339)
			err := repo.Delete(ctx, testData.UUID, deletedTime)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByUUID(ctx, testData.UUID)
			if err != ErrCarInspectionDeregistrationNotFound {
				t.Errorf("Delete: GetByUUID should return ErrCarInspectionDeregistrationNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record soft-deleted\n")
		})
	})
}
