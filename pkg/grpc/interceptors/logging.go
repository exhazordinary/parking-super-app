package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingUnaryServerInterceptor returns a server interceptor that logs requests
func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		log.Printf("gRPC method=%s duration=%v status=%s",
			info.FullMethod,
			duration,
			st.Code().String(),
		)

		return resp, err
	}
}

// LoggingUnaryClientInterceptor returns a client interceptor that logs requests
func LoggingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		log.Printf("gRPC client method=%s duration=%v status=%s",
			method,
			duration,
			st.Code().String(),
		)

		return err
	}
}

// LoggingStreamServerInterceptor returns a stream server interceptor that logs requests
func LoggingStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		log.Printf("gRPC stream method=%s duration=%v status=%s",
			info.FullMethod,
			duration,
			st.Code().String(),
		)

		return err
	}
}
