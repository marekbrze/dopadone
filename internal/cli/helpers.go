package cli

import (
	"log/slog"
)

type Closer interface {
	Close() error
}

func CloseWithLog(closer Closer, name string) {
	if closer == nil {
		return
	}
	if err := closer.Close(); err != nil {
		slog.Warn("failed to close resource", "name", name, "error", err)
	}
}
