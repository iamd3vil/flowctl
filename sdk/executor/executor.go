package executor

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Node struct {
	Hostname       string
	Port           int
	Username       string
	Auth           NodeAuth
	ConnectionType string
}

type NodeAuth struct {
	Method string
	Key    string
}

type ExecutionContext struct {
	// WithConfig is the yaml config passed to the executor
	WithConfig []byte
	Artifacts  []string
	Inputs     map[string]any
	Stdout     io.Writer
	Stderr     io.Writer
}

// CheckConnectivity can be used to check if a remote node is accessible at the given IP:Port
// The default connection timeout is 5 seconds
// Non-nil error is returned if the node is not accessible
func (n *Node) CheckConnectivity() error {
	address := fmt.Sprintf("%s:%d", n.Hostname, n.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	defer conn.Close()
	return nil
}

type Executor interface {
	Execute(ctx context.Context, execCtx ExecutionContext) (outputs map[string]string, err error)
	PushFile(ctx context.Context, localFilePath string, remoteFilePath string) error
	PullFile(ctx context.Context, remoteFilePath string, localFilePath string) error
	Close() error
}
