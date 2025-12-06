package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MockRow implements pgx.Row for testing
type MockRow struct {
	scanFunc func(dest ...any) error
}

func (r *MockRow) Scan(dest ...any) error {
	return r.scanFunc(dest...)
}

// MockRows implements pgx.Rows for testing
type MockRows struct {
	data    [][]any
	current int
	closed  bool
	err     error
}

func (r *MockRows) Close()                        { r.closed = true }
func (r *MockRows) Err() error                    { return r.err }
func (r *MockRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
func (r *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}
func (r *MockRows) Next() bool {
	if r.current >= len(r.data) {
		return false
	}
	r.current++
	return true
}
func (r *MockRows) Scan(dest ...any) error {
	if r.current == 0 || r.current > len(r.data) {
		return errors.New("no row")
	}
	row := r.data[r.current-1]
	for i, d := range dest {
		switch v := d.(type) {
		case *string:
			*v = row[i].(string)
		case *time.Time:
			*v = row[i].(time.Time)
		case **time.Time:
			if row[i] == nil {
				*v = nil
			} else {
				t := row[i].(time.Time)
				*v = &t
			}
		}
	}
	return nil
}
func (r *MockRows) Values() ([]any, error) { return nil, nil }
func (r *MockRows) RawValues() [][]byte    { return nil }
func (r *MockRows) Conn() *pgx.Conn        { return nil }

// MockDB implements DB interface for testing
type MockDB struct {
	queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
	queryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	execFunc     func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.queryRowFunc(ctx, sql, args...)
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return m.queryFunc(ctx, sql, args...)
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return m.execFunc(ctx, sql, args...)
}

func TestCreate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		input   struct{ name, slug string }
		mockRow *MockRow
		wantErr bool
	}{
		{
			name:  "success",
			input: struct{ name, slug string }{"Test Org", "test-org"},
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					*dest[0].(*string) = "uuid-123"
					*dest[1].(*string) = "Test Org"
					*dest[2].(*string) = "test-org"
					*dest[3].(*time.Time) = now
					*dest[4].(*time.Time) = now
					*dest[5].(**time.Time) = nil
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:  "db error",
			input: struct{ name, slug string }{"Test", "test"},
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					return errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				queryRowFunc: func(ctx context.Context, sql string, args ...any) pgx.Row {
					return tt.mockRow
				},
			}

			repo := NewOrganizationRepositoryWithDB(mockDB)
			org, err := repo.Create(context.Background(), tt.input.name, tt.input.slug)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && org == nil {
				t.Error("Create() returned nil organization")
			}

			if !tt.wantErr && org.Name != tt.input.name {
				t.Errorf("Create() name = %v, want %v", org.Name, tt.input.name)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		id      string
		mockRow *MockRow
		wantErr error
	}{
		{
			name: "found",
			id:   "uuid-123",
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					*dest[0].(*string) = "uuid-123"
					*dest[1].(*string) = "Test Org"
					*dest[2].(*string) = "test-org"
					*dest[3].(*time.Time) = now
					*dest[4].(*time.Time) = now
					*dest[5].(**time.Time) = nil
					return nil
				},
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   "nonexistent",
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					return pgx.ErrNoRows
				},
			},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				queryRowFunc: func(ctx context.Context, sql string, args ...any) pgx.Row {
					return tt.mockRow
				},
			}

			repo := NewOrganizationRepositoryWithDB(mockDB)
			org, err := repo.GetByID(context.Background(), tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && org == nil {
				t.Error("GetByID() returned nil organization")
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		input   struct{ id, name, slug string }
		mockRow *MockRow
		wantErr error
	}{
		{
			name:  "success",
			input: struct{ id, name, slug string }{"uuid-123", "Updated Org", "updated-org"},
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					*dest[0].(*string) = "uuid-123"
					*dest[1].(*string) = "Updated Org"
					*dest[2].(*string) = "updated-org"
					*dest[3].(*time.Time) = now
					*dest[4].(*time.Time) = now
					*dest[5].(**time.Time) = nil
					return nil
				},
			},
			wantErr: nil,
		},
		{
			name:  "not found",
			input: struct{ id, name, slug string }{"nonexistent", "Name", "slug"},
			mockRow: &MockRow{
				scanFunc: func(dest ...any) error {
					return pgx.ErrNoRows
				},
			},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				queryRowFunc: func(ctx context.Context, sql string, args ...any) pgx.Row {
					return tt.mockRow
				},
			}

			repo := NewOrganizationRepositoryWithDB(mockDB)
			org, err := repo.Update(context.Background(), tt.input.id, tt.input.name, tt.input.slug)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && org.Name != tt.input.name {
				t.Errorf("Update() name = %v, want %v", org.Name, tt.input.name)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		rowsAffected int64
		execErr     error
		wantErr     error
	}{
		{
			name:        "success",
			id:          "uuid-123",
			rowsAffected: 1,
			execErr:     nil,
			wantErr:     nil,
		},
		{
			name:        "not found",
			id:          "nonexistent",
			rowsAffected: 0,
			execErr:     nil,
			wantErr:     ErrNotFound,
		},
		{
			name:        "db error",
			id:          "uuid-123",
			rowsAffected: 0,
			execErr:     errors.New("db error"),
			wantErr:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				execFunc: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
					if tt.execErr != nil {
						return pgconn.CommandTag{}, tt.execErr
					}
					return pgconn.NewCommandTag("UPDATE " + string(rune('0'+tt.rowsAffected))), nil
				},
			}

			repo := NewOrganizationRepositoryWithDB(mockDB)
			err := repo.Delete(context.Background(), tt.id)

			if tt.wantErr == nil && err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
			}
			if tt.wantErr != nil && err == nil {
				t.Errorf("Delete() expected error %v, got nil", tt.wantErr)
			}
		})
	}
}

func TestList(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		limit    int
		offset   int
		mockRows *MockRows
		wantLen  int
		wantErr  bool
	}{
		{
			name:   "success with results",
			limit:  10,
			offset: 0,
			mockRows: &MockRows{
				data: [][]any{
					{"uuid-1", "Org 1", "org-1", now, now, nil},
					{"uuid-2", "Org 2", "org-2", now, now, nil},
				},
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:   "empty results",
			limit:  10,
			offset: 0,
			mockRows: &MockRows{
				data: [][]any{},
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:   "limit capped at 100",
			limit:  200,
			offset: 0,
			mockRows: &MockRows{
				data: [][]any{},
			},
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				queryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
					return tt.mockRows, nil
				},
			}

			repo := NewOrganizationRepositoryWithDB(mockDB)
			orgs, err := repo.List(context.Background(), tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(orgs) != tt.wantLen {
				t.Errorf("List() returned %d orgs, want %d", len(orgs), tt.wantLen)
			}
		})
	}
}
