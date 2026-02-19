package messengers

import (
	"fmt"
	"sync"
)

var (
	schemaRegistry = make(map[string]any)
	smu            sync.RWMutex
)

// RegisterSchema registers a messenger's config schema. Panics on duplicate.
func RegisterSchema(name string, schema any) {
	smu.Lock()
	defer smu.Unlock()

	if _, exists := schemaRegistry[name]; exists {
		panic(fmt.Sprintf("messenger schema '%s' is already registered", name))
	}
	schemaRegistry[name] = schema
}

// GetAllSchemas returns a map of all registered messenger names to their config schemas.
func GetAllSchemas() map[string]any {
	smu.RLock()
	defer smu.RUnlock()

	result := make(map[string]any, len(schemaRegistry))
	for k, v := range schemaRegistry {
		result[k] = v
	}
	return result
}
