package messengers

import (
	"fmt"
	"sync"
)

var (
	schemaRegistry = make(map[string]interface{})
	smu            sync.RWMutex
)

// RegisterSchema registers a messenger's config schema. Panics on duplicate.
func RegisterSchema(name string, schema interface{}) {
	smu.Lock()
	defer smu.Unlock()

	if _, exists := schemaRegistry[name]; exists {
		panic(fmt.Sprintf("messenger schema '%s' is already registered", name))
	}
	schemaRegistry[name] = schema
}

// GetAllSchemas returns a map of all registered messenger names to their config schemas.
func GetAllSchemas() map[string]interface{} {
	smu.RLock()
	defer smu.RUnlock()

	result := make(map[string]interface{}, len(schemaRegistry))
	for k, v := range schemaRegistry {
		result[k] = v
	}
	return result
}
