---
id: TASK-41
title: change width of the columns
status: Done
assignee: []
created_date: '2026-03-05 21:38'
updated_date: '2026-03-06 10:16'
labels: []
dependencies: []
references:
  - internal/tui/views/columns.go
  - internal/tui/app.go
  - .agents/skills/bubbletea/references/golden-rules.md
  - docs/TUI.md
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the three-column TUI layout to use proportional widths instead of equal thirds. Columns should allocate: Subareas 25%, Projects 25%, Tasks 50%.

The layout must be responsive to terminal size changes and should gracefully handle narrow terminals by stacking Subareas and Projects vertically while keeping Tasks as a separate column.

This improves usability by giving more space to tasks, which typically contain longer text and more information than subareas or projects.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Subareas column uses 25% of available width (with minimum width of 20 chars)
- [x] #2 Projects column uses 25% of available width (with minimum width of 20 chars)
- [x] #3 Tasks column uses 50% of available width (with minimum width of 40 chars)
- [x] #4 Column widths are responsive and adjust proportionally when terminal resizes
- [x] #5 On narrow terminals (<120 cols): Subareas and Projects stack vertically, Tasks remain as separate column
- [x] #6 Layout switches smoothly between side-by-side and stacked modes
- [x] #7 Column borders and content render correctly in both layout modes
- [x] #8 No text wrapping occurs in bordered panels (text is properly truncated)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Task-41: Change Width of Columns

## Overview

Task-41 has been split into two focused subtasks for better manageability and incremental delivery:

- **Task-42**: Implement proportional column widths (25/25/50) with responsive behavior
- **Task-43**: Add stacked layout for narrow terminals (<120 cols)

**Dependencies**: Task-43 depends on Task-42 (stacked layout requires proportional widths to be working)

---

## Task-42: Proportional Column Widths with Responsive Behavior

### Phase 1: Layout Calculation Updates (Sequential)

1. **Update Layout function in columns.go**
   - Replace equal division logic with weight-based layout (Golden Rule #4)
   - Define column weights: Subareas=1, Projects=1, Tasks=2 (25/25/50 ratio)
   - Calculate actual widths: `width * weight / totalWeight`
   - Follow Golden Rule #1: Account for borders in height calculations

2. **Add minimum width constraints**
   - Subareas: min 20 chars
   - Projects: min 20 chars  
   - Tasks: min 40 chars
   - Implement enforcement: if calculated width < min, use min

3. **Implement responsive width calculation**
   - Create helper function: `calculateColumnWidths(totalWidth int) (int, int, int)`
   - Use integer math to avoid rounding errors
   - Ensure total width matches available space exactly

### Phase 2: Text Truncation (Sequential)

4. **Add text truncation logic**
   - Follow Golden Rule #2: Never auto-wrap in bordered panels
   - Calculate max text width per column: `columnWidth - 4` (borders + padding)
   - Create/extend `truncateString` helper function
   - Apply truncation to ALL text before rendering (titles, content)

### Phase 3: Testing (Sequential)

5. **Create unit tests for layout calculations**
   - Test file: `internal/tui/views/columns_test.go`
   - Use table-driven tests pattern
   - Test cases:
     * Standard terminal width (120 cols)
     * Wide terminal (160 cols)
     * Narrow terminal approaching minimum (90 cols)
     * Edge case: exactly at minimum widths
     * Edge case: below minimum widths

6. **Test responsive behavior**
   - Simulate WindowSizeMsg at different widths
   - Verify proportional scaling
   - Check minimum width enforcement

7. **Test text truncation**
   - Verify no wrapping occurs
   - Test with long strings
   - Verify ellipsis appears for truncated text

### Phase 4: Manual Verification (Sequential)

8. **Manual testing at multiple terminal sizes**
   - Test at 80 cols (should hit minimums)
   - Test at 120 cols (standard)
   - Test at 160 cols (wide)
   - Verify visual alignment

---

## Task-43: Stacked Layout for Narrow Terminals

### Phase 1: Layout Mode Detection (Sequential)

1. **Add layout mode detection**
   - Create helper function: `shouldUseStackedLayout(width int) bool`
   - Threshold: 120 columns
   - Return true if width < 120

### Phase 2: Stacked Layout Implementation (Sequential)

2. **Implement stacked layout rendering**
   - Create new function: `LayoutStacked(columns []Column, width, height int) string`
   - Split screen into two sections:
     * Left: Tasks column (50% width)
     * Right: Subareas + Projects stacked vertically (50% width)
   - Calculate heights for stacked section:
     * Available height for stack = availableHeight / 2 (per column)
     * Follow Golden Rule #1: Account for borders

3. **Update Layout function to support both modes**
   - Check layout mode at start of Layout()
   - Delegate to LayoutStacked() if narrow terminal
   - Otherwise use proportional side-by-side layout

4. **Handle layout switching**
   - Store current layout mode in app state (optional, for smooth transitions)
   - Ensure View() recalculates on every WindowSizeMsg
   - No animation needed (instant switch is fine)

### Phase 3: Mouse Interaction Updates (Sequential)

5. **Update mouse click detection (if mouse support exists)**
   - Follow Golden Rule #3: Match mouse detection to layout
   - Check layout mode before processing coordinates
   - Use Y coordinates for stacked section
   - Use X coordinates for side-by-side

### Phase 4: Testing (Sequential)

6. **Create unit tests for stacked layout**
   - Test file: `internal/tui/views/columns_test.go` (extend existing)
   - Test cases:
     * Layout mode detection at 119 cols (stacked)
     * Layout mode detection at 120 cols (side-by-side)
     * Layout mode detection at 121 cols (side-by-side)
     * Stacked layout rendering correctness
     * Smooth switching between modes

7. **Test mouse interaction in both modes**
   - Verify clicks detected correctly in stacked mode
   - Verify clicks detected correctly in side-by-side mode

### Phase 5: Manual Verification (Sequential)

8. **Manual testing at threshold width**
   - Test at exactly 120 cols (boundary)
   - Test at 119 cols (stacked mode)
   - Test at 121 cols (side-by-side mode)
   - Resize terminal from 100 → 150 cols and verify smooth switching

---

## Testing Strategy

### Unit Tests (Required)

**File**: `internal/tui/views/columns_test.go` (new file)

Follow golang-testing patterns:
- Use table-driven tests for all test cases
- Use subtests for related scenarios
- Use t.Helper() for helper functions
- Test both success and edge cases

**Test Coverage Targets**:
- `calculateColumnWidths`: 100%
- `shouldUseStackedLayout`: 100%
- `Layout`: 90%+
- `LayoutStacked`: 90%+
- Text truncation: 100%

### Integration Tests

**File**: `internal/tui/integration_test.go` (extend existing)

- Test full TUI rendering at different terminal sizes
- Test WindowSizeMsg handling
- Verify no visual artifacts

### Manual Testing Checklist

For both tasks:
- [ ] Test at 80 cols (narrow)
- [ ] Test at 120 cols (standard)
- [ ] Test at 160 cols (wide)
- [ ] Verify column borders align correctly
- [ ] Verify no text wrapping
- [ ] Verify smooth resize transitions

For task-43:
- [ ] Test at 119 cols (stacked mode)
- [ ] Test at 120 cols (side-by-side boundary)
- [ ] Test at 121 cols (side-by-side mode)
- [ ] Verify stacked columns have equal heights
- [ ] Verify Tasks column spans full height
- [ ] Test mouse clicks in both modes (if applicable)

---

## Documentation Updates

### Code Documentation

1. **Update columns.go comments**
   - Document proportional layout algorithm
   - Document minimum width constraints
   - Document stacked layout behavior
   - Add examples in function comments

2. **Update TUI.md**
   - Update "Three-Column Browser" section:
     * Document 25/25/50 proportional widths
     * Document minimum width constraints
     * Document stacked layout for narrow terminals
     * Add visual diagram showing both layouts
   - Update "Responsive Design" section:
     * Document threshold width (120 cols)
     * Explain layout switching behavior

3. **Update START_HERE.md** (if needed)
   - Update any references to column layout

### Inline Documentation

Add godoc comments for new functions:
```go
// calculateColumnWidths calculates proportional column widths using weight-based layout.
// Weights: Subareas=1, Projects=1, Tasks=2 (25/25/50 ratio).
// Enforces minimum widths: Subareas=20, Projects=20, Tasks=40.
func calculateColumnWidths(totalWidth int) (subareasWidth, projectsWidth, tasksWidth int)

// shouldUseStackedLayout returns true if terminal width < 120 columns.
// In stacked mode, Subareas and Projects stack vertically on the right,
// while Tasks occupy the left column.
func shouldUseStackedLayout(width int) bool
```

---

## Execution Plan: Sequential vs Parallel

### Sequential Tasks (MUST be done in order)

1. **Task-42 → Task-43**: Task-43 depends on task-42's proportional layout
2. **Within Task-42**: Phases 1-4 must be sequential
3. **Within Task-43**: Phases 1-5 must be sequential

### Parallel Opportunities

**Within Task-42**:
- Phases 3 (Testing) and 4 (Manual Verification) can overlap
- While tests are being written, manual verification can start

**Within Task-43**:
- Phases 4 (Testing) and 5 (Manual Verification) can overlap

**Documentation**:
- Code documentation can be written alongside implementation
- TUI.md updates can be done in parallel with testing

---

## Code Quality Checklist

Follow golang-patterns and golang-pro guidelines:

- [ ] Use gofmt and goimports
- [ ] Pass go vet
- [ ] Pass golangci-lint
- [ ] No naked returns
- [ ] Explicit error handling (if applicable)
- [ ] Document all exported functions
- [ ] Use meaningful variable names
- [ ] Keep functions focused (single responsibility)
- [ ] Use table-driven tests
- [ ] Run tests with race detector: `go test -race ./internal/tui/...`

---

## Key Golden Rules to Follow

From bubbletea skill, apply these critical rules:

1. **Rule #1 - Account for Borders**: Subtract 2 from height BEFORE rendering
   ```go
   availableHeight := height - tabsHeight - footerHeight - 2  // -2 for borders
   ```

2. **Rule #2 - Never Auto-Wrap**: Truncate ALL text explicitly
   ```go
   maxTextWidth := columnWidth - 4  // -2 borders, -2 padding
   title = truncateString(title, maxTextWidth)
   ```

3. **Rule #3 - Match Mouse to Layout**: Different logic for stacked vs side-by-side
   ```go
   if shouldUseStackedLayout(width) {
       // Use Y coordinates for stacked section
   } else {
       // Use X coordinates for side-by-side
   }
   ```

4. **Rule #4 - Use Weights, Not Pixels**: Weight-based layout for perfect scaling
   ```go
   subareasWeight, projectsWeight, tasksWeight := 1, 1, 2
   totalWeight := subareasWeight + projectsWeight + tasksWeight
   subareasWidth := (availableWidth * subareasWeight) / totalWeight
   ```

---

## Risks and Mitigations

### Risk 1: Text Wrapping in Narrow Terminals
**Mitigation**: Implement strict text truncation with ellipsis. Test at minimum widths.

### Risk 2: Layout Switching Jank
**Mitigation**: Ensure View() is pure function that recalculates on every render. No caching needed.

### Risk 3: Mouse Clicks Not Working After Layout Change
**Mitigation**: Follow Golden Rule #3, add layout mode checks in mouse handlers.

### Risk 4: Border Misalignment
**Mitigation**: Follow Golden Rule #1 strictly. Use same border style for all panels.

---

## Success Criteria

Task-41 is complete when:

1. ✅ All acceptance criteria in task-42 are met
2. ✅ All acceptance criteria in task-43 are met
3. ✅ Unit tests pass with >90% coverage on new code
4. ✅ Manual testing completed at multiple terminal sizes
5. ✅ Code passes lint and vet
6. ✅ Documentation updated (TUI.md, code comments)
7. ✅ No regressions in existing TUI tests
8. ✅ Smooth user experience across all terminal sizes
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Split into two focused subtasks for better manageability:
- Task-42: Proportional column widths with responsive behavior
- Task-43: Stacked layout for narrow terminals

See task-42 and task-43 for detailed implementation.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Completed via subtasks task-42 (proportional column widths) and task-43 (stacked layout for narrow terminals). See those tasks for implementation details.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All acceptance criteria are met and verified
- [x] #2 Code passes go vet and golangci-lint
- [x] #3 Unit tests added for new layout calculation logic
- [x] #4 Manual testing performed at multiple terminal sizes
- [x] #5 No regressions in existing TUI tests
- [x] #6 Code follows Go patterns and Bubbletea golden rules
<!-- DOD:END -->
