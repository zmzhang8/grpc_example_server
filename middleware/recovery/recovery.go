package recovery

import (
	"context"
	"path"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/log"
)

func UnaryServerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				contextLogger := getLoggerWithContextFields(logger, ctx)
				service, method := splitServiceMethod(info.FullMethod)
				stats := []interface{}{
					"grpc.service", service,
					"grpc.method", method,
					"grpc.time", time.Now().UTC(),
					"error", r,
					"stack", string(debug.Stack()),
				}
				contextLogger.Errorw("Recover from panic", stats...)
				err = status.Errorf(codes.Internal, "panicked")
			}
		}()

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func StreamServerInterceptor(logger log.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				contextLogger := getLoggerWithContextFields(logger, stream.Context())
				service, method := splitServiceMethod(info.FullMethod)
				stats := []interface{}{
					"grpc.service", service,
					"grpc.method", method,
					"grpc.time", time.Now().UTC(),
					"error", r,
					"stack", string(debug.Stack()),
				}
				contextLogger.Errorw("Recover from panic", stats...)
				err = status.Errorf(codes.Internal, "panicked")
			}
		}()

		err = handler(srv, stream)
		return err
	}
}

func splitServiceMethod(fullMethod string) (string, string) {
	service := path.Dir(fullMethod)[1:]
	method := path.Base(fullMethod)
	return service, method
}

func getLoggerWithContextFields(logger log.Logger, ctx context.Context) log.Logger {
	ctxFields := []string{"trace-id"}
	args := make([]interface{}, 0)
	for _, key := range ctxFields {
		value, ok := ctx.Value(key).(string)
		if ok && value != "" {
			args = append(args, key)
			args = append(args, value)
		}
	}

	return logger.With(args...)
}
