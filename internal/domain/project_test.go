package domain

import (
	"testing"
	"time"
)

func TestNewProject(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)

	t.Run("creates root project with subarea_id", func(t *testing.T) {
		subareaID := "subarea-123"
		project, err := NewProject(NewProjectParams{
			Name:      "Test Project",
			Status:    ProjectStatusActive,
			Priority:  PriorityHigh,
			Progress:  Progress(0),
			SubareaID: &subareaID,
			Position:  0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.Name != "Test Project" {
			t.Errorf("expected Name 'Test Project', got %s", project.Name)
		}
		if project.Status != ProjectStatusActive {
			t.Errorf("expected Status 'active', got %s", project.Status)
		}
		if project.Priority != PriorityHigh {
			t.Errorf("expected Priority 'high', got %s", project.Priority)
		}
		if project.SubareaID == nil || *project.SubareaID != subareaID {
			t.Errorf("expected SubareaID %s, got %v", subareaID, project.SubareaID)
		}
		if project.ParentID != nil {
			t.Errorf("expected ParentID to be nil for root project, got %v", project.ParentID)
		}
	})

	t.Run("creates nested project with parent_id", func(t *testing.T) {
		parentID := "parent-123"
		project, err := NewProject(NewProjectParams{
			Name:     "Nested Project",
			Status:   ProjectStatusActive,
			Priority: PriorityMedium,
			Progress: Progress(0),
			ParentID: &parentID,
			Position: 0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.ParentID == nil || *project.ParentID != parentID {
			t.Errorf("expected ParentID %s, got %v", parentID, project.ParentID)
		}
		if project.SubareaID != nil {
			t.Errorf("expected SubareaID to be nil for nested project, got %v", project.SubareaID)
		}
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "",
			Status:    ProjectStatusActive,
			Priority:  PriorityMedium,
			Progress:  Progress(0),
			SubareaID: &subareaID,
		})
		if err != ErrProjectNameEmpty {
			t.Errorf("expected ErrProjectNameEmpty, got %v", err)
		}
	})

	t.Run("returns error when status is invalid", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatus("invalid"),
			Priority:  PriorityMedium,
			Progress:  Progress(0),
			SubareaID: &subareaID,
		})
		if err != ErrProjectInvalidStatus {
			t.Errorf("expected ErrProjectInvalidStatus, got %v", err)
		}
	})

	t.Run("returns error when priority is invalid", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatusActive,
			Priority:  Priority("invalid"),
			Progress:  Progress(0),
			SubareaID: &subareaID,
		})
		if err != ErrProjectInvalidPriority {
			t.Errorf("expected ErrProjectInvalidPriority, got %v", err)
		}
	})

	t.Run("returns error when progress is out of range", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatusActive,
			Priority:  PriorityMedium,
			Progress:  Progress(150),
			SubareaID: &subareaID,
		})
		if err != ErrProjectInvalidProgress {
			t.Errorf("expected ErrProjectInvalidProgress, got %v", err)
		}
	})

	t.Run("returns error when both parent_id and subarea_id are nil", func(t *testing.T) {
		_, err := NewProject(NewProjectParams{
			Name:     "Test",
			Status:   ProjectStatusActive,
			Priority: PriorityMedium,
			Progress: Progress(0),
		})
		if err != ErrProjectNoParent {
			t.Errorf("expected ErrProjectNoParent, got %v", err)
		}
	})

	t.Run("returns error when deadline is before start date", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatusActive,
			Priority:  PriorityMedium,
			Progress:  Progress(0),
			StartDate: &future,
			Deadline:  &now,
			SubareaID: &subareaID,
		})
		if err != ErrProjectInvalidDateRange {
			t.Errorf("expected ErrProjectInvalidDateRange, got %v", err)
		}
	})

	t.Run("returns error when deadline is set without start date", func(t *testing.T) {
		subareaID := "subarea-123"
		_, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatusActive,
			Priority:  PriorityMedium,
			Progress:  Progress(0),
			Deadline:  &future,
			SubareaID: &subareaID,
		})
		if err != ErrProjectInvalidDateRange {
			t.Errorf("expected ErrProjectInvalidDateRange, got %v", err)
		}
	})

	t.Run("accepts valid date range", func(t *testing.T) {
		subareaID := "subarea-123"
		project, err := NewProject(NewProjectParams{
			Name:      "Test",
			Status:    ProjectStatusActive,
			Priority:  PriorityMedium,
			Progress:  Progress(0),
			StartDate: &now,
			Deadline:  &future,
			SubareaID: &subareaID,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.StartDate == nil || !project.StartDate.Equal(now) {
			t.Errorf("StartDate mismatch")
		}
		if project.Deadline == nil || !project.Deadline.Equal(future) {
			t.Errorf("Deadline mismatch")
		}
	})
}

func TestProjectIsDeleted(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	if project.IsDeleted() {
		t.Error("expected project to not be deleted initially")
	}

	now := time.Now()
	project.DeletedAt = &now
	if !project.IsDeleted() {
		t.Error("expected project to be deleted after setting DeletedAt")
	}
}

func TestProjectIsCompleted(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	if project.IsCompleted() {
		t.Error("expected project to not be completed initially")
	}

	project.Status = ProjectStatusCompleted
	if !project.IsCompleted() {
		t.Error("expected project to be completed after setting status")
	}
}

func TestProjectIsNested(t *testing.T) {
	subareaID := "subarea-123"
	rootProject, _ := NewProject(NewProjectParams{
		Name:      "Root",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	if rootProject.IsNested() {
		t.Error("expected root project to not be nested")
	}

	parentID := "parent-123"
	nestedProject, _ := NewProject(NewProjectParams{
		Name:     "Nested",
		Status:   ProjectStatusActive,
		Priority: PriorityMedium,
		Progress: Progress(0),
		ParentID: &parentID,
	})

	if !nestedProject.IsNested() {
		t.Error("expected nested project to be nested")
	}
}

func TestProjectMarkCompleted(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(50),
		SubareaID: &subareaID,
	})

	completedAt := time.Now()
	project.MarkCompleted(completedAt)

	if project.Status != ProjectStatusCompleted {
		t.Errorf("expected Status 'completed', got %s", project.Status)
	}
	if project.Progress != 100 {
		t.Errorf("expected Progress 100, got %d", project.Progress)
	}
	if project.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestProjectSetProgress(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	err := project.SetProgress(Progress(75))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Progress != 75 {
		t.Errorf("expected Progress 75, got %d", project.Progress)
	}

	err = project.SetProgress(Progress(150))
	if err != ErrProjectInvalidProgress {
		t.Errorf("expected ErrProjectInvalidProgress, got %v", err)
	}
}

func TestProjectSetPriority(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	err := project.SetPriority(PriorityUrgent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Priority != PriorityUrgent {
		t.Errorf("expected Priority 'urgent', got %s", project.Priority)
	}

	err = project.SetPriority(Priority("invalid"))
	if err != ErrProjectInvalidPriority {
		t.Errorf("expected ErrProjectInvalidPriority, got %v", err)
	}
}

func TestProjectSetStatus(t *testing.T) {
	subareaID := "subarea-123"
	project, _ := NewProject(NewProjectParams{
		Name:      "Test",
		Status:    ProjectStatusActive,
		Priority:  PriorityMedium,
		Progress:  Progress(0),
		SubareaID: &subareaID,
	})

	err := project.SetStatus(ProjectStatusOnHold)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Status != ProjectStatusOnHold {
		t.Errorf("expected Status 'on_hold', got %s", project.Status)
	}

	err = project.SetStatus(ProjectStatus("invalid"))
	if err != ErrProjectInvalidStatus {
		t.Errorf("expected ErrProjectInvalidStatus, got %v", err)
	}
}
