package core

import (
	"context"
	"errors"

	"github.com/cvhariharan/autopilot/internal/core/models"
	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/google/uuid"
)

func (c *Core) CreateNamespace(ctx context.Context, namespace *models.Namespace) (*models.Namespace, error) {
	if namespace.Name == "" {
		return nil, errors.New("namespace name is required")
	}

	created, err := c.store.CreateNamespace(ctx, namespace.Name)
	if err != nil {
		return nil, err
	}

	return &models.Namespace{
		ID:   created.Uuid.String(),
		Name: created.Name,
	}, nil
}

func (c *Core) GetNamespaceByID(ctx context.Context, id string) (*models.Namespace, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	namespace, err := c.store.GetNamespaceByUUID(ctx, uuidID)
	if err != nil {
		return nil, err
	}

	return &models.Namespace{
		ID:   namespace.Uuid.String(),
		Name: namespace.Name,
	}, nil
}

func (c *Core) ListNamespaces(ctx context.Context, limit, offset int) ([]*models.Namespace, int64, int64, error) {
	namespaces, err := c.store.ListNamespaces(ctx, repo.ListNamespacesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, -1, -1, err
	}

	results := make([]*models.Namespace, len(namespaces))
	for i, n := range namespaces {
		results[i] = &models.Namespace{
			ID:   n.Uuid.String(),
			Name: n.Name,
		}
	}

	if len(namespaces) > 0 {
		return results, namespaces[0].PageCount, namespaces[0].TotalCount, nil
	}
	return results, 0, 0, nil
}

func (c *Core) UpdateNamespace(ctx context.Context, id string, namespace *models.Namespace) (*models.Namespace, error) {
	if namespace.Name == "" {
		return nil, errors.New("namespace name is required")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	updated, err := c.store.UpdateNamespace(ctx, repo.UpdateNamespaceParams{
		Uuid: uuidID,
		Name: namespace.Name,
	})
	if err != nil {
		return nil, err
	}

	return &models.Namespace{
		ID:   updated.Uuid.String(),
		Name: updated.Name,
	}, nil
}

func (c *Core) DeleteNamespace(ctx context.Context, id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return c.store.DeleteNamespace(ctx, uuidID)
}