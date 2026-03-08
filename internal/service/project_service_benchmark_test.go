package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db"
)

func BenchmarkListBySubareaRecursive(b *testing.B) {
	now := time.Now()
	subareaID := "subarea-benchmark"

	sizes := []int{100, 500, 1000}

	for _, size := range sizes {
		b.Run("new_implementation_"+string(rune(size)), func(b *testing.B) {
			mock := &mockProjectQuerier{
				listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
					filterRatio := 0.1
					returnCount := int(float64(size) * filterRatio)
					rows := make([]db.ListProjectsBySubareaRecursiveRow, returnCount)
					for i := 0; i < returnCount; i++ {
						rows[i] = projectToRow(db.Project{
							ID:        string(rune(i)),
							SubareaID: subareaID,
							Status:    "active",
							Priority:  "high",
							Progress:  0,
							CreatedAt: now,
							UpdatedAt: now,
						})
					}
					return rows, nil
				},
			}

			svc := NewProjectService(mock, nil)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = svc.ListBySubareaRecursive(context.Background(), subareaID)
			}
		})
	}
}
