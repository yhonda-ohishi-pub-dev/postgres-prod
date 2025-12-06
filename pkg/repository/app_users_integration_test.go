//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_AppUsers_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewAppUserRepository(pool)
	ctx := context.Background()

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		uniqueEmail := fmt.Sprintf("test-appuser-%s@example.com", uuid.New().String()[:8])
		user, err := repo.Create(ctx, uniqueEmail, "Test AppUser", false)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if user.ID == "" {
			t.Error("Create: ID should not be empty")
		}
		if user.IamEmail != uniqueEmail {
			t.Errorf("Create: IamEmail = %s, want %s", user.IamEmail, uniqueEmail)
		}
		if user.DisplayName != "Test AppUser" {
			t.Errorf("Create: DisplayName = %s, want Test AppUser", user.DisplayName)
		}
		if user.IsSuperadmin != false {
			t.Errorf("Create: IsSuperadmin = %v, want false", user.IsSuperadmin)
		}
		fmt.Printf("✓ Create: ID=%s, IamEmail=%s, DisplayName=%s, IsSuperadmin=%v\n", user.ID, user.IamEmail, user.DisplayName, user.IsSuperadmin)

		// 2. GetByID
		t.Run("GetByID", func(t *testing.T) {
			fetched, err := repo.GetByID(ctx, user.ID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}
			if fetched.IamEmail != user.IamEmail {
				t.Errorf("GetByID: IamEmail = %s, want %s", fetched.IamEmail, user.IamEmail)
			}
			if fetched.DisplayName != user.DisplayName {
				t.Errorf("GetByID: DisplayName = %s, want %s", fetched.DisplayName, user.DisplayName)
			}
			fmt.Printf("✓ GetByID: ID=%s, IamEmail=%s\n", fetched.ID, fetched.IamEmail)
		})

		// 3. GetByIamEmail
		t.Run("GetByIamEmail", func(t *testing.T) {
			fetched, err := repo.GetByIamEmail(ctx, user.IamEmail)
			if err != nil {
				t.Fatalf("GetByIamEmail failed: %v", err)
			}
			if fetched.ID != user.ID {
				t.Errorf("GetByIamEmail: ID = %s, want %s", fetched.ID, user.ID)
			}
			if fetched.DisplayName != user.DisplayName {
				t.Errorf("GetByIamEmail: DisplayName = %s, want %s", fetched.DisplayName, user.DisplayName)
			}
			fmt.Printf("✓ GetByIamEmail: ID=%s, IamEmail=%s\n", fetched.ID, fetched.IamEmail)
		})

		// 4. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := repo.Update(ctx, user.ID, "Updated AppUser", true)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.DisplayName != "Updated AppUser" {
				t.Errorf("Update: DisplayName = %s, want Updated AppUser", updated.DisplayName)
			}
			if updated.IsSuperadmin != true {
				t.Errorf("Update: IsSuperadmin = %v, want true", updated.IsSuperadmin)
			}
			if updated.IamEmail != user.IamEmail {
				t.Errorf("Update: IamEmail should not change, got %s", updated.IamEmail)
			}
			fmt.Printf("✓ Update: ID=%s, DisplayName=%s, IsSuperadmin=%v\n", updated.ID, updated.DisplayName, updated.IsSuperadmin)
		})

		// 5. List
		t.Run("List", func(t *testing.T) {
			users, err := repo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(users) == 0 {
				t.Error("List: should return at least one user")
			}
			fmt.Printf("✓ List: returned %d users\n", len(users))
		})

		// 6. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, user.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify soft delete
			_, err = repo.GetByID(ctx, user.ID)
			if err != ErrAppUserNotFound {
				t.Errorf("Delete: GetByID should return ErrAppUserNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: user soft deleted\n")
		})
	})
}
