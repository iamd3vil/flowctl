package remoteclient

import (
	"net"
)

// RemoteClient defines an interface for interacting with a remote machine.
// This abstraction allows for swapping the underlying client implementation
type RemoteClient interface {
	// RunCommand executes a command on the remote machine and returns the combined output.
	RunCommand(command string) (string, error)
	// Download copies a file from the remote path to a local path.
	Download(remotePath, localPath string) error
	// Upload copies a file from the local path to a remote path.
	Upload(localPath, remotePath string) error
	// Dial opens a connection to the given network and address on the remote machine.
	Dial(network, address string) (net.Conn, error)
	// Close terminates the connection to the remote machine.
	Close() error
}
