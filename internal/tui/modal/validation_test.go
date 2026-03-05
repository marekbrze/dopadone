package modal

import (
	"strings"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		shouldError bool
		expectedErr error
	}{
		{
			name:        "valid title",
			title:       "Valid Title",
			shouldError: false,
		},
		{
			name:        "valid title with spaces",
			title:       "  Valid Title  ",
			shouldError: false,
		},
		{
			name:        "empty string",
			title:       "",
			shouldError: true,
			expectedErr: ErrTitleEmpty,
		},
		{
			name:        "whitespace only",
			title:       "   ",
			shouldError: true,
			expectedErr: ErrTitleEmpty,
		},
		{
			name:        "with newline",
			title:       "Title\n",
			shouldError: true,
			expectedErr: ErrTitleInvalidChars,
		},
		{
			name:        "with carriage return",
			title:       "Title\r",
			shouldError: true,
			expectedErr: ErrTitleInvalidChars,
		},
		{
			name:        "with tab",
			title:       "Title\t",
			shouldError: true,
			expectedErr: ErrTitleInvalidChars,
		},
		{
			name:        "with null byte",
			title:       "Title\x00",
			shouldError: true,
			expectedErr: ErrTitleInvalidChars,
		},
		{
			name:        "with escape char",
			title:       "Title\x1b",
			shouldError: true,
			expectedErr: ErrTitleInvalidChars,
		},
		{
			name:        "valid unicode",
			title:       "Title 你好 مرحبا",
			shouldError: false,
		},
		{
			name:        "valid emoji",
			title:       "Title 🚀 ✨",
			shouldError: false,
		},
		{
			name:        "max length",
			title:       strings.Repeat("a", MaxTitleLength),
			shouldError: false,
		},
		{
			name:        "over max length",
			title:       strings.Repeat("a", MaxTitleLength+1),
			shouldError: true,
			expectedErr: ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
			}
		})
	}
}

func TestValidateTitleEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		title string
		valid bool
	}{
		{
			name:  "single character",
			title: "a",
			valid: true,
		},
		{
			name:  "numbers only",
			title: "12345",
			valid: true,
		},
		{
			name:  "special characters",
			title: "Test-Feature_2.0!",
			valid: true,
		},
		{
			name:  "multiple spaces in middle",
			title: "Test   Title",
			valid: true,
		},
		{
			name:  "CRLF line ending",
			title: "Title\r\n",
			valid: false,
		},
		{
			name:  "vertical tab",
			title: "Title\v",
			valid: false,
		},
		{
			name:  "form feed",
			title: "Title\f",
			valid: false,
		},
		{
			name:  "backspace",
			title: "Title\x08",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			isValid := err == nil

			if isValid != tt.valid {
				t.Errorf("expected valid=%v, got valid=%v (err=%v)", tt.valid, isValid, err)
			}
		})
	}
}

func TestValidateTitleMessages(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		expectedError string
	}{
		{
			name:          "empty title error message",
			title:         "",
			expectedError: "empty",
		},
		{
			name:          "too long error message",
			title:         strings.Repeat("a", MaxTitleLength+1),
			expectedError: "255",
		},
		{
			name:          "invalid chars error message",
			title:         "Test\nTitle",
			expectedError: "newlines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if err == nil {
				t.Fatal("expected error but got none")
			}

			errMsg := err.Error()
			if !strings.Contains(strings.ToLower(errMsg), strings.ToLower(tt.expectedError)) {
				t.Errorf("expected error message to contain %q, got %q", tt.expectedError, errMsg)
			}
		})
	}
}
