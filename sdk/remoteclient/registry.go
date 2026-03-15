package remoteclient

import (
	"fmt"
)

// NodeConfig contains the configuration needed to connect to a remote node
type NodeConfig struct {
	Hostname string
	Port     int
	Username string
	Auth     NodeAuth
}

// NodeAuth contains authentication information for a node
type NodeAuth struct {
	Method string
	Key    string
}

// NewRemoteClientFunc defines the signature for creating a new RemoteClient.
type NewRemoteClientFunc func(config NodeConfig) (RemoteClient, error)

var registry = map[string]NewRemoteClientFunc{
	"ssh":  newSSHClient,
	"qssh": newQSSHClient,
}

// GetClient is called by executors to get a client for a specific protocol.
func GetClient(protocolName string, config NodeConfig) (RemoteClient, error) {
	factory, ok := registry[protocolName]
	if !ok {
		return nil, fmt.Errorf("remote client for protocol '%s' is not registered", protocolName)
	}
	return factory(config)
}
