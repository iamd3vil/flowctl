package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cvhariharan/flowctl/executors/docker"
	"github.com/cvhariharan/flowctl/executors/flow"
	"github.com/cvhariharan/flowctl/executors/script"
	"github.com/cvhariharan/flowctl/internal/core"
	"github.com/cvhariharan/flowctl/sdk/executor"
	sdkplugin "github.com/cvhariharan/flowctl/sdk/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

// pluginClients holds active go-plugin clients for external plugins (for cleanup)
var pluginClients []*goplugin.Client

// registerExecutorPlugin registers a single ExecutorPlugin into the executor registries
// and generates an API token for it, returning the token.
func registerExecutorPlugin(name string, plugin executor.ExecutorPlugin, signingKey []byte) string {
	executor.RegisterExecutor(name, plugin.New)
	schema := plugin.GetSchema()
	if schema != nil {
		executor.RegisterSchema(name, schema)
	}
	executor.RegisterCapabilities(name, plugin.GetCapabilities())

	token, err := core.GenerateExecutorToken(name, signingKey)
	if err != nil {
		log.Fatalf("failed to generate token for executor %s: %v", name, err)
	}
	return token
}

// registerPlugins registers all executors and remote clients.
// It generates an API token per executor and returns them as a map.
func registerPlugins(pluginDir string, signingKey []byte) map[string]string {
	builtins := map[string]executor.ExecutorPlugin{
		"docker": &docker.DockerExecutorPlugin{},
		"script": &script.ScriptExecutorPlugin{},
		"flow":   &flow.FlowExecutorPlugin{},
	}

	executorKeys := make(map[string]string)
	for name, plugin := range builtins {
		executorKeys[name] = registerExecutorPlugin(name, plugin, signingKey)
	}

	// Load external plugins
	if pluginDir != "" {
		externalKeys := loadExternalPlugins(pluginDir, signingKey)
		for k, v := range externalKeys {
			executorKeys[k] = v
		}
	}

	return executorKeys
}

// loadExternalPlugins loads executor plugin binaries from the given directory.
func loadExternalPlugins(dir string, signingKey []byte) map[string]string {
	executorKeys := make(map[string]string)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("failed to read plugin directory %s: %v", dir, err)
		}
		return executorKeys
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		client, plugin, err := sdkplugin.LoadPlugin(path)
		if err != nil {
			log.Printf("failed to load plugin %s: %v", path, err)
			continue
		}

		name := plugin.GetName()
		if name == "" {
			log.Printf("plugin %s returned empty name, skipping", path)
			client.Kill()
			continue
		}

		pluginClients = append(pluginClients, client)
		executorKeys[name] = registerExecutorPlugin(name, plugin, signingKey)
		log.Printf("loaded external executor plugin: %s", name)
	}

	return executorKeys
}

// CleanupPlugins kills all external plugin processes.
func CleanupPlugins() {
	for _, c := range pluginClients {
		c.Kill()
	}
}
