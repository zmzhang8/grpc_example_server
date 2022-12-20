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

const ContextKey = "logger"

func MustGetLogger(ctx context.Context) log.Logger {
	loggerValue := ctx.Value(ContextKey)
	if loggerValue == nil {
		panic("logger doesn't exist in context")
	}
	logger, ok := loggerValue.(log.Logger)
	if !ok {
		panic("bad logger in context")
	}
	return logger
}

func UnaryServerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now().UTC()
		contextLogger := getLoggerWithContextFields(logger, ctx)
		newCtx := context.WithValue(ctx, ContextKey, contextLogger)

		service, method := splitServiceMethod(info.FullMethod)
		stats := []interface{}{
			"grpc.service", service,
			"grpc.method", method,
			"grpc.start_time", startTime,
		}
		contextLogger.Infow("Started unary call", stats...)

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

func StreamServerInterceptor(logger log.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now().UTC()
		ctx := stream.Context()
		contextLogger := getLoggerWithContextFields(logger, ctx)
		newCtx := context.WithValue(ctx, ContextKey, contextLogger)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		service, method := splitServiceMethod(info.FullMethod)
		stats := []interface{}{
			"grpc.service", service,
			"grpc.method", method,
			"grpc.start_time", startTime,
		}
		contextLogger.Infow("Started stream call", stats...)

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
