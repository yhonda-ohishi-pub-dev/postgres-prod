//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestIntegration_UserOrganizations_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	uoRepo := NewUserOrganizationRepository(pool)
	userRepo := NewAppUserRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Setup: Create test user and organization
	uniqueEmail := fmt.Sprintf("test-uo-user-%d@example.com", time.Now().UnixNano())
	testUser, err := userRepo.Create(ctx, &uniqueEmail, "Test UO User", nil, false)
	if err != nil {
		t.Fatalf("Setup: failed to create test user: %v", err)
	}
	defer userRepo.Delete(ctx, testUser.ID)

	// slug is auto-generated
	testOrg, err := orgRepo.Create(ctx, "Test UO Organization")
	if err != nil {
		t.Fatalf("Setup: failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		uo, err := uoRepo.Create(ctx, testUser.ID, testOrg.ID, "admin", true)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if uo.ID == "" {
			t.Error("Create: ID should not be empty")
		}
		if uo.UserID != testUser.ID {
			t.Errorf("Create: UserID = %s, want %s", uo.UserID, testUser.ID)
		}
		if uo.OrganizationID != testOrg.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", uo.OrganizationID, testOrg.ID)
		}
		if uo.Role != "admin" {
			t.Errorf("Create: Role = %s, want admin", uo.Role)
		}
		if uo.IsDefault != true {
			t.Errorf("Create: IsDefault = %v, want true", uo.IsDefault)
		}
		fmt.Printf("✓ Create: ID=%s, UserID=%s, OrganizationID=%s, Role=%s, IsDefault=%v\n", uo.ID, uo.UserID, uo.OrganizationID, uo.Role, uo.IsDefault)

		// 2. GetByID
		t.Run("GetByID", func(t *testing.T) {
			fetched, err := uoRepo.GetByID(ctx, uo.ID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}
			if fetched.UserID != uo.UserID {
				t.Errorf("GetByID: UserID = %s, want %s", fetched.UserID, uo.UserID)
			}
			if fetched.Role != uo.Role {
				t.Errorf("GetByID: Role = %s, want %s", fetched.Role, uo.Role)
			}
			fmt.Printf("✓ GetByID: ID=%s, Role=%s\n", fetched.ID, fetched.Role)
		})

		// 3. GetByUserAndOrganization
		t.Run("GetByUserAndOrganization", func(t *testing.T) {
			fetched, err := uoRepo.GetByUserAndOrganization(ctx, testUser.ID, testOrg.ID)
			if err != nil {
				t.Fatalf("GetByUserAndOrganization failed: %v", err)
			}
			if fetched.ID != uo.ID {
				t.Errorf("GetByUserAndOrganization: ID = %s, want %s", fetched.ID, uo.ID)
			}
			fmt.Printf("✓ GetByUserAndOrganization: ID=%s\n", fetched.ID)
		})

		// 4. ListByUserID
		t.Run("ListByUserID", func(t *testing.T) {
			uos, err := uoRepo.ListByUserID(ctx, testUser.ID)
			if err != nil {
				t.Fatalf("ListByUserID failed: %v", err)
			}
			if len(uos) == 0 {
				t.Error("ListByUserID: should return at least one record")
			}
			fmt.Printf("✓ ListByUserID: returned %d records\n", len(uos))
		})

		// 5. ListByOrganizationID
		t.Run("ListByOrganizationID", func(t *testing.T) {
			uos, err := uoRepo.ListByOrganizationID(ctx, testOrg.ID)
			if err != nil {
				t.Fatalf("ListByOrganizationID failed: %v", err)
			}
			if len(uos) == 0 {
				t.Error("ListByOrganizationID: should return at least one record")
			}
			fmt.Printf("✓ ListByOrganizationID: returned %d records\n", len(uos))
		})

		// 6. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := uoRepo.Update(ctx, uo.ID, "member", false)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Role != "member" {
				t.Errorf("Update: Role = %s, want member", updated.Role)
			}
			if updated.IsDefault != false {
				t.Errorf("Update: IsDefault = %v, want false", updated.IsDefault)
			}
			fmt.Printf("✓ Update: ID=%s, Role=%s, IsDefault=%v\n", updated.ID, updated.Role, updated.IsDefault)
		})

		// 7. List
		t.Run("List", func(t *testing.T) {
			uos, err := uoRepo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(uos) == 0 {
				t.Error("List: should return at least one record")
			}
			fmt.Printf("✓ List: returned %d records\n", len(uos))
		})

		// 8. Delete
		t.Run("Delete", func(t *testing.T) {
			err := uoRepo.Delete(ctx, uo.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify delete
			_, err = uoRepo.GetByID(ctx, uo.ID)
			if err != ErrUserOrganizationNotFound {
				t.Errorf("Delete: GetByID should return ErrUserOrganizationNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
