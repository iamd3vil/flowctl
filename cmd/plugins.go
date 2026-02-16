package cmd

import (
	"github.com/cvhariharan/flowctl/executors/docker"
	"github.com/cvhariharan/flowctl/executors/script"
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

func registerPlugins() {
	executors["docker"] = executorDef{New: docker.NewDockerExecutor, Schema: docker.GetSchema()}
	executors["script"] = executorDef{New: script.NewScriptExecutor, Schema: script.GetSchema()}

	for name, e := range executors {
		executor.RegisterExecutor(name, e.New)
		if e.Schema != nil {
			executor.RegisterSchema(name, e.Schema)
		}
	}

	remoteClients["ssh"] = sshclient.NewRemoteClient
	remoteClients["qssh"] = qsshclient.NewRemoteClient

	for name, newClient := range remoteClients {
		remoteclient.Register(name, newClient)
	}
}
