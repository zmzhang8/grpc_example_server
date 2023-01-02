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

func TestMustGetLogger_success(t *testing.T) {
	wantLogger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), contextKey{}, wantLogger)

	gotLogger := MustGetLogger(ctx)

	if gotLogger == nil {
		t.Errorf("logger <nil>; want %v", wantLogger)
	}
}

func TestMustGetLogger_failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panicked false; want true")
		}
	}()

	MustGetLogger(context.TODO())
}

func TestUnaryServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if _, ok := ctx.Value(contextKey{}).(log.Logger); !ok {
			return nil, errors.New("")
		}
		return nil, nil
	}

	_, err := UnaryServerInterceptor(logger, nil)(ctx, nil, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	stream := test.ServerStreamMock{Ctx: context.TODO()}
	info := grpc.StreamServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		if _, ok := stream.Context().Value(contextKey{}).(log.Logger); !ok {
			return errors.New("")
		}
		return nil
	}

	err := StreamServerInterceptor(logger, nil)(nil, stream, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}
