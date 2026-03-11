package migrate

import (
	"database/sql"
	"fmt"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

type MigrationSyncer struct {
	driver driver.DatabaseDriver
}

func NewMigrationSyncer(d driver.DatabaseDriver) *MigrationSyncer {
	return &MigrationSyncer{driver: d}
}

func (s *MigrationSyncer) RunAndSync(db *sql.DB, command string) error {
	if err := Run(db, command); err != nil {
		return err
	}

	if s.driver.Type() == driver.DriverTursoReplica {
		if syncer, ok := s.driver.(interface{ Sync() error }); ok {
			if err := syncer.Sync(); err != nil {
				return fmt.Errorf("failed to sync schema to remote: %w", err)
			}
		}
	}

	return nil
}
