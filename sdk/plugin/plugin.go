// Package plugin provides go-plugin infrastructure for loading and serving
// executor plugins out-of-process via gRPC with server-side streaming.
//
// Plugin authors implement executor.ExecutorPlugin and serve it via:
//
//	func main() {
//	    goplugin.Serve(&goplugin.ServeConfig{
//	        HandshakeConfig: sdkplugin.HandshakeConfig,
//	        GRPCServer:      goplugin.DefaultGRPCServer,
//	        Plugins: map[string]goplugin.Plugin{
//	            "executor": &sdkplugin.GRPCExecutorPlugin{Impl: &MyPlugin{}},
//	        },
//	    })
//	}
package plugin

import (
	"context"
	"os/exec"

	"github.com/cvhariharan/flowctl/sdk/executor"
	"github.com/cvhariharan/flowctl/sdk/plugin/proto"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	MagicCookieKey   = "FLOWCTL_EXECUTOR_PLUGIN"
	MagicCookieValue = "flowctl-executor-v1"
)

// HandshakeConfig is the go-plugin handshake configuration.
var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   MagicCookieKey,
	MagicCookieValue: MagicCookieValue,
}

// PluginMap maps the plugin name to the go-plugin Plugin implementation.
var PluginMap = map[string]goplugin.Plugin{
	"executor": &GRPCExecutorPlugin{},
}

// GRPCExecutorPlugin is the go-plugin Plugin implementation using gRPC.
// Set Impl on the server side (plugin binary); leave nil on the client side.
type GRPCExecutorPlugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl executor.ExecutorPlugin
}

func (p *GRPCExecutorPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterExecutorPluginServer(s, &grpcServer{
		impl:      p.Impl,
		executors: make(map[string]executor.Executor),
	})
	return nil
}

func (p *GRPCExecutorPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: proto.NewExecutorPluginClient(conn)}, nil
}

// LoadPlugin starts an external plugin binary and returns the go-plugin client
// and the ExecutorPlugin interface to interact with it.
func LoadPlugin(path string) (*goplugin.Client, executor.ExecutorPlugin, error) {
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []goplugin.Protocol{
			goplugin.ProtocolGRPC,
		},
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, err
	}

	raw, err := rpcClient.Dispense("executor")
	if err != nil {
		client.Kill()
		return nil, nil, err
	}

	p, ok := raw.(executor.ExecutorPlugin)
	if !ok {
		client.Kill()
		return nil, nil, err
	}

	return client, p, nil
}
