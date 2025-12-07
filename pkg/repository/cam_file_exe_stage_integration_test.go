//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestCamFileExeStageIntegration_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileExeStageRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization first (slug is auto-generated)
	testOrg, err := orgRepo.Create(ctx, "Test Org for CamFileExeStage")
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		stage, err := repo.Create(ctx, 1, testOrg.ID, "Initial Stage")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if stage.Stage != 1 {
			t.Errorf("Create: Stage = %d, want 1", stage.Stage)
		}
		if stage.OrganizationID != testOrg.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", stage.OrganizationID, testOrg.ID)
		}
		if stage.Name != "Initial Stage" {
			t.Errorf("Create: Name = %s, want Initial Stage", stage.Name)
		}
		fmt.Printf("✓ Create: Stage=%d, OrgID=%s, Name=%s\n", stage.Stage, stage.OrganizationID, stage.Name)

		// 2. GetByStageAndOrg
		t.Run("GetByStageAndOrg", func(t *testing.T) {
			fetched, err := repo.GetByStageAndOrg(ctx, 1, testOrg.ID)
			if err != nil {
				t.Fatalf("GetByStageAndOrg failed: %v", err)
			}
			if fetched.Name != stage.Name {
				t.Errorf("GetByStageAndOrg: Name = %s, want %s", fetched.Name, stage.Name)
			}
			if fetched.Stage != stage.Stage {
				t.Errorf("GetByStageAndOrg: Stage = %d, want %d", fetched.Stage, stage.Stage)
			}
			fmt.Printf("✓ GetByStageAndOrg: Stage=%d, OrgID=%s, Name=%s\n", fetched.Stage, fetched.OrganizationID, fetched.Name)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := repo.Update(ctx, 1, testOrg.ID, "Updated Stage Name")
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Name != "Updated Stage Name" {
				t.Errorf("Update: Name = %s, want Updated Stage Name", updated.Name)
			}
			fmt.Printf("✓ Update: Stage=%d, OrgID=%s, Name=%s\n", updated.Stage, updated.OrganizationID, updated.Name)
		})

		// 4. Create additional stages for list test
		t.Run("CreateMultiple", func(t *testing.T) {
			_, err := repo.Create(ctx, 2, testOrg.ID, "Second Stage")
			if err != nil {
				t.Fatalf("Create second stage failed: %v", err)
			}
			_, err = repo.Create(ctx, 3, testOrg.ID, "Third Stage")
			if err != nil {
				t.Fatalf("Create third stage failed: %v", err)
			}
			fmt.Printf("✓ CreateMultiple: created stages 2 and 3\n")
		})

		// 5. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			stages, err := repo.ListByOrganization(ctx, testOrg.ID)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(stages) != 3 {
				t.Errorf("ListByOrganization: got %d stages, want 3", len(stages))
			}
			// Verify order (should be ASC by stage)
			if len(stages) == 3 {
				if stages[0].Stage != 1 || stages[1].Stage != 2 || stages[2].Stage != 3 {
					t.Errorf("ListByOrganization: stages not in correct order")
				}
			}
			fmt.Printf("✓ ListByOrganization: returned %d stages\n", len(stages))
		})

		// 6. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, 1, testOrg.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByStageAndOrg(ctx, 1, testOrg.ID)
			if err != ErrCamFileExeStageNotFound {
				t.Errorf("Delete: GetByStageAndOrg should return ErrCamFileExeStageNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: stage deleted\n")
		})

		// Clean up remaining stages
		repo.Delete(ctx, 2, testOrg.ID)
		repo.Delete(ctx, 3, testOrg.ID)
	})
}

func TestCamFileExeStageIntegration_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileExeStageRepository(pool)
	ctx := context.Background()

	t.Run("GetByStageAndOrg_NotFound", func(t *testing.T) {
		_, err := repo.GetByStageAndOrg(ctx, 999, "00000000-0000-0000-0000-000000000000")
		if err != ErrCamFileExeStageNotFound {
			t.Errorf("GetByStageAndOrg: expected ErrCamFileExeStageNotFound, got %v", err)
		}
		fmt.Printf("✓ GetByStageAndOrg_NotFound: returned correct error\n")
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		_, err := repo.Update(ctx, 999, "00000000-0000-0000-0000-000000000000", "Test")
		if err != ErrCamFileExeStageNotFound {
			t.Errorf("Update: expected ErrCamFileExeStageNotFound, got %v", err)
		}
		fmt.Printf("✓ Update_NotFound: returned correct error\n")
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		err := repo.Delete(ctx, 999, "00000000-0000-0000-0000-000000000000")
		if err != ErrCamFileExeStageNotFound {
			t.Errorf("Delete: expected ErrCamFileExeStageNotFound, got %v", err)
		}
		fmt.Printf("✓ Delete_NotFound: returned correct error\n")
	})
}

func TestCamFileExeStageIntegration_ForeignKeyConstraint(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileExeStageRepository(pool)
	ctx := context.Background()

	t.Run("Create_InvalidOrgID", func(t *testing.T) {
		_, err := repo.Create(ctx, 1, "00000000-0000-0000-0000-000000000000", "Test Stage")
		if err == nil {
			t.Error("Create: expected error for invalid organization_id, got nil")
		}
		fmt.Printf("✓ Create_InvalidOrgID: correctly rejected invalid foreign key\n")
	})
}
