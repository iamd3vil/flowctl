package executor

// ExecutorPlugin is implemented by all executors (built-in and external).
// go-plugin wraps this interface with gRPC for out-of-process executors.
type ExecutorPlugin interface {
	GetName() string
	GetSchema() interface{}
	GetCapabilities() Capability
	New(name string, node Node, execID string) (Executor, error)
}
