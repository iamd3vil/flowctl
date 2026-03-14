package executor

import (
	"fmt"
	"regexp"
	"sync"
)

// NewExecutorFunc defines the signature for a function that can create an executor.
type NewExecutorFunc func(name string, node Node, execID string) (Executor, error)

var (
	registry       = make(map[string]NewExecutorFunc)
	schemaRegistry = make(map[string]interface{})
	capRegistry    = make(map[string]Capability)
	mu             sync.RWMutex
	smu            sync.RWMutex
	cmu            sync.RWMutex
)

var validNameRegex = regexp.MustCompile(`^[a-zA-Z_]+$`)

func isValidName(name string) bool {
	return validNameRegex.MatchString(name)
}

// RegisterExecutor registers an executor new func under the given name.
func RegisterExecutor(name string, factory NewExecutorFunc) {
	mu.Lock()
	defer mu.Unlock()

	if !isValidName(name) {
		panic("executor name can only include alphabets and underscore")
	}

	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("executor with name '%s' is already registered", name))
	}
	registry[name] = factory
}

// RegisterSchema should be called by executor modules to register their config schema
func RegisterSchema(name string, schema interface{}) {
	smu.Lock()
	defer smu.Unlock()

	if !isValidName(name) {
		panic("executor name can only include alphabets and underscore")
	}

	mu.RLock()
	if _, exists := registry[name]; !exists {
		panic(fmt.Sprintf("executor '%s' is not registered, cannot register schema", name))
	}
	mu.RUnlock()

	if _, exists := schemaRegistry[name]; exists {
		panic(fmt.Sprintf("schema with name '%s' is already registered", name))
	}
	schemaRegistry[name] = schema
}

// RegisterCapabilities registers the capabilities for a named executor.
func RegisterCapabilities(name string, caps Capability) {
	cmu.Lock()
	defer cmu.Unlock()

	mu.RLock()
	if _, exists := registry[name]; !exists {
		panic(fmt.Sprintf("executor '%s' is not registered, cannot register capabilities", name))
	}
	mu.RUnlock()

	capRegistry[name] = caps
}

// GetCapabilities returns the capabilities of a named executor.
func GetCapabilities(name string) (Capability, error) {
	cmu.RLock()
	defer cmu.RUnlock()

	caps, ok := capRegistry[name]
	if !ok {
		return 0, fmt.Errorf("capabilities for executor '%s' are not registered", name)
	}
	return caps, nil
}

// GetSchema returns the config schema of the executor
func GetSchema(name string) (interface{}, error) {
	smu.RLock()
	defer smu.RUnlock()

	schema, ok := schemaRegistry[name]
	if !ok {
		return nil, fmt.Errorf("schema for executor '%s' is not registered", name)
	}
	return schema, nil
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

type ExecutorEntry struct {
	Name         string
	Capabilities []string
}

func GetAllExecutors() []ExecutorEntry {
	mu.RLock()
	defer mu.RUnlock()

	entries := make([]ExecutorEntry, 0, len(registry))
	for name := range registry {
		cmu.RLock()
		caps := capRegistry[name]
		cmu.RUnlock()
		entries = append(entries, ExecutorEntry{
			Name:         name,
			Capabilities: CapabilityStrings(caps),
		})
	}
	return entries
}
