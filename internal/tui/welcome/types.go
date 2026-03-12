package welcome

import "github.com/marekbrze/dopadone/internal/domain"

type SubmitMsg struct {
	Name  string
	Color domain.Color
}

type ExitMsg struct{}
