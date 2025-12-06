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

func (c *Core) CreateFlowSecret(ctx context.Context, flowID string, secret models.FlowSecret, namespaceID string) (models.FlowSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	// Get the flow to get its database ID
	flow, err := c.GetFlowByID(flowID, namespaceID)
	if err != nil {
		return models.FlowSecret{}, fmt.Errorf("flow not found: %w", err)
	}

	if secret.Key == "" {
		return models.FlowSecret{}, errors.New("secret key is required")
	}

	if secret.Value == "" {
		return models.FlowSecret{}, errors.New("secret value is required")
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(secret.Value))
	if err != nil {
		return models.FlowSecret{}, err
	}
	encryptedValue := hex.EncodeToString(enc)

	var description sql.NullString
	if secret.Description != "" {
		description = sql.NullString{String: secret.Description, Valid: true}
	}

	created, err := c.store.CreateFlowSecret(ctx, repo.CreateFlowSecretParams{
		FlowID:         flow.Meta.DBID,
		Key:            secret.Key,
		EncryptedValue: encryptedValue,
		Description:    description,
		Uuid:           namespaceUUID,
	})
	if err != nil {
		return models.FlowSecret{}, err
	}

	return models.RepoFlowSecretToFlowSecret(created), nil
}

func (c *Core) GetFlowSecretByID(ctx context.Context, id string, namespaceID string) (models.FlowSecret, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.FlowSecret{}, err
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	secret, err := c.store.GetFlowSecretByUUID(ctx, repo.GetFlowSecretByUUIDParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
	if err != nil {
		return models.FlowSecret{}, err
	}

	return models.RepoFlowSecretByUUIDToFlowSecret(secret), nil
}

func (c *Core) ListFlowSecrets(ctx context.Context, flowID string, namespaceID string) ([]models.FlowSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	// Get the flow to get its database ID
	flow, err := c.GetFlowByID(flowID, namespaceID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	secrets, err := c.store.ListFlowSecrets(ctx, repo.ListFlowSecretsParams{
		FlowID: flow.Meta.DBID,
		Uuid:   namespaceUUID,
	})
	if err != nil {
		return nil, err
	}

	return models.RepoFlowSecretListToFlowSecret(secrets), nil
}

func (c *Core) UpdateFlowSecret(ctx context.Context, id string, secret models.FlowSecret, namespaceID string) (models.FlowSecret, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return models.FlowSecret{}, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	if secret.Value == "" {
		return models.FlowSecret{}, errors.New("secret value is required")
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.FlowSecret{}, err
	}

	enc, err := c.keeper.Encrypt(ctx, []byte(secret.Value))
	if err != nil {
		return models.FlowSecret{}, err
	}
	encryptedValue := hex.EncodeToString(enc)

	var description sql.NullString
	if secret.Description != "" {
		description = sql.NullString{String: secret.Description, Valid: true}
	}

	updated, err := c.store.UpdateFlowSecret(ctx, repo.UpdateFlowSecretParams{
		Uuid:           uuidID,
		Uuid_2:         namespaceUUID,
		EncryptedValue: encryptedValue,
		Description:    description,
	})
	if err != nil {
		return models.FlowSecret{}, err
	}

	return models.RepoFlowSecretToFlowSecret(updated), nil
}

func (c *Core) DeleteFlowSecret(ctx context.Context, id string, namespaceID string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid secret UUID: %w", err)
	}

	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	return c.store.DeleteFlowSecret(ctx, repo.DeleteFlowSecretParams{
		Uuid:   uuidID,
		Uuid_2: namespaceUUID,
	})
}

func (c *Core) GetDecryptedFlowSecrets(ctx context.Context, flowID string, namespaceID string) (map[string]string, error) {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %w", err)
	}

	// Get the flow to get its database ID
	flow, err := c.GetFlowByID(flowID, namespaceID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	secrets, err := c.store.ListFlowSecrets(ctx, repo.ListFlowSecretsParams{
		FlowID: flow.Meta.DBID,
		Uuid:   namespaceUUID,
	})
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
