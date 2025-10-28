package core

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

func (c *Core) CreateCredential(ctx context.Context, cred models.Credential, namespaceID string) (models.Credential, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Credential{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if cred.Name == "" {
		return models.Credential{}, errors.New("credential name is required")
	}

	if cred.KeyData == "" {
		return models.Credential{}, errors.New("key data is required")
	}

	if cred.KeyType == "" {
		return models.Credential{}, errors.New("key type is required")
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(cred.KeyData))
	if err != nil {
		return models.Credential{}, err
	}
	encryptedKeyData := hex.EncodeToString(enc)

	created, err := c.store.CreateCredential(ctx, repo.CreateCredentialParams{
		Name:    cred.Name,
		KeyType: cred.KeyType,
		KeyData: encryptedKeyData,
		Uuid:    namespaceUUID,
	})
	if err != nil {
		return models.Credential{}, err
	}

	var lastAccessed string
	if created.LastAccessed.Valid {
		lastAccessed = created.LastAccessed.Time.Format(TimeFormat)
	}

	return models.Credential{
		ID:           created.Uuid.String(),
		Name:         created.Name,
		KeyType:      created.KeyType,
		KeyData:      created.KeyData,
		LastAccessed: lastAccessed,
	}, nil
}

func (c *Core) GetCredentialByID(ctx context.Context, id string, namespaceID string) (models.Credential, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Credential{}, err
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Credential{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	cred, err := c.store.GetCredentialByUUID(ctx, repo.GetCredentialByUUIDParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.Credential{}, err
	}

	var lastAccessed string
	if cred.LastAccessed.Valid {
		lastAccessed = cred.LastAccessed.Time.Format(TimeFormat)
	}

	return models.Credential{
		ID:           cred.Uuid.String(),
		Name:         cred.Name,
		KeyType:      cred.KeyType,
		KeyData:      cred.KeyData,
		LastAccessed: lastAccessed,
	}, nil
}

func (c *Core) SearchCredentials(ctx context.Context, filter string, limit, offset int, namespaceID string) ([]models.Credential, int64, int64, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, -1, -1, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	creds, err := c.store.SearchCredentials(ctx, repo.SearchCredentialsParams{
		Uuid:    namespaceUUID,
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column4: filter,
	})
	if err != nil {
		return nil, -1, -1, err
	}

	results := make([]models.Credential, 0, len(creds))
	var pageCount, totalCount int64
	for _, cred := range creds {
		var lastAccessed string
		if cred.LastAccessed.Valid {
			lastAccessed = cred.LastAccessed.Time.Format(TimeFormat)
		}
		results = append(results, models.Credential{
			ID:           cred.Uuid.String(),
			Name:         cred.Name,
			KeyType:      cred.KeyType,
			KeyData:      cred.KeyData,
			LastAccessed: lastAccessed,
		})
		pageCount = cred.PageCount
		totalCount = cred.TotalCount
	}

	return results, pageCount, totalCount, nil
}

func (c *Core) UpdateCredential(ctx context.Context, id string, cred *models.Credential, namespaceID string) (models.Credential, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.Credential{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if cred.Name == "" {
		return models.Credential{}, errors.New("credential name is required")
	}
	if cred.KeyData == "" {
		return models.Credential{}, errors.New("key data is required")
	}

	if cred.KeyType == "" {
		return models.Credential{}, errors.New("key type is required")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Credential{}, err
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(cred.KeyData))
	if err != nil {
		return models.Credential{}, err
	}
	encryptedKeyData := hex.EncodeToString(enc)

	updated, err := c.store.UpdateCredential(ctx, repo.UpdateCredentialParams{
		Uuid:    uuidID,
		Name:    cred.Name,
		KeyType: cred.KeyType,
		KeyData: encryptedKeyData,
		Uuid_2:  namespaceUUID,
	})
	if err != nil {
		return models.Credential{}, err
	}

	var lastAccessed string
	if updated.LastAccessed.Valid {
		lastAccessed = updated.LastAccessed.Time.Format(TimeFormat)
	}

	return models.Credential{
		ID:           updated.Uuid.String(),
		Name:         updated.Name,
		KeyType:      updated.KeyType,
		KeyData:      updated.KeyData,
		LastAccessed: lastAccessed,
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
