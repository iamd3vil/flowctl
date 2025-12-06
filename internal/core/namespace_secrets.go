package core

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

func (c *Core) CreateNamespaceSecret(ctx context.Context, secret models.NamespaceSecret, namespaceID string) (models.NamespaceSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.NamespaceSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if secret.Key == "" {
		return models.NamespaceSecret{}, errors.New("secret key is required")
	}

	if secret.Value == "" {
		return models.NamespaceSecret{}, errors.New("secret value is required")
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(secret.Value))
	if err != nil {
		return models.NamespaceSecret{}, err
	}
	encryptedValue := hex.EncodeToString(enc)

	var description sql.NullString
	if secret.Description != "" {
		description = sql.NullString{String: secret.Description, Valid: true}
	}

	created, err := c.store.CreateNamespaceSecret(ctx, repo.CreateNamespaceSecretParams{
		Key:            secret.Key,
		EncryptedValue: encryptedValue,
		Description:    description,
		Uuid:           namespaceUUID,
	})
	if err != nil {
		return models.NamespaceSecret{}, err
	}

	return models.RepoNamespaceSecretToNamespaceSecret(created), nil
}

func (c *Core) GetNamespaceSecretByID(ctx context.Context, id string, namespaceID string) (models.NamespaceSecret, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.NamespaceSecret{}, err
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.NamespaceSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	secret, err := c.store.GetNamespaceSecretByUUID(ctx, repo.GetNamespaceSecretByUUIDParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.NamespaceSecret{}, err
	}

	return models.RepoNamespaceSecretByUUIDToNamespaceSecret(secret), nil
}

func (c *Core) ListNamespaceSecrets(ctx context.Context, namespaceID string) ([]models.NamespaceSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	secrets, err := c.store.ListNamespaceSecrets(ctx, namespaceUUID)
	if err != nil {
		return nil, err
	}

	return models.RepoNamespaceSecretListToNamespaceSecret(secrets), nil
}

func (c *Core) UpdateNamespaceSecret(ctx context.Context, id string, secret models.NamespaceSecret, namespaceID string) (models.NamespaceSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.NamespaceSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if secret.Value == "" {
		return models.NamespaceSecret{}, errors.New("secret value is required")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.NamespaceSecret{}, err
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(secret.Value))
	if err != nil {
		return models.NamespaceSecret{}, err
	}
	encryptedValue := hex.EncodeToString(enc)

	var description sql.NullString
	if secret.Description != "" {
		description = sql.NullString{String: secret.Description, Valid: true}
	}

	updated, err := c.store.UpdateNamespaceSecret(ctx, repo.UpdateNamespaceSecretParams{
		Uuid:           uuidID,
		Uuid_2:         namespaceUUID,
		EncryptedValue: encryptedValue,
		Description:    description,
	})
	if err != nil {
		return models.NamespaceSecret{}, err
	}

	return models.RepoNamespaceSecretToNamespaceSecret(updated), nil
}

func (c *Core) DeleteNamespaceSecret(ctx context.Context, id string, namespaceID string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid secret UUID: %w", err)
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	return c.store.DeleteNamespaceSecret(ctx, repo.DeleteNamespaceSecretParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
}

// getDecryptedNamespaceSecrets returns decrypted namespace secrets as a map
func (c *Core) getDecryptedNamespaceSecrets(ctx context.Context, namespaceID string) (map[string]string, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	secrets, err := c.store.GetDecryptedNamespaceSecrets(ctx, namespaceUUID)
	if err != nil {
		return nil, err
	}

	decryptedSecrets := make(map[string]string)
	for _, secret := range secrets {
		encryptedBytes, err := hex.DecodeString(secret.EncryptedValue)
		if err != nil {
			return nil, fmt.Errorf("could not decode encrypted value for secret %s: %w", secret.Key, err)
		}

		decryptedValue, err := c.keeper.Decrypt(ctx, encryptedBytes)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt value for secret %s: %w", secret.Key, err)
		}

		decryptedSecrets[secret.Key] = string(decryptedValue)
	}

	return decryptedSecrets, nil
}

// GetMergedSecretsForFlow returns merged namespace + flow secrets (flow overrides namespace)
// This is the SecretsProviderFn implementation that should be used by the scheduler
func (c *Core) GetMergedSecretsForFlow(ctx context.Context, flowID string, namespaceID string) (map[string]string, error) {
	merged := make(map[string]string)

	// 1. Get namespace secrets first (base layer)
	// Errors are ignored - namespace secrets might not exist or might fail to decrypt
	nsSecrets, _ := c.getDecryptedNamespaceSecrets(ctx, namespaceID)
	for k, v := range nsSecrets {
		merged[k] = v
	}

	// 2. Get flow secrets and override (flow secrets take precedence)
	// Errors are ignored - flow secrets might not exist or might fail to decrypt
	flowSecrets, _ := c.GetDecryptedFlowSecrets(ctx, flowID, namespaceID)
	for k, v := range flowSecrets {
		merged[k] = v // Flow secrets override namespace secrets with same key
	}

	return merged, nil
}
