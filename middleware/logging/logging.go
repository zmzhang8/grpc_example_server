package logging

import (
	"context"
	"path"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/log"
)

type contextKey struct{}

type LoggerFunc func(ctx context.Context, logger log.Logger) log.Logger

func ContextKey() contextKey {
	return contextKey{}
}

func MustGetLogger(ctx context.Context) log.Logger {
	logger, ok := ctx.Value(contextKey{}).(log.Logger)
	if !ok {
		panic("cannot get logger in context")
	}
	return logger
}

func UnaryServerInterceptor(logger log.Logger, loggerFunc LoggerFunc) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now().UTC()
		contextLogger := logger
		if loggerFunc != nil {
			contextLogger = loggerFunc(ctx, logger)
		}
		newCtx := context.WithValue(ctx, contextKey{}, contextLogger)

		service, method := splitServiceMethod(info.FullMethod)
		stats := []interface{}{
			"grpc.service", service,
			"grpc.method", method,
			"grpc.start_time", startTime,
		}

		resp, err := handler(newCtx, req)

		code := status.Code(err)
		duration := float64(time.Since(startTime)) / float64(time.Millisecond)
		stats = append(stats,
			"grpc.code", code,
			"grpc.duration_ms", duration,
		)
		logwCodeToLevel(contextLogger, code, "Finished unary call", stats...)

		return resp, err
	}
}

func StreamServerInterceptor(logger log.Logger, loggerFunc LoggerFunc) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now().UTC()
		ctx := stream.Context()
		contextLogger := logger
		if loggerFunc != nil {
			contextLogger = loggerFunc(ctx, logger)
		}
		newCtx := context.WithValue(ctx, contextKey{}, contextLogger)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		service, method := splitServiceMethod(info.FullMethod)
		stats := []interface{}{
			"grpc.service", service,
			"grpc.method", method,
			"grpc.start_time", startTime,
		}

		err := handler(srv, wrapped)

		code := status.Code(err)
		duration := float64(time.Since(startTime)) / float64(time.Millisecond)
		stats = append(stats,
			"grpc.code", code,
			"grpc.duration_ms", duration,
		)
		logwCodeToLevel(contextLogger, code, "Finished stream call", stats...)

		return err
	}
}

func splitServiceMethod(fullMethod string) (string, string) {
	service := path.Dir(fullMethod)[1:]
	method := path.Base(fullMethod)
	return service, method
}

func logwCodeToLevel(
	logger log.Logger,
	code codes.Code,
	msg string,
	keysAndValues ...interface{},
) {
	switch code {
	case codes.OK:
		logger.Infow(msg, keysAndValues...)
	case codes.Canceled:
		logger.Infow(msg, keysAndValues...)
	case codes.Unknown:
		logger.Errorw(msg, keysAndValues...)
	case codes.InvalidArgument:
		logger.Infow(msg, keysAndValues...)
	case codes.DeadlineExceeded:
		logger.Warnw(msg, keysAndValues...)
	case codes.NotFound:
		logger.Infow(msg, keysAndValues...)
	case codes.AlreadyExists:
		logger.Infow(msg, keysAndValues...)
	case codes.PermissionDenied:
		logger.Warnw(msg, keysAndValues...)
	case codes.Unauthenticated:
		logger.Infow(msg, keysAndValues...) // unauthenticated requests can happen
	case codes.ResourceExhausted:
		logger.Warnw(msg, keysAndValues...)
	case codes.FailedPrecondition:
		logger.Warnw(msg, keysAndValues...)
	case codes.Aborted:
		logger.Warnw(msg, keysAndValues...)
	case codes.OutOfRange:
		logger.Warnw(msg, keysAndValues...)
	case codes.Unimplemented:
		logger.Errorw(msg, keysAndValues...)
	case codes.Internal:
		logger.Errorw(msg, keysAndValues...)
	case codes.Unavailable:
		logger.Warnw(msg, keysAndValues...)
	case codes.DataLoss:
		logger.Errorw(msg, keysAndValues...)
	default:
		logger.Errorw(msg, keysAndValues...)
	}
}
