# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

#### Tree Visual Design (Task-45)

**Modernized project tree rendering with arrow-based indicators**

The project tree component now uses a clean, minimalist design with arrow indicators instead of traditional box-drawing characters.

**Changes**:
- Replaced box-drawing characters (├─└│) with simple 2-space indentation
- Replaced expand/collapse indicators `[-]`/`[+]` with arrows `▾`/`▸`
- Removed vertical connector lines for cleaner visual appearance
- Improved readability with consistent indentation at all depth levels

**Files Modified**:
- `internal/tui/tree/constants.go`: Updated tree character constants
- `internal/tui/tree/renderer.go`: Simplified rendering logic
- `internal/tui/tree/renderer_test.go`: Updated test expectations
- `docs/TUI.md`: Added documentation for new tree styling

**Visual Comparison**:

Before (box-drawing):
```
├─ Project A
│  ├─ Subproject A1
│  └─ Subproject A2
└─ Project B
```

After (arrow indicators):
```
▾ Project A
  Subproject A1
  ▸ Subproject A2
Project B
```

**Benefits**:
- Reduced visual clutter with no vertical connector lines
- Clearer expand/collapse state with intuitive arrow indicators
- Modern, minimalist appearance
- Better readability on high-DPI displays
- Customizable through `TreeStyle` struct

**Backward Compatibility**:
- All existing tree navigation and functionality preserved
- Only visual rendering changed, no API changes
- Custom tree styles can still be applied via `TreeStyle` struct

**Testing**:
- All 45 tree renderer tests updated and passing
- Visual verification in TUI confirms modern appearance
- Expand/collapse functionality verified working
- Navigation preserved across all tree operations

**Documentation**:
- Added tree styling section to `docs/TUI.md`
- Updated tree rendering examples
- Documented customization options via `TreeStyle`
