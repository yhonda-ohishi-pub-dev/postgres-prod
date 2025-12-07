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
		avatarURL := "https://example.com/avatar.png"
		user, err := repo.Create(ctx, &uniqueEmail, "Test AppUser", &avatarURL, false)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if user.ID == "" {
			t.Error("Create: ID should not be empty")
		}
		if user.Email == nil || *user.Email != uniqueEmail {
			t.Errorf("Create: Email = %v, want %s", user.Email, uniqueEmail)
		}
		if user.DisplayName != "Test AppUser" {
			t.Errorf("Create: DisplayName = %s, want Test AppUser", user.DisplayName)
		}
		if user.AvatarURL == nil || *user.AvatarURL != avatarURL {
			t.Errorf("Create: AvatarURL = %v, want %s", user.AvatarURL, avatarURL)
		}
		if user.IsSuperadmin != false {
			t.Errorf("Create: IsSuperadmin = %v, want false", user.IsSuperadmin)
		}
		fmt.Printf("✓ Create: ID=%s, Email=%s, DisplayName=%s, IsSuperadmin=%v\n", user.ID, *user.Email, user.DisplayName, user.IsSuperadmin)

		// 2. GetByID
		t.Run("GetByID", func(t *testing.T) {
			fetched, err := repo.GetByID(ctx, user.ID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}
			if fetched.Email == nil || *fetched.Email != *user.Email {
				t.Errorf("GetByID: Email = %v, want %s", fetched.Email, *user.Email)
			}
			if fetched.DisplayName != user.DisplayName {
				t.Errorf("GetByID: DisplayName = %s, want %s", fetched.DisplayName, user.DisplayName)
			}
			fmt.Printf("✓ GetByID: ID=%s, Email=%s\n", fetched.ID, *fetched.Email)
		})

		// 3. GetByEmail
		t.Run("GetByEmail", func(t *testing.T) {
			fetched, err := repo.GetByEmail(ctx, *user.Email)
			if err != nil {
				t.Fatalf("GetByEmail failed: %v", err)
			}
			if fetched.ID != user.ID {
				t.Errorf("GetByEmail: ID = %s, want %s", fetched.ID, user.ID)
			}
			if fetched.DisplayName != user.DisplayName {
				t.Errorf("GetByEmail: DisplayName = %s, want %s", fetched.DisplayName, user.DisplayName)
			}
			fmt.Printf("✓ GetByEmail: ID=%s, Email=%s\n", fetched.ID, *fetched.Email)
		})

		// 4. Update
		t.Run("Update", func(t *testing.T) {
			newAvatarURL := "https://example.com/new-avatar.png"
			updated, err := repo.Update(ctx, user.ID, "Updated AppUser", &newAvatarURL, true)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.DisplayName != "Updated AppUser" {
				t.Errorf("Update: DisplayName = %s, want Updated AppUser", updated.DisplayName)
			}
			if updated.IsSuperadmin != true {
				t.Errorf("Update: IsSuperadmin = %v, want true", updated.IsSuperadmin)
			}
			if updated.AvatarURL == nil || *updated.AvatarURL != newAvatarURL {
				t.Errorf("Update: AvatarURL = %v, want %s", updated.AvatarURL, newAvatarURL)
			}
			if updated.Email == nil || *updated.Email != *user.Email {
				t.Errorf("Update: Email should not change, got %v", updated.Email)
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

	// Test with nil email (LINE case)
	t.Run("CreateWithNilEmail", func(t *testing.T) {
		user, err := repo.Create(ctx, nil, "LINE User", nil, false)
		if err != nil {
			t.Fatalf("Create with nil email failed: %v", err)
		}
		if user.Email != nil {
			t.Errorf("Create: Email should be nil, got %v", user.Email)
		}
		fmt.Printf("✓ CreateWithNilEmail: ID=%s, Email=nil, DisplayName=%s\n", user.ID, user.DisplayName)

		// Cleanup
		_ = repo.Delete(ctx, user.ID)
	})
}
