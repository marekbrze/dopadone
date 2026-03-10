package cli

import (
	"context"
	"fmt"

	"github.com/marekbrze/dopadone/internal/cli/output"
)

// Deleteable defines the interface for entities that support soft and hard deletion.
// This interface allows the RunDelete helper to work with any entity type
// (areas, subareas, projects, tasks) that implements these operations.
type Deleteable interface {
	// GetByID retrieves an entity to verify its existence before deletion.
	// Returns a nil pointer if the entity is not found.
	GetByID(ctx context.Context, id string) (any, error)

	// SoftDelete marks an entity as deleted without removing it from the database.
	SoftDelete(ctx context.Context, id string) error

	// HardDelete permanently removes an entity from the database.
	HardDelete(ctx context.Context, id string) error
}

// DeleteParams contains the parameters needed for the RunDelete helper.
type DeleteParams struct {
	// ID is the unique identifier of the entity to delete.
	ID string

	// Permanent indicates whether to perform a hard delete (true) or soft delete (false).
	Permanent bool

	// EntityName is the human-readable name of the entity type (e.g., "project", "subarea").
	// Used in error and success messages.
	EntityName string

	// NotFoundErr is the sentinel error returned by the service when an entity is not found.
	// Used to provide a specific "not found" error message.
	NotFoundErr error
}

// RunDelete executes a delete operation (soft or hard) on an entity using the provided service.
// It handles the common pattern of:
//  1. Verify the entity exists
//  2. Call soft delete or hard delete based on the permanent flag
//  3. Return appropriate success/error messages
//
// Example usage:
//
//	params := DeleteParams{
//		ID:          "project-123",
//		Permanent:   false,
//		EntityName:  "project",
//		NotFoundErr: service.ErrProjectNotFound,
//	}
//	if err := RunDelete(ctx, services.Projects, params); err != nil {
//		cli.ExitWithError(err)
//	}
func RunDelete(ctx context.Context, svc Deleteable, params DeleteParams) error {
	// Verify entity exists
	_, err := svc.GetByID(ctx, params.ID)
	if err != nil {
		if params.NotFoundErr != nil && err == params.NotFoundErr {
			return fmt.Errorf("%s not found: %s", params.EntityName, params.ID)
		}
		return WrapError(err, fmt.Sprintf("failed to get %s", params.EntityName))
	}

	// Perform deletion
	if params.Permanent {
		if err := svc.HardDelete(ctx, params.ID); err != nil {
			return WrapError(err, fmt.Sprintf("failed to permanently delete %s", params.EntityName))
		}
		output.PrintSuccess(fmt.Sprintf("%s permanently deleted: %s", params.EntityName, params.ID))
		return nil
	}

	if err := svc.SoftDelete(ctx, params.ID); err != nil {
		if params.NotFoundErr != nil && err == params.NotFoundErr {
			return fmt.Errorf("%s not found: %s", params.EntityName, params.ID)
		}
		return WrapError(err, fmt.Sprintf("failed to delete %s", params.EntityName))
	}

	output.PrintSuccess(fmt.Sprintf("%s deleted: %s", params.EntityName, params.ID))
	return nil
}
