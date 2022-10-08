package trace_id

import (
	"context"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ContextKey = "trace-id"

func TraceIdFromContext(ctx context.Context) string {
	traceId, ok := ctx.Value(ContextKey).(string)
	if !ok {
		return ""
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
		newCtx := context.WithValue(ctx, ContextKey, traceId)
		grpc.SetHeader(newCtx, metadata.Pairs(
			ContextKey, traceId,
		))

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
		newCtx := context.WithValue(stream.Context(), ContextKey, traceId)
		stream.SetHeader(metadata.Pairs(
			ContextKey, traceId,
		))
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
