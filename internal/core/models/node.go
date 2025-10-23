package models

type AuthMethod string

const (
	AuthMethodPrivateKey AuthMethod = "private_key"
	AuthMethodPassword   AuthMethod = "password"
)

type Node struct {
	ID             string
	Name           string
	Hostname       string
	Port           int
	Username       string
	OSFamily       string
	ConnectionType string
	Tags           []string
	Auth           NodeAuth
	NamespaceUUID  string
}

type NodeAuth struct {
	CredentialID string
	Method       AuthMethod
	Key          string
}

type NodeStats struct {
	TotalHosts int64 `json:"total_hosts"`
	SSHHosts   int64 `json:"ssh_hosts"`
	QSSHHosts  int64 `json:"qssh_hosts"`
}
