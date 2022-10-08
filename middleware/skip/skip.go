package skip

import (
	"context"
	"path"

	"google.golang.org/grpc"
)

type SkipFunc func(ctx context.Context, service string, method string) bool

func UnaryServerInterceptor(in grpc.UnaryServerInterceptor, skipFunc SkipFunc) grpc.UnaryServerInterceptor {
	if skipFunc == nil {
		return in
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		service, method := splitServiceMethod(info.FullMethod)
		if skipFunc(ctx, service, method) {
			return handler(ctx, req)
		} else {
			return in(ctx, req, info, handler)
		}
	}
}

func StreamServerInterceptor(in grpc.StreamServerInterceptor, skipFunc SkipFunc) grpc.StreamServerInterceptor {
	if skipFunc == nil {
		return in
	}

	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		service, method := splitServiceMethod(info.FullMethod)
		if skipFunc(stream.Context(), service, method) {
			return handler(srv, stream)
		} else {
			return in(srv, stream, info, handler)
		}
	}
}

func splitServiceMethod(fullMethod string) (string, string) {
	service := path.Dir(fullMethod)[1:]
	method := path.Base(fullMethod)
	return service, method
}
