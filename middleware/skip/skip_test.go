package skip

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"

	"github.com/zmzhang8/grpc_example/test"
)

func dummyUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		return nil, errors.New("interceptor")
	}
}

func dummyStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return errors.New("interceptor")
	}
}

func TestUnaryServerInterceptor_skip(t *testing.T) {
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{FullMethod: "/grpc.DummyService/DummyMethod"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	skipFunc := func(ctx context.Context, service string, method string) bool {
		return service == "grpc.DummyService"
	}

	_, err := UnaryServerInterceptor(dummyUnaryServerInterceptor(), skipFunc)(ctx, nil, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}

func TestUnaryServerInterceptor_noSkip(t *testing.T) {
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{FullMethod: "/grpc.DummyService/DummyMethod"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	skipFunc := func(ctx context.Context, service string, method string) bool {
		return service == "grpc.OtherService"
	}

	_, err := UnaryServerInterceptor(dummyUnaryServerInterceptor(), skipFunc)(ctx, nil, &info, handler)

	if err.Error() != "interceptor" {
		t.Errorf("err %v; want interceptor", err)
	}
}

func TestStreamServerInterceptor_skip(t *testing.T) {
	ctx := context.TODO()
	stream := test.ServerStreamMock{Ctx: ctx}
	info := grpc.StreamServerInfo{FullMethod: "/grpc.DummyService/DummyMethod"}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
	skipFunc := func(ctx context.Context, service string, method string) bool {
		return service == "grpc.DummyService"
	}

	err := StreamServerInterceptor(dummyStreamServerInterceptor(), skipFunc)(nil, stream, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}

func TestStreamServerInterceptor_noSkip(t *testing.T) {
	ctx := context.TODO()
	stream := test.ServerStreamMock{Ctx: ctx}
	info := grpc.StreamServerInfo{FullMethod: "/grpc.DummyService/DummyMethod"}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
	skipFunc := func(ctx context.Context, service string, method string) bool {
		return service == "grpc.OtherService"
	}

	err := StreamServerInterceptor(dummyStreamServerInterceptor(), skipFunc)(nil, stream, &info, handler)

	if err.Error() != "interceptor" {
		t.Errorf("err %v; want interceptor", err)
	}
}
