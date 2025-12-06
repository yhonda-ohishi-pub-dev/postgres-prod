//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_CarInspectionFiles_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCarInspectionFilesRepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	typeVal := "inspection"
	electCertMgNo := fmt.Sprintf("cert-%s", uuid.New().String()[:8])
	electCertPublishdateE := "R"
	electCertPublishdateY := "05"
	electCertPublishdateM := "12"
	electCertPublishdateD := "07"
	now := time.Now().Format(time.RFC3339)

	var createdUUID string

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, orgID, typeVal, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, now, now)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.UUID == "" {
			t.Error("Create: UUID should not be empty")
		}
		if record.OrganizationID != orgID {
			t.Errorf("Create: OrganizationID = %s, want %s", record.OrganizationID, orgID)
		}
		createdUUID = record.UUID
		fmt.Printf("✓ Create: UUID=%s, OrgID=%s\n", record.UUID, record.OrganizationID)

		// 2. GetByUUID
		t.Run("GetByUUID", func(t *testing.T) {
			fetched, err := repo.GetByUUID(ctx, createdUUID)
			if err != nil {
				t.Fatalf("GetByUUID failed: %v", err)
			}
			if fetched.UUID != createdUUID {
				t.Errorf("GetByUUID: UUID = %s, want %s", fetched.UUID, createdUUID)
			}
			fmt.Printf("✓ GetByUUID: UUID=%s, Type=%s\n", fetched.UUID, fetched.Type)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updatedType := "updated-inspection"
			modifiedTime := time.Now().Format(time.RFC3339)
			updated, err := repo.Update(ctx, createdUUID, updatedType, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, modifiedTime)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Type != updatedType {
				t.Errorf("Update: Type = %s, want %s", updated.Type, updatedType)
			}
			fmt.Printf("✓ Update: Type=%s\n", updated.Type)
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
			err := repo.Delete(ctx, createdUUID, deletedTime)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByUUID(ctx, createdUUID)
			if err != ErrCarInspectionFileNotFound {
				t.Errorf("Delete: GetByUUID should return ErrCarInspectionFileNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record soft-deleted\n")
		})
	})
}
