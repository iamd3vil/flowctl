package models

import "time"

type Namespace struct {
	ID        int32     `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupNamespaceAccess struct {
	ID            int32     `json:"id"`
	GroupUUID     string    `json:"group_uuid"`
	NamespaceUUID string    `json:"namespace_uuid"`
	CreatedAt     time.Time `json:"created_at"`
}