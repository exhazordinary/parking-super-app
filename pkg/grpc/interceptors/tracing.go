package interceptors

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// TracingUnaryServerInterceptor returns a server interceptor for OpenTelemetry tracing
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

// TracingUnaryClientInterceptor returns a client interceptor for OpenTelemetry tracing
func TracingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor()
}

// TracingStreamServerInterceptor returns a stream server interceptor for OpenTelemetry tracing
func TracingStreamServerInterceptor() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}

// TracingStreamClientInterceptor returns a stream client interceptor for OpenTelemetry tracing
func TracingStreamClientInterceptor() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor()
}
