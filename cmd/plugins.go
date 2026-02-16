package cmd

import (
	"log"

	"github.com/cvhariharan/flowctl/executors/docker"
	"github.com/cvhariharan/flowctl/executors/flow"
	"github.com/cvhariharan/flowctl/executors/script"
	"github.com/cvhariharan/flowctl/internal/core"
	qsshclient "github.com/cvhariharan/flowctl/remoteclients/qssh"
	sshclient "github.com/cvhariharan/flowctl/remoteclients/ssh"
	"github.com/cvhariharan/flowctl/sdk/executor"
	"github.com/cvhariharan/flowctl/sdk/remoteclient"
)

type executorDef struct {
	New    executor.NewExecutorFunc
	Schema interface{}
}

var executors = make(map[string]executorDef)
var remoteClients = make(map[string]remoteclient.NewRemoteClientFunc)

// registerPlugins registers all executors and remote clients.
// It generates an API token per executor and returns them as a map.
func registerPlugins(signingKey []byte) map[string]string {
	executors["docker"] = executorDef{New: docker.NewDockerExecutor, Schema: docker.GetSchema()}
	executors["script"] = executorDef{New: script.NewScriptExecutor, Schema: script.GetSchema()}
	executors["flow"] = executorDef{New: flow.NewFlowExecutor, Schema: flow.GetSchema()}

	executorKeys := make(map[string]string)
	for name, e := range executors {
		executor.RegisterExecutor(name, e.New)
		if e.Schema != nil {
			executor.RegisterSchema(name, e.Schema)
		}

		token, err := core.GenerateExecutorToken(name, signingKey)
		if err != nil {
			log.Fatalf("failed to generate token for executor %s: %v", name, err)
		}
		executorKeys[name] = token
	}

	remoteClients["ssh"] = sshclient.NewRemoteClient
	remoteClients["qssh"] = qsshclient.NewRemoteClient

	for name, newClient := range remoteClients {
		remoteclient.Register(name, newClient)
	}

	return executorKeys
}
