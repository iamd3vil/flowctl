package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

func (c *Core) CreateNode(ctx context.Context, node *models.Node, namespaceID string) (models.Node, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Node{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if node.Name == "" {
		return models.Node{}, errors.New("node name is required")
	}
	if node.Hostname == "" {
		return models.Node{}, errors.New("hostname is required")
	}

	credID, err := uuid.Parse(node.Auth.CredentialID)
	if err != nil {
		return models.Node{}, errors.New("invalid credential ID format")
	}

	credential, err := c.store.GetCredentialByUUID(ctx, repo.GetCredentialByUUIDParams{
		Uuid:   credID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.Node{}, errors.New("credential not found")
	}

	created, err := c.store.CreateNode(ctx, repo.CreateNodeParams{
		Name:           node.Name,
		Hostname:       node.Hostname,
		Port:           int32(node.Port),
		Username:       node.Username,
		OsFamily:       node.OSFamily,
		Tags:           node.Tags,
		AuthMethod:     repo.AuthenticationMethod(node.Auth.Method),
		ConnectionType: repo.ConnectionType(node.ConnectionType),
		CredentialID:   sql.NullInt32{Int32: credential.ID, Valid: true},
		Uuid:           namespaceUUID,
	})
	if err != nil {
		return models.Node{}, err
	}

	key := credential.KeyData

	return models.Node{
		ID:             created.Uuid.String(),
		Name:           created.Name,
		Hostname:       created.Hostname,
		Port:           int(created.Port),
		Username:       created.Username,
		OSFamily:       created.OsFamily,
		ConnectionType: string(created.ConnectionType),
		Tags:           created.Tags,
		Auth: models.NodeAuth{
			Method:       node.Auth.Method,
			CredentialID: credential.Uuid.String(),
			Key:          key,
		},
	}, nil
}

func (c *Core) GetNodeByID(ctx context.Context, id string, namespaceID string) (models.Node, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Node{}, err
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Node{}, err
	}

	node, err := c.store.GetNodeByUUID(ctx, repo.GetNodeByUUIDParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.Node{}, err
	}

	credential, err := c.store.GetCredentialByID(ctx, repo.GetCredentialByIDParams{
		ID:   node.CredentialID.Int32,
		Uuid: namespaceUUID,
	})
	if err != nil {
		return models.Node{}, errors.New("credential not found")
	}

	key := credential.KeyData

	return models.Node{
		ID:             node.Uuid.String(),
		Name:           node.Name,
		Hostname:       node.Hostname,
		Port:           int(node.Port),
		Username:       node.Username,
		OSFamily:       node.OsFamily,
		ConnectionType: string(node.ConnectionType),
		Tags:           node.Tags,
		Auth: models.NodeAuth{
			Method:       models.AuthMethod(node.AuthMethod),
			CredentialID: credential.Uuid.String(),
			Key:          key,
		},
	}, nil
}

func (c *Core) SearchNodes(ctx context.Context, filter string, limit, offset int, namespaceID string) ([]models.Node, int64, int64, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, -1, -1, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	nodes, err := c.store.SearchNodes(ctx, repo.SearchNodesParams{
		Uuid:   namespaceUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
		Column4: filter,
	})
	if err != nil {
		return nil, -1, -1, err
	}

	results := make([]models.Node, 0)
	var pageCount, totalCount int64
	for _, n := range nodes {
		res, err := c.GetNodeByID(ctx, n.Uuid.String(), namespaceID)
		if err != nil {
			return nil, -1, -1, err
		}
		results = append(results, res)
		pageCount = n.PageCount
		totalCount = n.TotalCount
	}


	return results, pageCount, totalCount, nil
}

func (c *Core) UpdateNode(ctx context.Context, id string, node *models.Node, namespaceID string) (models.Node, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Node{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if node.Name == "" {
		return models.Node{}, errors.New("node name is required")
	}
	if node.Hostname == "" {
		return models.Node{}, errors.New("hostname is required")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Node{}, err
	}

	credID, _ := uuid.Parse(node.Auth.CredentialID)
	credential, err := c.store.GetCredentialByUUID(ctx, repo.GetCredentialByUUIDParams{
		Uuid:   credID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.Node{}, errors.New("credential not found")
	}

	updated, err := c.store.UpdateNode(ctx, repo.UpdateNodeParams{
		Uuid:           uuidID,
		Name:           node.Name,
		Hostname:       node.Hostname,
		Port:           int32(node.Port),
		Username:       node.Username,
		OsFamily:       node.OSFamily,
		Tags:           node.Tags,
		AuthMethod:     repo.AuthenticationMethod(node.Auth.Method),
		ConnectionType: repo.ConnectionType(node.ConnectionType),
		CredentialID:   sql.NullInt32{Int32: credential.ID, Valid: true},
		Uuid_2:         namespaceUUID,
	})
	if err != nil {
		return models.Node{}, err
	}

	key := credential.KeyData

	return models.Node{
		ID:             updated.Uuid.String(),
		Name:           updated.Name,
		Hostname:       updated.Hostname,
		Port:           int(updated.Port),
		Username:       updated.Username,
		OSFamily:       updated.OsFamily,
		ConnectionType: string(updated.ConnectionType),
		Tags:           updated.Tags,
		Auth: models.NodeAuth{
			Method:       models.AuthMethod(updated.AuthMethod),
			CredentialID: credential.Uuid.String(),
			Key:          key,
		},
	}, nil
}

func (c *Core) DeleteNode(ctx context.Context, id string, namespaceID string) error {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid node UUID: %w", err)
	}
	return c.store.DeleteNode(ctx, repo.DeleteNodeParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
}

func (c *Core) GetNodeStats(ctx context.Context, namespaceID string) (models.NodeStats, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.NodeStats{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	stats, err := c.store.GetNodeStats(ctx, namespaceUUID)
	if err != nil {
		return models.NodeStats{}, fmt.Errorf("error getting node stats: %w", err)
	}

	return models.NodeStats{
		TotalHosts: stats.TotalHosts,
		SSHHosts:   stats.SshHosts,
		QSSHHosts:  stats.QsshHosts,
	}, nil
}
