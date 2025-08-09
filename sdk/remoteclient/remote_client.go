package remoteclient

import (
	"context"
	"io"
	"net"
)

// RemoteClient defines an interface for interacting with a remote machine.
// This abstraction allows for swapping the underlying client implementation
type RemoteClient interface {
	// RunCommand executes a command on the remote machine
	RunCommand(ctx context.Context, command string, stdout io.Writer, stderr io.Writer) error
	// Download copies a file from the remote path to a local path
	Download(ctx context.Context, remotePath, localPath string) error
	// Upload copies a file from the local path to a remote path
	Upload(ctx context.Context, localPath, remotePath string) error
	// Dial opens a connection to the given network and address on the remote machine.
	Dial(network, address string) (net.Conn, error)
	// Close terminates the connection to the remote machine.
	Close() error
}
