package modal

import (
	"errors"
	"strings"
	"unicode"
)

const MaxTitleLength = 255

var (
	ErrTitleEmpty        = errors.New("title cannot be empty")
	ErrTitleTooLong      = errors.New("title cannot exceed 255 characters")
	ErrTitleInvalidChars = errors.New("title cannot contain newlines or control characters")
)

func ValidateTitle(title string) error {
	trimmed := strings.TrimSpace(title)

	if trimmed == "" {
		return ErrTitleEmpty
	}

	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}

	for _, r := range title {
		if r == '\n' || r == '\r' {
			return ErrTitleInvalidChars
		}
		if unicode.IsControl(r) {
			return ErrTitleInvalidChars
		}
	}

	return nil
}
