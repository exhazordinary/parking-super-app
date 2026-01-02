package interceptors

import (
	"google.golang.org/grpc"
)

// DefaultServerInterceptors returns the recommended chain of server interceptors
// Order: Recovery -> Tracing -> Logging
// Recovery is first to catch panics from all other interceptors
func DefaultServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		RecoveryUnaryServerInterceptor(),
		TracingUnaryServerInterceptor(),
		LoggingUnaryServerInterceptor(),
	}
}

// DefaultStreamServerInterceptors returns the recommended chain of stream server interceptors
func DefaultStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		RecoveryStreamServerInterceptor(),
		TracingStreamServerInterceptor(),
		LoggingStreamServerInterceptor(),
	}
}

// DefaultClientInterceptors returns the recommended chain of client interceptors
// Order: Tracing -> Logging
func DefaultClientInterceptors() []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{
		TracingUnaryClientInterceptor(),
		LoggingUnaryClientInterceptor(),
	}
}

// DefaultStreamClientInterceptors returns the recommended chain of stream client interceptors
func DefaultStreamClientInterceptors() []grpc.StreamClientInterceptor {
	return []grpc.StreamClientInterceptor{
		TracingStreamClientInterceptor(),
	}
}

// NewServerWithDefaults creates a gRPC server with default interceptors
func NewServerWithDefaults(opts ...grpc.ServerOption) *grpc.Server {
	serverOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(DefaultServerInterceptors()...),
		grpc.ChainStreamInterceptor(DefaultStreamServerInterceptors()...),
	}
	serverOpts = append(serverOpts, opts...)
	return grpc.NewServer(serverOpts...)
}
