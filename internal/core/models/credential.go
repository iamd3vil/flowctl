package models

const TimeFormat = "2006-01-02T15:04:05Z"

type Credential struct {
	ID            string
	Name          string
	KeyType       string
	KeyData       string
	NamespaceUUID string
	LastAccessed  string
}
