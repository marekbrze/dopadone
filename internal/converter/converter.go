package converter

import (
	"database/sql"
	"time"

	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func DbAreaToDomain(dbArea db.Area) domain.Area {
	var deletedAt *time.Time
	if dbArea.DeletedAt != nil {
		if t, ok := dbArea.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(dbArea.Color.String)

	return domain.Area{
		ID:        dbArea.ID,
		Name:      dbArea.Name,
		Color:     color,
		SortOrder: int(dbArea.SortOrder),
		CreatedAt: dbArea.CreatedAt,
		UpdatedAt: dbArea.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbListAreasRowToDomain(row db.ListAreasRow) domain.Area {
	var deletedAt *time.Time
	if row.DeletedAt != nil {
		if t, ok := row.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbGetAreaByIDRowToDomain(row db.GetAreaByIDRow) domain.Area {
	var deletedAt *time.Time
	if row.DeletedAt != nil {
		if t, ok := row.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbCreateAreaRowToDomain(row db.CreateAreaRow) domain.Area {
	var deletedAt *time.Time
	if row.DeletedAt != nil {
		if t, ok := row.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbUpdateAreaRowToDomain(row db.UpdateAreaRow) domain.Area {
	var deletedAt *time.Time
	if row.DeletedAt != nil {
		if t, ok := row.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(row.Color.String)

	return domain.Area{
		ID:        row.ID,
		Name:      row.Name,
		Color:     color,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbSubareaToDomain(dbSubarea db.Subarea) domain.Subarea {
	var deletedAt *time.Time
	if dbSubarea.DeletedAt != nil {
		if t, ok := dbSubarea.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	color, _ := domain.ParseColor(dbSubarea.Color.String)

	return domain.Subarea{
		ID:        dbSubarea.ID,
		Name:      dbSubarea.Name,
		AreaID:    dbSubarea.AreaID,
		Color:     color,
		CreatedAt: dbSubarea.CreatedAt,
		UpdatedAt: dbSubarea.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func DbProjectToDomain(dbProject db.Project) domain.Project {
	var deletedAt *time.Time
	if dbProject.DeletedAt != nil {
		if t, ok := dbProject.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	var completedAt *time.Time
	if dbProject.CompletedAt != nil {
		if t, ok := dbProject.CompletedAt.(time.Time); ok {
			completedAt = &t
		}
	}

	var deadline *time.Time
	if dbProject.Deadline != nil {
		if t, ok := dbProject.Deadline.(time.Time); ok {
			deadline = &t
		}
	}

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
		Deadline:    deadline,
		Color:       color,
		ParentID:    parentID,
		SubareaID:   subareaID,
		Position:    int(dbProject.Position),
		CreatedAt:   dbProject.CreatedAt,
		UpdatedAt:   dbProject.UpdatedAt,
		CompletedAt: completedAt,
		DeletedAt:   deletedAt,
	}
}

func DbTaskToDomain(dbTask db.Task) domain.Task {
	var deletedAt *time.Time
	if dbTask.DeletedAt != nil {
		if t, ok := dbTask.DeletedAt.(time.Time); ok {
			deletedAt = &t
		}
	}

	var startDate *time.Time
	if dbTask.StartDate != nil {
		if t, ok := dbTask.StartDate.(time.Time); ok {
			startDate = &t
		}
	}

	var deadline *time.Time
	if dbTask.Deadline != nil {
		if t, ok := dbTask.Deadline.(time.Time); ok {
			deadline = &t
		}
	}

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
		StartDate:         startDate,
		Deadline:          deadline,
		Priority:          priority,
		Context:           nullStringToString(dbTask.Context),
		EstimatedDuration: estimatedDuration,
		Status:            status,
		IsNext:            isNext,
		CreatedAt:         dbTask.CreatedAt,
		UpdatedAt:         dbTask.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}
