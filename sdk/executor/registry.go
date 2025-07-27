package executor

import (
	"fmt"
	"sync"
)

// NewExecutorFunc defines the signature for a function that can create an executor.
type NewExecutorFunc func(name string, node Node) (Executor, error)

var (
	registry = make(map[string]NewExecutorFunc)
	mu       sync.RWMutex
)

// RegisterExecutor should be called by executor modules their init() function.
func RegisterExecutor(name string, factory NewExecutorFunc) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("executor with name '%s' is already registered", name))
	}
	registry[name] = factory
}

// GetNewExecutorFunc is used to retrieve the NewExecutorFunc for an executor
func GetNewExecutorFunc(name string) (NewExecutorFunc, error) {
	mu.RLock()
	defer mu.RUnlock()

	f, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("executor '%s' is not registered or not included in the build", name)
	}
	return f, nil
}
