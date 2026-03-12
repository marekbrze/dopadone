package welcome

import (
	"errors"
	"strings"
)

const MaxNameLength = 100

var ErrNameEmpty = errors.New("area name is required")

func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrNameEmpty
	}
	if len(trimmed) > MaxNameLength {
		return errors.New("area name must be 100 characters or less")
	}
	return nil
}
