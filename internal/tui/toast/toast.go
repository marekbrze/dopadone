package toast

import (
	"time"
)

type Toast struct {
	Type      string
	Message   string
	CreatedAt time.Time
}

func NewError(message string) Toast {
	return Toast{
		Type:      TypeError,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

func NewSuccess(message string) Toast {
	return Toast{
		Type:      TypeSuccess,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

func NewInfo(message string) Toast {
	return Toast{
		Type:      TypeInfo,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

func (t Toast) IsExpired() bool {
	return time.Since(t.CreatedAt) > time.Duration(ToastDuration)*time.Millisecond
}
