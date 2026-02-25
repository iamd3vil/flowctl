package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

// ResolvePrefixID looks up a flow prefix by name in the given namespace (by UUID).
// Returns sql.NullInt32{Valid: false} when prefix is empty.
// Returns an error if the prefix name is non-empty but not found.
func (c *Core) ResolvePrefixID(ctx context.Context, prefix, namespaceID string) (sql.NullInt32, error) {
	if prefix == "" {
		return sql.NullInt32{}, nil
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return sql.NullInt32{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	fp, err := c.store.GetFlowPrefixByName(ctx, repo.GetFlowPrefixByNameParams{
		Name: prefix,
		Uuid: namespaceUUID,
	})
	if err != nil {
		return sql.NullInt32{}, fmt.Errorf("flow prefix %q not found in namespace %s: %w", prefix, namespaceID, err)
	}

	return sql.NullInt32{Int32: fp.ID, Valid: true}, nil
}

// CreateFlowPrefix creates a new flow prefix in the given namespace.
func (c *Core) CreateFlowPrefix(ctx context.Context, namespaceID, name, description string) (models.FlowPrefix, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	fp, err := c.store.CreateFlowPrefix(ctx, repo.CreateFlowPrefixParams{
		Uuid:        namespaceUUID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("failed to create flow prefix: %w", err)
	}

	return repoFlowPrefixToModel(fp), nil
}

// UpdateFlowPrefix updates an existing flow prefix.
func (c *Core) UpdateFlowPrefix(ctx context.Context, prefixUUID, namespaceID, name, description string) (models.FlowPrefix, error) {
	pUUID, err := uuid.Parse(prefixUUID)
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("invalid prefix UUID: %w", err)
	}
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	fp, err := c.store.UpdateFlowPrefix(ctx, repo.UpdateFlowPrefixParams{
		Uuid:        pUUID,
		Uuid_2:      namespaceUUID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("failed to update flow prefix: %w", err)
	}

	return repoFlowPrefixToModel(fp), nil
}

// DeleteFlowPrefix deletes a flow prefix by UUID and all flows in the group.
func (c *Core) DeleteFlowPrefix(ctx context.Context, prefixUUID, namespaceID string) error {
	pUUID, err := uuid.Parse(prefixUUID)
	if err != nil {
		return fmt.Errorf("invalid prefix UUID: %w", err)
	}
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	// Get all flows in this group before deleting the prefix
	flowRows, err := c.store.GetFlowsByPrefixUUID(ctx, repo.GetFlowsByPrefixUUIDParams{
		Uuid:   pUUID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return fmt.Errorf("failed to get flows for prefix: %w", err)
	}

	// Delete the prefix (sets flows' prefix_id to NULL via ON DELETE SET NULL)
	if err := c.store.DeleteFlowPrefix(ctx, repo.DeleteFlowPrefixParams{
		Uuid:   pUUID,
		Uuid_2: namespaceUUID,
	}); err != nil {
		return err
	}

	// Delete all flows in the group in the background
	go func() {
		for _, f := range flowRows {
			if err := c.DeleteFlow(context.Background(), f.Slug, namespaceID); err != nil {
				log.Printf("error deleting flow %s during group deletion: %v", f.Slug, err)
			}
		}
	}()

	return nil
}

// GetFlowPrefix returns a flow prefix by UUID.
func (c *Core) GetFlowPrefix(ctx context.Context, prefixUUID, namespaceID string) (models.FlowPrefix, error) {
	pUUID, err := uuid.Parse(prefixUUID)
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("invalid prefix UUID: %w", err)
	}
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	fp, err := c.store.GetFlowPrefixByUUID(ctx, repo.GetFlowPrefixByUUIDParams{
		Uuid:   pUUID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.FlowPrefix{}, fmt.Errorf("flow prefix not found: %w", err)
	}

	return repoFlowPrefixToModel(fp), nil
}

// ListFlowPrefixes returns all flow prefixes in a namespace.
func (c *Core) ListFlowPrefixes(ctx context.Context, namespaceID string) ([]models.FlowPrefix, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	rows, err := c.store.ListFlowPrefixes(ctx, namespaceUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list flow prefixes: %w", err)
	}

	result := make([]models.FlowPrefix, 0, len(rows))
	for _, row := range rows {
		result = append(result, repoFlowPrefixToModel(row))
	}
	return result, nil
}

// getOrCreateFlowPrefix looks up a prefix by name, auto-creating it if missing.
// Used during filesystem-based flow import where prefixes are inferred from directory names.
func (c *Core) getOrCreateFlowPrefix(ctx context.Context, prefix string, namespaceUUID uuid.UUID) (sql.NullInt32, error) {
	if prefix == "" {
		return sql.NullInt32{}, nil
	}

	fp, err := c.store.GetFlowPrefixByName(ctx, repo.GetFlowPrefixByNameParams{
		Name: prefix,
		Uuid: namespaceUUID,
	})
	if err != nil {
		fp, err = c.store.CreateFlowPrefix(ctx, repo.CreateFlowPrefixParams{
			Uuid:        namespaceUUID,
			Name:        prefix,
			Description: "",
		})
		if err != nil {
			return sql.NullInt32{}, fmt.Errorf("failed to create flow prefix %s: %w", prefix, err)
		}
	}

	return sql.NullInt32{Int32: fp.ID, Valid: true}, nil
}

func repoFlowPrefixToModel(fp repo.FlowPrefix) models.FlowPrefix {
	return models.FlowPrefix{
		ID:          fp.Uuid.String(),
		Name:        fp.Name,
		Description: fp.Description,
	}
}
