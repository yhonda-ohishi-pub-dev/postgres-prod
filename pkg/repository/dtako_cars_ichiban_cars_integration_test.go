//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_DtakoCarsIchibanCars_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewDtakoCarsIchibanCarsRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create test organization
	uniqueSlug := fmt.Sprintf("dtako-cars-%s", uuid.New().String()[:8])
	org, err := orgRepo.Create(ctx, "Test Org for DtakoCars", uniqueSlug)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, org.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		idValue := "ichiban-123"
		entry, err := repo.Create(ctx, "dtako-001", org.ID, &idValue)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if entry.IdDtako != "dtako-001" {
			t.Errorf("Create: IdDtako = %s, want dtako-001", entry.IdDtako)
		}
		if entry.OrganizationID != org.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", entry.OrganizationID, org.ID)
		}
		if entry.Id == nil || *entry.Id != "ichiban-123" {
			t.Errorf("Create: Id = %v, want ichiban-123", entry.Id)
		}
		fmt.Printf("✓ Create: IdDtako=%s, OrganizationID=%s, Id=%v\n", entry.IdDtako, entry.OrganizationID, *entry.Id)

		// 2. GetByDtakoAndOrg
		t.Run("GetByDtakoAndOrg", func(t *testing.T) {
			fetched, err := repo.GetByDtakoAndOrg(ctx, entry.IdDtako, entry.OrganizationID)
			if err != nil {
				t.Fatalf("GetByDtakoAndOrg failed: %v", err)
			}
			if fetched.IdDtako != entry.IdDtako {
				t.Errorf("GetByDtakoAndOrg: IdDtako = %s, want %s", fetched.IdDtako, entry.IdDtako)
			}
			if fetched.OrganizationID != entry.OrganizationID {
				t.Errorf("GetByDtakoAndOrg: OrganizationID = %s, want %s", fetched.OrganizationID, entry.OrganizationID)
			}
			if fetched.Id == nil || *fetched.Id != *entry.Id {
				t.Errorf("GetByDtakoAndOrg: Id = %v, want %v", fetched.Id, *entry.Id)
			}
			fmt.Printf("✓ GetByDtakoAndOrg: IdDtako=%s, OrganizationID=%s, Id=%v\n", fetched.IdDtako, fetched.OrganizationID, *fetched.Id)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updatedIdValue := "ichiban-456"
			updated, err := repo.Update(ctx, entry.IdDtako, entry.OrganizationID, &updatedIdValue)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Id == nil || *updated.Id != "ichiban-456" {
				t.Errorf("Update: Id = %v, want ichiban-456", updated.Id)
			}
			if updated.IdDtako != entry.IdDtako {
				t.Errorf("Update: IdDtako should not change, got %s", updated.IdDtako)
			}
			if updated.OrganizationID != entry.OrganizationID {
				t.Errorf("Update: OrganizationID should not change, got %s", updated.OrganizationID)
			}
			fmt.Printf("✓ Update: IdDtako=%s, OrganizationID=%s, Id=%v\n", updated.IdDtako, updated.OrganizationID, *updated.Id)
		})

		// 4. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			entries, err := repo.ListByOrganization(ctx, org.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(entries) == 0 {
				t.Error("ListByOrganization: should return at least one entry")
			}
			fmt.Printf("✓ ListByOrganization: returned %d entries\n", len(entries))
		})

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, entry.IdDtako, entry.OrganizationID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByDtakoAndOrg(ctx, entry.IdDtako, entry.OrganizationID)
			if err != ErrDtakoCarsIchibanCarsNotFound {
				t.Errorf("Delete: GetByDtakoAndOrg should return ErrDtakoCarsIchibanCarsNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: entry deleted\n")
		})
	})
}

func TestIntegration_DtakoCarsIchibanCars_CreateWithNullId(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewDtakoCarsIchibanCarsRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create test organization
	nullSlug := fmt.Sprintf("dtako-null-%s", uuid.New().String()[:8])
	org, err := orgRepo.Create(ctx, "Test Org for DtakoCars Null", nullSlug)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, org.ID)

	// Create entry with null id
	entry, err := repo.Create(ctx, "dtako-002", org.ID, nil)
	if err != nil {
		t.Fatalf("Create with null id failed: %v", err)
	}
	defer repo.Delete(ctx, entry.IdDtako, entry.OrganizationID)

	if entry.Id != nil {
		t.Errorf("Create with null: Id should be nil, got %v", entry.Id)
	}
	fmt.Printf("✓ Create with null Id: IdDtako=%s, OrganizationID=%s, Id=nil\n", entry.IdDtako, entry.OrganizationID)

	// Fetch and verify
	fetched, err := repo.GetByDtakoAndOrg(ctx, entry.IdDtako, entry.OrganizationID)
	if err != nil {
		t.Fatalf("GetByDtakoAndOrg failed: %v", err)
	}
	if fetched.Id != nil {
		t.Errorf("Fetched entry: Id should be nil, got %v", fetched.Id)
	}
	fmt.Printf("✓ GetByDtakoAndOrg with null Id: verified Id=nil\n")
}
