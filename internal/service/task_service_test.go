package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

type mockTaskQuerier struct {
	createTaskFunc          func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
	getTaskByIDFunc         func(ctx context.Context, id string) (db.Task, error)
	listTasksByProjectFunc  func(ctx context.Context, projectID string) ([]db.Task, error)
	listNextTasksFunc       func(ctx context.Context) ([]db.Task, error)
	listTasksByStatusFunc   func(ctx context.Context, status string) ([]db.Task, error)
	listTasksByPriorityFunc func(ctx context.Context, priority string) ([]db.Task, error)
	listAllTasksFunc        func(ctx context.Context) ([]db.Task, error)
	updateTaskFunc          func(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error)
	softDeleteTaskFunc      func(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error)
	hardDeleteTaskFunc      func(ctx context.Context, id string) error
	toggleIsNextFunc        func(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error)
}

func (m *mockTaskQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, arg)
	}
	return db.Task{}, nil
}

func (m *mockTaskQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
	if m.getTaskByIDFunc != nil {
		return m.getTaskByIDFunc(ctx, id)
	}
	return db.Task{}, nil
}

func (m *mockTaskQuerier) ListTasksByProject(ctx context.Context, projectID string) ([]db.Task, error) {
	if m.listTasksByProjectFunc != nil {
		return m.listTasksByProjectFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *mockTaskQuerier) ListNextTasks(ctx context.Context) ([]db.Task, error) {
	if m.listNextTasksFunc != nil {
		return m.listNextTasksFunc(ctx)
	}
	return nil, nil
}

func (m *mockTaskQuerier) ListTasksByStatus(ctx context.Context, status string) ([]db.Task, error) {
	if m.listTasksByStatusFunc != nil {
		return m.listTasksByStatusFunc(ctx, status)
	}
	return nil, nil
}

func (m *mockTaskQuerier) ListTasksByPriority(ctx context.Context, priority string) ([]db.Task, error) {
	if m.listTasksByPriorityFunc != nil {
		return m.listTasksByPriorityFunc(ctx, priority)
	}
	return nil, nil
}

func (m *mockTaskQuerier) ListAllTasks(ctx context.Context) ([]db.Task, error) {
	if m.listAllTasksFunc != nil {
		return m.listAllTasksFunc(ctx)
	}
	return nil, nil
}

func (m *mockTaskQuerier) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	if m.updateTaskFunc != nil {
		return m.updateTaskFunc(ctx, arg)
	}
	return db.Task{}, nil
}

func (m *mockTaskQuerier) SoftDeleteTask(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error) {
	if m.softDeleteTaskFunc != nil {
		return m.softDeleteTaskFunc(ctx, arg)
	}
	return db.Task{}, nil
}

func (m *mockTaskQuerier) HardDeleteTask(ctx context.Context, id string) error {
	if m.hardDeleteTaskFunc != nil {
		return m.hardDeleteTaskFunc(ctx, id)
	}
	return nil
}

func (m *mockTaskQuerier) ToggleIsNext(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
	if m.toggleIsNextFunc != nil {
		return m.toggleIsNextFunc(ctx, arg)
	}
	return db.Task{}, nil
}

func (m *mockTaskQuerier) ListAllSubareas(ctx context.Context) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockTaskQuerier) ListProjectsByPriority(ctx context.Context, priority string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockTaskQuerier) CreateArea(ctx context.Context, arg db.CreateAreaParams) (db.CreateAreaRow, error) {
	return db.CreateAreaRow{}, nil
}

func (m *mockTaskQuerier) GetAreaByID(ctx context.Context, id string) (db.GetAreaByIDRow, error) {
	return db.GetAreaByIDRow{}, nil
}

func (m *mockTaskQuerier) ListAreas(ctx context.Context) ([]db.ListAreasRow, error) {
	return nil, nil
}

func (m *mockTaskQuerier) UpdateArea(ctx context.Context, arg db.UpdateAreaParams) (db.UpdateAreaRow, error) {
	return db.UpdateAreaRow{}, nil
}

func (m *mockTaskQuerier) UpdateAreaSortOrder(ctx context.Context, arg db.UpdateAreaSortOrderParams) error {
	return nil
}

func (m *mockTaskQuerier) SoftDeleteArea(ctx context.Context, arg db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
	return db.SoftDeleteAreaRow{}, nil
}

func (m *mockTaskQuerier) HardDeleteArea(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskQuerier) CountSubareasByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) CountProjectsByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) CountTasksByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) CountProjectsByParent(ctx context.Context, parentID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) CountProjectsBySubarea(ctx context.Context, subareaID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) CountTasksByProject(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}

func (m *mockTaskQuerier) DeleteTasksByProject(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockTaskQuerier) DeleteProjectsBySubarea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockTaskQuerier) DeleteSubareasByArea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockTaskQuerier) CreateSubarea(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockTaskQuerier) GetSubareaByID(ctx context.Context, id string) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockTaskQuerier) ListSubareasByArea(ctx context.Context, areaID string) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockTaskQuerier) UpdateSubarea(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockTaskQuerier) SoftDeleteSubarea(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockTaskQuerier) HardDeleteSubarea(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskQuerier) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockTaskQuerier) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockTaskQuerier) ListProjectsBySubarea(ctx context.Context, subareaID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockTaskQuerier) ListProjectsByParent(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockTaskQuerier) ListAllProjects(ctx context.Context) ([]db.Project, error) {
	return nil, nil
}

func (m *mockTaskQuerier) GetProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockTaskQuerier) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockTaskQuerier) SoftDeleteProject(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockTaskQuerier) HardDeleteProject(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskQuerier) DeleteProjectsByParentID(ctx context.Context, parentID sql.NullString) error {
	return nil
}

func (m *mockTaskQuerier) DeleteProjectsBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockTaskQuerier) DeleteTasksBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockTaskQuerier) DeleteTasksByProjectID(ctx context.Context, projectID string) error {
	return nil
}

func TestTaskService_Create(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		params  CreateTaskParams
		mock    func() *mockTaskQuerier
		wantErr bool
	}{
		{
			name: "creates task successfully",
			params: CreateTaskParams{
				ProjectID: "project-1",
				Title:     "Test Task",
				Status:    domain.TaskStatusTodo,
				Priority:  domain.TaskPriorityHigh,
				IsNext:    false,
			},
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					createTaskFunc: func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
						return db.Task{
							ID:        arg.ID,
							ProjectID: arg.ProjectID,
							Title:     arg.Title,
							Status:    arg.Status,
							Priority:  arg.Priority,
							IsNext:    arg.IsNext,
							CreatedAt: now,
							UpdatedAt: now,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "rejects empty title",
			params: CreateTaskParams{
				ProjectID: "project-1",
				Title:     "",
				Status:    domain.TaskStatusTodo,
				Priority:  domain.TaskPriorityHigh,
			},
			mock:    func() *mockTaskQuerier { return &mockTaskQuerier{} },
			wantErr: true,
		},
		{
			name: "rejects empty project ID",
			params: CreateTaskParams{
				ProjectID: "",
				Title:     "Test",
				Status:    domain.TaskStatusTodo,
				Priority:  domain.TaskPriorityHigh,
			},
			mock:    func() *mockTaskQuerier { return &mockTaskQuerier{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.Create(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TaskService.Create() returned nil task")
			}
		})
	}
}

func TestTaskService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockTaskQuerier
		wantErr bool
	}{
		{
			name: "retrieves task by ID",
			id:   "task-1",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						if id == "task-1" {
							return db.Task{
								ID:        "task-1",
								ProjectID: "project-1",
								Title:     "Test Task",
								Status:    "todo",
								Priority:  "high",
								IsNext:    0,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							}, nil
						}
						return db.Task{}, sql.ErrNoRows
					},
				}
			},
			wantErr: false,
		},
		{
			name: "returns error for non-existent task",
			id:   "nonexistent",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{}, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.GetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TaskService.GetByID() returned nil task")
			}
		})
	}
}

func TestTaskService_ListByProject(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		mock      func() *mockTaskQuerier
		want      int
		wantErr   bool
	}{
		{
			name:      "lists tasks by project",
			projectID: "project-1",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					listTasksByProjectFunc: func(ctx context.Context, projectID string) ([]db.Task, error) {
						return []db.Task{
							{ID: "task-1", ProjectID: "project-1", Title: "Task 1", Status: "todo", Priority: "high", IsNext: 0},
							{ID: "task-2", ProjectID: "project-1", Title: "Task 2", Status: "in_progress", Priority: "medium", IsNext: 1},
						}, nil
					},
				}
			},
			want:    2,
			wantErr: false,
		},
		{
			name:      "returns empty list",
			projectID: "project-2",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					listTasksByProjectFunc: func(ctx context.Context, projectID string) ([]db.Task, error) {
						return []db.Task{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.ListByProject(context.Background(), tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.ListByProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("TaskService.ListByProject() returned %d tasks, want %d", len(got), tt.want)
			}
		})
	}
}

func TestTaskService_SetStatus(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		status  domain.TaskStatus
		mock    func() *mockTaskQuerier
		wantErr bool
	}{
		{
			name:   "sets task status",
			id:     "task-1",
			status: domain.TaskStatusDone,
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{
							ID:        "task-1",
							ProjectID: "project-1",
							Title:     "Test Task",
							Status:    "todo",
							Priority:  "high",
							IsNext:    0,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
					updateTaskFunc: func(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
						return db.Task{
							ID:        arg.ID,
							ProjectID: "project-1",
							Title:     arg.Title,
							Status:    arg.Status,
							Priority:  arg.Priority,
							IsNext:    arg.IsNext,
							UpdatedAt: arg.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "returns error for non-existent task",
			id:     "nonexistent",
			status: domain.TaskStatusDone,
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{}, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.SetStatus(context.Background(), tt.id, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.SetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TaskService.SetStatus() returned nil task")
			}
		})
	}
}

func TestTaskService_ToggleIsNext(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockTaskQuerier
		wantErr bool
	}{
		{
			name: "toggles is_next from false to true",
			id:   "task-1",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{
							ID:        "task-1",
							ProjectID: "project-1",
							Title:     "Test Task",
							Status:    "todo",
							Priority:  "high",
							IsNext:    0,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
					toggleIsNextFunc: func(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
						return db.Task{
							ID:        arg.ID,
							ProjectID: "project-1",
							Title:     "Test Task",
							Status:    "todo",
							Priority:  "high",
							IsNext:    1,
							UpdatedAt: arg.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "returns error for non-existent task",
			id:   "nonexistent",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{}, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.ToggleIsNext(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.ToggleIsNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TaskService.ToggleIsNext() returned nil task")
			}
		})
	}
}

func TestTaskService_MarkCompleted(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockTaskQuerier
		wantErr bool
	}{
		{
			name: "marks task as completed",
			id:   "task-1",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					getTaskByIDFunc: func(ctx context.Context, id string) (db.Task, error) {
						return db.Task{
							ID:        "task-1",
							ProjectID: "project-1",
							Title:     "Test Task",
							Status:    "in_progress",
							Priority:  "high",
							IsNext:    0,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
					updateTaskFunc: func(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
						if arg.Status != "done" {
							t.Errorf("expected status 'done', got %s", arg.Status)
						}
						return db.Task{
							ID:        arg.ID,
							ProjectID: "project-1",
							Title:     arg.Title,
							Status:    arg.Status,
							Priority:  arg.Priority,
							IsNext:    arg.IsNext,
							UpdatedAt: arg.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.MarkCompleted(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.MarkCompleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TaskService.MarkCompleted() returned nil task")
			}
		})
	}
}

func TestTaskService_ListAll(t *testing.T) {
	tests := []struct {
		name    string
		mock    func() *mockTaskQuerier
		want    int
		wantErr bool
	}{
		{
			name: "lists all tasks",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					listAllTasksFunc: func(ctx context.Context) ([]db.Task, error) {
						return []db.Task{
							{ID: "task-1", ProjectID: "project-1", Title: "Task 1", Status: "todo", Priority: "high", IsNext: 0},
							{ID: "task-2", ProjectID: "project-1", Title: "Task 2", Status: "in_progress", Priority: "medium", IsNext: 1},
							{ID: "task-3", ProjectID: "project-2", Title: "Task 3", Status: "todo", Priority: "low", IsNext: 0},
						}, nil
					},
				}
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "returns empty list",
			mock: func() *mockTaskQuerier {
				return &mockTaskQuerier{
					listAllTasksFunc: func(ctx context.Context) ([]db.Task, error) {
						return []db.Task{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTaskService(tt.mock(), nil)
			got, err := svc.ListAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("TaskService.ListAll() returned %d tasks, want %d", len(got), tt.want)
			}
		})
	}
}
