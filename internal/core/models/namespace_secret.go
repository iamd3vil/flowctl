package models

import (
	"github.com/cvhariharan/flowctl/internal/repo"
)

type NamespaceSecret struct {
	ID            string
	Key           string
	Value         string // Only populated when creating/updating, not when listing
	Description   string
	NamespaceUUID string
	CreatedAt     string
	UpdatedAt     string
}

func RepoNamespaceSecretToNamespaceSecret(repoSecret repo.NamespaceSecret) NamespaceSecret {
	var description string
	if repoSecret.Description.Valid {
		description = repoSecret.Description.String
	}

	return NamespaceSecret{
		ID:          repoSecret.Uuid.String(),
		Key:         repoSecret.Key,
		Description: description,
		CreatedAt:   repoSecret.CreatedAt.Format(TimeFormat),
		UpdatedAt:   repoSecret.UpdatedAt.Format(TimeFormat),
	}
}

func RepoNamespaceSecretByUUIDToNamespaceSecret(repoSecret repo.GetNamespaceSecretByUUIDRow) NamespaceSecret {
	var description string
	if repoSecret.Description.Valid {
		description = repoSecret.Description.String
	}

	return NamespaceSecret{
		ID:            repoSecret.Uuid.String(),
		Key:           repoSecret.Key,
		Description:   description,
		NamespaceUUID: repoSecret.NamespaceUuid.String(),
		CreatedAt:     repoSecret.CreatedAt.Format(TimeFormat),
		UpdatedAt:     repoSecret.UpdatedAt.Format(TimeFormat),
	}
}

func RepoNamespaceSecretListToNamespaceSecret(repoSecrets []repo.ListNamespaceSecretsRow) []NamespaceSecret {
	results := make([]NamespaceSecret, 0)
	for _, secret := range repoSecrets {
		var description string
		if secret.Description.Valid {
			description = secret.Description.String
		}

		results = append(results, NamespaceSecret{
			ID:            secret.Uuid.String(),
			Key:           secret.Key,
			Description:   description,
			NamespaceUUID: secret.NamespaceUuid.String(),
			CreatedAt:     secret.CreatedAt.Format(TimeFormat),
			UpdatedAt:     secret.UpdatedAt.Format(TimeFormat),
		})
	}
	return results
}
