package models

import "time"

type Schedule struct {
	UUID          string                 `json:"uuid"`
	FlowSlug      string                 `json:"flow_slug"`
	FlowName      string                 `json:"flow_name"`
	Cron          string                 `json:"cron"`
	Timezone      string                 `json:"timezone"`
	Inputs        map[string]interface{} `json:"inputs"`
	CreatedByID   string                 `json:"created_by_id"`
	CreatedByName string                 `json:"created_by_name"`
	IsActive      bool                   `json:"is_active"`
	IsUserCreated bool                   `json:"is_user_created"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
