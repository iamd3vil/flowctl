package core

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cvhariharan/autopilot/internal/core/models"
	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/google/uuid"
)

func (c *Core) CreateCredential(ctx context.Context, cred *models.Credential, namespaceID string) (*models.Credential, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if cred.Name == "" {
		return nil, errors.New("credential name is required")
	}

	// At least one of password or private key should be present
	if cred.Password == "" && cred.PrivateKey == "" {
		return nil, errors.New("either password or private key is required")
	}

	if cred.Password != "" && cred.PrivateKey != "" {
		return nil, errors.New("only one of password or private key can be set at a time")
	}

	var encryptedPassword, encryptedPrivateKey string
	if cred.Password != "" {
		enc, err := c.keeper.Encrypt(ctx, []byte(cred.Password))
		if err != nil {
			return nil, err
		}
		encryptedPassword = hex.EncodeToString(enc)
	} else if cred.PrivateKey != "" {
		enc, err := c.keeper.Encrypt(ctx, []byte(cred.PrivateKey))
		if err != nil {
			return nil, err
		}
		encryptedPrivateKey = hex.EncodeToString(enc)
	}

	created, err := c.store.CreateCredential(ctx, repo.CreateCredentialParams{
		Name: cred.Name,
		PrivateKey: sql.NullString{
			String: encryptedPrivateKey,
			Valid:  encryptedPrivateKey != "",
		},
		Password: sql.NullString{
			String: encryptedPassword,
			Valid:  encryptedPassword != "",
		},
		Uuid: namespaceUUID,
	})
	if err != nil {
		return nil, err
	}

	return &models.Credential{
		ID:         created.Uuid.String(),
		Name:       created.Name,
		PrivateKey: created.PrivateKey.String,
		Password:   created.Password.String,
	}, nil
}

func (c *Core) GetCredentialByID(ctx context.Context, id string, namespaceID string) (*models.Credential, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	cred, err := c.store.GetCredentialByUUID(ctx, repo.GetCredentialByUUIDParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return nil, err
	}

	return &models.Credential{
		ID:         cred.Uuid.String(),
		Name:       cred.Name,
		PrivateKey: cred.PrivateKey.String,
		Password:   cred.Password.String,
	}, nil
}

func (c *Core) ListCredentials(ctx context.Context, limit, offset int, namespaceID string) ([]*models.Credential, int64, int64, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, -1, -1, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	creds, err := c.store.ListCredentials(ctx, repo.ListCredentialsParams{
		Uuid:   namespaceUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, -1, -1, err
	}

	results := make([]*models.Credential, 0)
	for _, cred := range creds {
		results = append(results, &models.Credential{
			ID:         cred.Uuid.String(),
			Name:       cred.Name,
			PrivateKey: cred.PrivateKey.String,
			Password:   cred.Password.String,
		})
	}

	if len(creds) > 0 {
		return results, creds[0].PageCount, creds[0].TotalCount, nil
	}

	return results, 0, 0, nil
}

func (c *Core) UpdateCredential(ctx context.Context, id string, cred *models.Credential, namespaceID string) (*models.Credential, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if cred.Name == "" {
		return nil, errors.New("credential name is required")
	}
	if cred.Password == "" && cred.PrivateKey == "" {
		return nil, errors.New("either password or private key is required")
	}

	if cred.Password != "" && cred.PrivateKey != "" {
		return nil, errors.New("only one of password or private key can be set at a time")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var encryptedPassword, encryptedPrivateKey string
	if cred.Password != "" {
		enc, err := c.keeper.Encrypt(ctx, []byte(cred.Password))
		if err != nil {
			return nil, err
		}
		encryptedPassword = hex.EncodeToString(enc)
	} else if cred.PrivateKey != "" {
		enc, err := c.keeper.Encrypt(ctx, []byte(cred.PrivateKey))
		if err != nil {
			return nil, err
		}
		encryptedPrivateKey = hex.EncodeToString(enc)
	}

	updated, err := c.store.UpdateCredential(ctx, repo.UpdateCredentialParams{
		Uuid: uuidID,
		Name: cred.Name,
		PrivateKey: sql.NullString{
			String: encryptedPrivateKey,
			Valid:  encryptedPrivateKey != "",
		},
		Password: sql.NullString{
			String: encryptedPassword,
			Valid:  encryptedPassword != "",
		},
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return nil, err
	}

	return &models.Credential{
		ID:         updated.Uuid.String(),
		Name:       updated.Name,
		PrivateKey: updated.PrivateKey.String,
		Password:   updated.Password.String,
	}, nil
}

func (c *Core) DeleteCredential(ctx context.Context, id string, namespaceID string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid credential UUID: %w", err)
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	return c.store.DeleteCredential(ctx, repo.DeleteCredentialParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
}
