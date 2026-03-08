---
id: TASK-58
title: 'TUI Rendering: Grouped Display (51D)'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 21:30'
updated_date: '2026-03-07 15:25'
labels:
  - tui
  - rendering
dependencies:
  - TASK-54
  - TASK-57
references:
  - task-51
  - internal/tui/renderer.go
  - .agents/skills/bubbletea/references/golden-rules.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement visual rendering of tasks grouped by subproject with indented headers, styling, and proper text truncation. Depends on tasks 51B (TASK-54) and 51C (TASK-57). Part of task-51 nested task grouping feature.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 #1 Render tasks grouped by subproject with indented headers showing subproject name
- [x] #2 #2 Show direct project tasks at the top without header (ungrouped)
- [x] #3 #3 Use indentation (2 spaces) for tasks under each subproject group header
- [x] #4 #4 Use subtle styling for group headers (dimmed color, no reverse highlight)
- [x] #5 #5 Add text truncation to prevent wrapping in narrow columns
- [x] #6 #6 Write rendering tests (empty, direct only, groups, mixed, styling, truncation)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Detailed Implementation Plan: Task-58 - TUI Rendering - Grouped Display

## Overview
This task implements visual rendering of tasks grouped by subproject, integrating the GroupedTasks domain model into the TUI rendering layer. It completes the nested task grouping feature (task-51).

## Task Assessment
- **Scope**: Focused rendering task with 6 ACs
- **Complexity**: Medium (rendering logic, visual design, testing)
- **Dependencies**: Task-54 ✅ DONE, Task-57 ✅ DONE
- **Decision**: ✅ NO SPLIT NEEDED - Appropriately sized (2-3 hours)

## Sequential Phases (MUST be done in order)

### Phase 1: Helper Methods (30 min)
Create rendering utilities: renderGroupHeader, renderTaskLine, truncateString

### Phase 2: Core Rendering Logic (45 min)
Update RenderTasks to use grouped structure with direct tasks + grouped tasks

### Phase 3: Integration (10 min)
Update View() and ensure width calculation

### Phase 4: Comprehensive Testing (60 min)
Write test suite covering: empty, direct only, grouped, mixed, truncation, styling

### Phase 5: Code Quality (10 min)
Format, lint, fix issues

### Phase 6: Documentation (15 min)
Update godoc comments, TUI.md

## Estimated Total: 2h 50min + 30min buffer = ~3.5 hours

## Dependencies
- Upstream: Task-54 (✅), Task-57 (✅)
- Downstream: Task-56 (TUI Interaction)

## Test Strategy
- Table-driven tests for all scenarios
- Target coverage: 80%+
- Test suites: empty state, direct only, grouped, mixed, truncation, styling

## Definition of Done
- All 6 ACs checked
- Tests passing with race detection
- Coverage ≥ 80%
- Code formatted and linted
- Documentation updated
- Task status: Done
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
IMPLEMENTATION DETAILS:

## File: internal/tui/renderer.go

### Method: RenderTasks (NEW)

```go
func (m *Model) RenderTasks() string {
    if m.isLoadingTasks {
        return m.spinner.View() + " " + LoadingMessageTasks
    }
    
    if m.groupedTasks.TotalCount == 0 {
        return EmptyStateNoTasks
    }
    
    var lines []string
    taskIndex := 0
    
    // Render direct tasks (no header)
    for _, task := range m.groupedTasks.DirectTasks {
        lines = append(lines, m.renderTaskLine(task, taskIndex, 0))
        taskIndex++
    }
    
    // Add separator if we have both direct and grouped tasks
    if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
        lines = append(lines, "")
    }
    
    // Render grouped tasks
    for _, group := range m.groupedTasks.Groups {
        // Add group header
        lines = append(lines, m.renderGroupHeader(group))
        
        // Add tasks if group is expanded
        if m.expandedTaskGroups[group.ProjectID] {
            for _, task := range group.Tasks {
                lines = append(lines, m.renderTaskLine(task, taskIndex, 1))
                taskIndex++
            }
        }
    }
    
    return strings.Join(lines, "
")
}
```

### Method: renderGroupHeader (NEW)

```go
func (m *Model) renderGroupHeader(group domain.TaskGroup) string {
    // Expand/collapse icon
    icon := "▾"
    if !m.expandedTaskGroups[group.ProjectID] {
        icon = "▸"
    }
    
    // Task count
    taskCount := len(group.Tasks)
    taskLabel := "task"
    if taskCount != 1 {
        taskLabel = "tasks"
    }
    
    // Build header text
    header := fmt.Sprintf("%s %s (%d %s)", icon, group.ProjectName, taskCount, taskLabel)
    
    // Apply styling (dimmed, no reverse)
    style := lipgloss.NewStyle().
        Foreground(m.theme.Dimmed).
        PaddingLeft(2)  // Indent from left border
        
    return style.Render(header)
}
```

### Method: renderTaskLine (UPDATED)

```go
func (m *Model) renderTaskLine(task domain.Task, index int, indentLevel int) string {
    // Calculate indentation
    indent := strings.Repeat("  ", indentLevel+1)  // +1 for base indent
    
    // Determine prefix and style
    var prefix string
    var text string
    var style lipgloss.Style
    
    if task.IsCompleted() {
        prefix = "✓ "
        text = task.Title
        style = lipgloss.NewStyle().
            Strikethrough(true).
            Foreground(m.theme.Muted)
    } else {
        prefix = "  "
        text = task.Title
        style = lipgloss.NewStyle().
            Foreground(m.theme.Foreground)
    }
    
    // Highlight selected task
    if index == m.selectedTaskIndex {
        style = style.Bold(true).Reverse(true)
    }
    
    // Truncate text to prevent wrapping
    maxTextWidth := m.taskColumnWidth - len(indent) - len(prefix) - 4  // -4 for borders and padding
    if maxTextWidth > 0 && len(text) > maxTextWidth {
        text = text[:maxTextWidth-1] + "…"
    }
    
    return indent + style.Render(prefix+text)
}
```

## File: internal/tui/tui.go

### Update View method

```go
func (m Model) View() string {
    // ... existing layout code ...
    
    // Render tasks column (UPDATED to use RenderTasks)
    tasksContent := m.RenderTasks()  // NEW METHOD
    
    tasksPanel := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(0, 1).
        Height(tasksColumnHeight).
        Render(tasksContent)
    
    // ... rest of layout ...
}
```

## Visual Design Guidelines

Following bubbletea skill Golden Rules:

1. **Always Account for Borders**
   - maxTextWidth = columnWidth - 4 (2 borders + 2 padding)

2. **Never Auto-Wrap**
   - Always truncate with "…" character
   - Use explicit width calculations

3. **Consistent Styling**
   - Group headers: dimmed, no reverse highlight
   - Tasks: normal foreground color
   - Selected task: bold + reverse
   - Completed tasks: strikethrough + muted

4. **Indentation**
   - Direct tasks: 2 spaces (1 level)
   - Nested tasks: 4 spaces (2 levels)
   - Group headers: 2 spaces (aligned with direct tasks)

## Testing Strategy

### File: internal/tui/renderer_test.go

```go
func TestRenderTasks(t *testing.T) {
    tests := []struct {
        name     string
        model    Model
        wantContains []string
    }{
        {
            name: "empty state",
            model: Model{
                groupedTasks: domain.GroupedTasks{},
            },
            wantContains: []string{"No tasks"},
        },
        {
            name: "direct tasks only",
            model: Model{
                groupedTasks: domain.GroupedTasks{
                    DirectTasks: []domain.Task{
                        {Title: "Task 1"},
                        {Title: "Task 2"},
                    },
                    TotalCount: 2,
                },
                theme: DefaultTheme(),
            },
            wantContains: []string{"Task 1", "Task 2"},
        },
        {
            name: "grouped with expanded",
            model: Model{
                groupedTasks: domain.GroupedTasks{
                    Groups: []domain.TaskGroup{
                        {
                            ProjectID:   "p2",
                            ProjectName: "Backend",
                            Tasks:       []domain.Task{{Title: "API"}},
                        },
                    },
                    TotalCount: 1,
                },
                expandedTaskGroups: map[string]bool{"p2": true},
                theme:              DefaultTheme(),
            },
            wantContains: []string{"Backend", "▾", "API"},
        },
        {
            name: "grouped with collapsed",
            model: Model{
                groupedTasks: domain.GroupedTasks{
                    Groups: []domain.TaskGroup{
                        {
                            ProjectID:   "p2",
                            ProjectName: "Backend",
                            Tasks:       []domain.Task{{Title: "API"}},
                        },
                    },
                    TotalCount: 1,
                },
                expandedTaskGroups: map[string]bool{"p2": false},
                theme:              DefaultTheme(),
            },
            wantContains: []string{"Backend", "▸"},
            wantNotContains: []string{"API"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.model.RenderTasks()
            
            for _, want := range tt.wantContains {
                assert.Contains(t, got, want)
            }
        })
    }
}
```

## Performance Considerations

- O(n) rendering where n = total visible tasks
- String builder for efficiency
- No nested loops
- Truncation before styling (avoid re-rendering)

AC1,

Completed AC rendering: GroupedTasks:

 implementing visual rendering of tasks grouped by subproject with indented headers, styling, and truncation.

Implementation notes updated
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TUI rendering: tasks grouped by subproject with indented headers and dimmed styling, proper text truncation.

- internal/tui/renderer.go: Updated RenderTasks() to use GroupedTasks domain model from task service
- Updated taskColumnWidth() helper with proper column width calculation
- All acceptance criteria verified
<!-- SECTION:FINAL_SUMMARY:END -->
