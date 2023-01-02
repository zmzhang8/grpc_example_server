package trace_id

import (
	"context"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey struct{}

func MustGetTraceID(ctx context.Context) string {
	traceId, ok := ctx.Value(contextKey{}).(string)
	if !ok {
		panic("cannot get trace id in context")
	}
	return traceId
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		traceId := uuid.NewString()
		newCtx := context.WithValue(ctx, contextKey{}, traceId)
		grpc.SetHeader(newCtx, metadata.Pairs("trace-id", traceId))

		return handler(newCtx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		traceId := uuid.NewString()
		newCtx := context.WithValue(stream.Context(), contextKey{}, traceId)
		stream.SetHeader(metadata.Pairs("trace-id", traceId))
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
