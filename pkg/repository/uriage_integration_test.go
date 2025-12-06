//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_Uriage_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewUriageRepository(pool)
	ctx := context.Background()

	// Create unique test data
	name := fmt.Sprintf("test-name-%s", uuid.New().String()[:8])
	bumon := fmt.Sprintf("test-bumon-%s", uuid.New().String()[:8])
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	date := time.Now().Format("2006-01-02")
	kingaku := int32(50000)
	typeVal := int32(1)
	cam := int32(100)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, name, bumon, orgID, &kingaku, &typeVal, &cam, date)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.Bumon != bumon {
			t.Errorf("Create: Bumon = %s, want %s", record.Bumon, bumon)
		}
		if record.Name != name {
			t.Errorf("Create: Name = %s, want %s", record.Name, name)
		}
		fmt.Printf("✓ Create: Name=%s, Bumon=%s, OrgID=%s, Date=%s\n", record.Name, record.Bumon, record.OrganizationID, record.Date)

		// 2. GetByPrimaryKey
		t.Run("GetByPrimaryKey", func(t *testing.T) {
			fetched, err := repo.GetByPrimaryKey(ctx, name, bumon, date, orgID)
			if err != nil {
				t.Fatalf("GetByPrimaryKey failed: %v", err)
			}
			if fetched.Bumon != bumon {
				t.Errorf("GetByPrimaryKey: Bumon = %s, want %s", fetched.Bumon, bumon)
			}
			fmt.Printf("✓ GetByPrimaryKey: Name=%s, Bumon=%s, OrgID=%s\n", fetched.Name, fetched.Bumon, fetched.OrganizationID)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			newKingaku := int32(100000)
			newType := int32(2)
			newCam := int32(200)
			updated, err := repo.Update(ctx, name, bumon, date, orgID, &newKingaku, &newType, &newCam)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if *updated.Kingaku != newKingaku {
				t.Errorf("Update: Kingaku = %d, want %d", *updated.Kingaku, newKingaku)
			}
			fmt.Printf("✓ Update: Kingaku=%d, Type=%d\n", *updated.Kingaku, *updated.Type)
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
			err := repo.Delete(ctx, name, bumon, date, orgID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByPrimaryKey(ctx, name, bumon, date, orgID)
			if err != ErrUriageNotFound {
				t.Errorf("Delete: GetByPrimaryKey should return ErrUriageNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
