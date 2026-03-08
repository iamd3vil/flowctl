package executor

import (
	"context"
	"io"
)

type Node struct {
	Hostname       string
	Port           int
	Username       string
	Auth           NodeAuth
	ConnectionType string
	OSFamily       string
}

type NodeAuth struct {
	Method string
	Key    string
}

type ExecutionContext struct {
	// WithConfig is the yaml config passed to the executor
	WithConfig    []byte
	Inputs        map[string]any
	Stdout        io.Writer
	Stderr        io.Writer
	UserUUID      string
	NamespaceName string // human-readable namespace name for API calls
	APIKey        string // executor API key for authenticating with the server
	APIBaseURL    string // server base URL for API calls
}

type Capability uint64

const (
	RemoteExecution Capability = 1 << iota
	EnvironmentVariables
	FileTransfer
	StreamingOutput
)

type Executor interface {
	Execute(ctx context.Context, execCtx ExecutionContext) (outputs map[string]string, err error)
	GetArtifactsDir() string
	Close() error
}
