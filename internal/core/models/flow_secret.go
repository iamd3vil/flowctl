package models

import (
	"github.com/cvhariharan/flowctl/internal/repo"
)

type FlowSecret struct {
	ID            string
	FlowID        int32
	Key           string
	Value         string // Only populated when creating/updating, not when listing
	Description   string
	NamespaceUUID string
	CreatedAt     string
	UpdatedAt     string
}

func RepoFlowSecretToFlowSecret(repoSecret repo.FlowSecret) FlowSecret {
	var description string
	if repoSecret.Description.Valid {
		description = repoSecret.Description.String
	}

	return FlowSecret{
		ID:          repoSecret.Uuid.String(),
		FlowID:      repoSecret.FlowID,
		Key:         repoSecret.Key,
		Description: description,
		CreatedAt:   repoSecret.CreatedAt.Format(TimeFormat),
		UpdatedAt:   repoSecret.UpdatedAt.Format(TimeFormat),
	}
}

func RepoFlowSecretByUUIDToFlowSecret(repoSecret repo.GetFlowSecretByUUIDRow) FlowSecret {
	var description string
	if repoSecret.Description.Valid {
		description = repoSecret.Description.String
	}

	return FlowSecret{
		ID:            repoSecret.Uuid.String(),
		FlowID:        repoSecret.FlowID,
		Key:           repoSecret.Key,
		Description:   description,
		NamespaceUUID: repoSecret.NamespaceUuid.String(),
		CreatedAt:     repoSecret.CreatedAt.Format(TimeFormat),
		UpdatedAt:     repoSecret.UpdatedAt.Format(TimeFormat),
	}
}

func RepoFlowSecretListToFlowSecret(repoSecrets []repo.ListFlowSecretsRow) []FlowSecret {
	results := make([]FlowSecret, 0)
	for _, secret := range repoSecrets {
		var description string
		if secret.Description.Valid {
			description = secret.Description.String
		}

		results = append(results, FlowSecret{
			ID:            secret.Uuid.String(),
			FlowID:        secret.FlowID,
			Key:           secret.Key,
			Description:   description,
			NamespaceUUID: secret.NamespaceUuid.String(),
			CreatedAt:     secret.CreatedAt.Format(TimeFormat),
			UpdatedAt:     secret.UpdatedAt.Format(TimeFormat),
		})
	}
	return results
}
