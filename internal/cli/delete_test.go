package cli

import (
	"context"
	"errors"
	"testing"
)

// mockDeleter is a mock implementation of Deleteable for testing
type mockDeleter struct {
	getByIDFunc    func(ctx context.Context, id string) (any, error)
	softDeleteFunc func(ctx context.Context, id string) error
	hardDeleteFunc func(ctx context.Context, id string) error
}

func (m *mockDeleter) GetByID(ctx context.Context, id string) (any, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockDeleter) SoftDelete(ctx context.Context, id string) error {
	if m.softDeleteFunc != nil {
		return m.softDeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockDeleter) HardDelete(ctx context.Context, id string) error {
	if m.hardDeleteFunc != nil {
		return m.hardDeleteFunc(ctx, id)
	}
	return nil
}

func TestRunDelete(t *testing.T) {
	tests := []struct {
		name         string
		params       DeleteParams
		setupDeleter func(*mockDeleter)
		wantErr      bool
		errContains  string
	}{
		{
			name: "soft delete success",
			params: DeleteParams{
				ID:         "test-123",
				Permanent:  false,
				EntityName: "test",
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return "entity", nil
				}
				m.softDeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "hard delete success",
			params: DeleteParams{
				ID:         "test-456",
				Permanent:  true,
				EntityName: "test",
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return "entity", nil
				}
				m.hardDeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "entity not found with custom error",
			params: DeleteParams{
				ID:          "missing-123",
				Permanent:   false,
				EntityName:  "test",
				NotFoundErr: errors.New("not found"),
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "entity not found without custom error",
			params: DeleteParams{
				ID:         "missing-456",
				Permanent:  false,
				EntityName: "test",
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return nil, errors.New("random error")
				}
			},
			wantErr:     true,
			errContains: "failed to get test",
		},
		{
			name: "soft delete error",
			params: DeleteParams{
				ID:         "test-789",
				Permanent:  false,
				EntityName: "test",
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return "entity", nil
				}
				m.softDeleteFunc = func(ctx context.Context, id string) error {
					return errors.New("delete failed")
				}
			},
			wantErr:     true,
			errContains: "failed to delete test",
		},
		{
			name: "hard delete error",
			params: DeleteParams{
				ID:         "test-789",
				Permanent:  true,
				EntityName: "test",
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return "entity", nil
				}
				m.hardDeleteFunc = func(ctx context.Context, id string) error {
					return errors.New("delete failed")
				}
			},
			wantErr:     true,
			errContains: "failed to permanently delete test",
		},
		{
			name: "soft delete not found error",
			params: DeleteParams{
				ID:          "test-789",
				Permanent:   false,
				EntityName:  "test",
				NotFoundErr: errors.New("not found"),
			},
			setupDeleter: func(m *mockDeleter) {
				m.getByIDFunc = func(ctx context.Context, id string) (any, error) {
					return "entity", nil
				}
				m.softDeleteFunc = func(ctx context.Context, id string) error {
					return errors.New("not found")
				}
			},
			wantErr:     true,
			errContains: "failed to delete test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDeleter{}
			if tt.setupDeleter != nil {
				tt.setupDeleter(mock)
			}

			ctx := context.Background()
			err := RunDelete(ctx, mock, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errContains != "" && !errorsContains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func errorsContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
