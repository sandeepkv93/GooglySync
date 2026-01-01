package ipc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial returns a gRPC client connection over a Unix domain socket.
func Dial(ctx context.Context, socketPath string) (*grpc.ClientConn, error) {
	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", socketPath)
	}

	return grpc.DialContext(
		ctx,
		socketPath,
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
