package logging

import (
	"context"
	"errors"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/zmzhang8/grpc_example/lib/log"
	"github.com/zmzhang8/grpc_example/test"
)

func TestLoggerFromContext_success(t *testing.T) {
	wantLogger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), ContextKey, wantLogger)

	gotLogger := LoggerFromContext(ctx)

	if gotLogger == nil {
		t.Errorf("logger <nil>; want %v", wantLogger)
	}
}

func TestLoggerFromContext_failure(t *testing.T) {
	gotLogger := LoggerFromContext(context.TODO())

	if gotLogger != nil {
		t.Errorf("logger %v; want <nil>", gotLogger)
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), "request-id", "")
	info := grpc.UnaryServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "yyy", errors.New("zzz")
	}

	resp, err := UnaryServerInterceptor(logger)(ctx, nil, &info, handler)

	if resp != "yyy" {
		t.Errorf("resp %v; want yyy", resp)
	}
	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), "request-id", "")
	stream := test.ServerStreamMock{Ctx: ctx}
	info := grpc.StreamServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New("zzz")
	}

	err := StreamServerInterceptor(logger)(nil, stream, &info, handler)

	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}
