package converter

import (
	"database/sql"
	"time"

	"github.com/example/dopadone/internal/db"
	"github.com/example/dopadone/internal/domain"
)

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func DbAreaToDomain(dbArea db.Area) domain.Area {
	color, _ := domain.ParseColor(dbArea.Color.String)

	return domain.Area{
		ID:        dbArea.ID,
		Name:      dbArea.Name,
		Color:     color,
		SortOrder: int(dbArea.SortOrder),
		CreatedAt: dbArea.CreatedAt,
		UpdatedAt: dbArea.UpdatedAt,
		DeletedAt: dbArea.DeletedAt,
	}
}

func DbListAreasRowToDomain(row db.ListAreasRow) domain.Area {
	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}

func DbGetAreaByIDRowToDomain(row db.GetAreaByIDRow) domain.Area {
	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}

func DbCreateAreaRowToDomain(row db.CreateAreaRow) domain.Area {
	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}

func DbUpdateAreaRowToDomain(row db.UpdateAreaRow) domain.Area {
	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: row.DeletedAt,
	}
}

func DbSubareaToDomain(dbSubarea db.Subarea) domain.Subarea {
	color, _ := domain.ParseColor(dbSubarea.Color.String)

	return domain.Subarea{
		ID:        dbSubarea.ID,
		Name:      dbSubarea.Name,
		AreaID:    dbSubarea.AreaID,
		Color:     color,
		CreatedAt: dbSubarea.CreatedAt,
		UpdatedAt: dbSubarea.UpdatedAt,
		DeletedAt: dbSubarea.DeletedAt,
	}
}

func DbProjectToDomain(dbProject db.Project) domain.Project {
	var parentID *string
	if dbProject.ParentID.Valid {
		parentID = &dbProject.ParentID.String
	}

	var subareaID *string
	if dbProject.SubareaID.Valid {
		subareaID = &dbProject.SubareaID.String
	}

	status, _ := domain.ParseProjectStatus(dbProject.Status)
	priority, _ := domain.ParsePriority(dbProject.Priority)
	progress, _ := domain.ParseProgress(int(dbProject.Progress))
	color, _ := domain.ParseColor(dbProject.Color.String)

	return domain.Project{
		ID:          dbProject.ID,
		Name:        dbProject.Name,
		Description: nullStringToString(dbProject.Description),
		Goal:        nullStringToString(dbProject.Goal),
		Status:      status,
		Priority:    priority,
		Progress:    progress,
		StartDate:   nil,
		Deadline:    dbProject.Deadline,
		Color:       color,
		ParentID:    parentID,
		SubareaID:   subareaID,
		Position:    int(dbProject.Position),
		CreatedAt:   dbProject.CreatedAt,
		UpdatedAt:   dbProject.UpdatedAt,
		CompletedAt: dbProject.CompletedAt,
		DeletedAt:   dbProject.DeletedAt,
	}
}

func DbTaskToDomain(dbTask db.Task) domain.Task {
	status, _ := domain.ParseTaskStatus(dbTask.Status)
	priority, _ := domain.ParseTaskPriority(dbTask.Priority)

	var estimatedDuration domain.TaskDuration
	if dbTask.EstimatedDuration.Valid {
		estimatedDuration, _ = domain.ParseTaskDuration(int(dbTask.EstimatedDuration.Int64))
	}

	var isNext bool
	if dbTask.IsNext == 1 {
		isNext = true
	}

	return domain.Task{
		ID:                dbTask.ID,
		ProjectID:         dbTask.ProjectID,
		Title:             dbTask.Title,
		Description:       nullStringToString(dbTask.Description),
		StartDate:         dbTask.StartDate,
		Deadline:          dbTask.Deadline,
		Priority:          priority,
		Context:           nullStringToString(dbTask.Context),
		EstimatedDuration: estimatedDuration,
		Status:            status,
		IsNext:            isNext,
		CreatedAt:         dbTask.CreatedAt,
		UpdatedAt:         dbTask.UpdatedAt,
		DeletedAt:         dbTask.DeletedAt,
	}
}

func DbProjectRowToDomain(row db.ListProjectsBySubareaRecursiveRow) domain.Project {
	var parentID *string
	if row.ParentID.Valid {
		parentID = &row.ParentID.String
	}

	var subareaID *string
	if row.SubareaID.Valid {
		subareaID = &row.SubareaID.String
	}

	status, _ := domain.ParseProjectStatus(row.Status)
	priority, _ := domain.ParsePriority(row.Priority)
	progress, _ := domain.ParseProgress(int(row.Progress))
	color, _ := domain.ParseColor(row.Color.String)

	var deadline *time.Time
	if t, ok := row.Deadline.(*time.Time); ok {
		deadline = t
	}

	var completedAt *time.Time
	if t, ok := row.CompletedAt.(*time.Time); ok {
		completedAt = t
	}

	var deletedAt *time.Time
	if t, ok := row.DeletedAt.(*time.Time); ok {
		deletedAt = t
	}

	return domain.Project{
		ID:          row.ID,
		Name:        row.Name,
		Description: nullStringToString(row.Description),
		Goal:        nullStringToString(row.Goal),
		Status:      status,
		Priority:    priority,
		Progress:    progress,
		StartDate:   nil,
		Deadline:    deadline,
		Color:       color,
		ParentID:    parentID,
		SubareaID:   subareaID,
		Position:    int(row.Position),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		CompletedAt: completedAt,
		DeletedAt:   deletedAt,
	}
}
