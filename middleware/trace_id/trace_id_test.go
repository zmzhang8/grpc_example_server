package trace_id

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"

	"github.com/zmzhang8/grpc_example/test"
)

func TestGetTraceID_success(t *testing.T) {
	ctx := context.WithValue(context.TODO(), contextKey{}, "dummy")

	gotTraceID := MustGetTraceID(ctx)

	if gotTraceID != "dummy" {
		t.Errorf("trace id %v; want dummy", gotTraceID)
	}
}

func TestGetTraceID_failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panicked false; want true")
		}
	}()

	MustGetTraceID(context.TODO())
}

func TestUnaryServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if _, ok := ctx.Value(contextKey{}).(string); !ok {
			return nil, errors.New("")
		}
		return nil, nil
	}

	_, err := UnaryServerInterceptor()(ctx, nil, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	stream := test.ServerStreamMock{Ctx: ctx}
	info := grpc.StreamServerInfo{}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		if _, ok := stream.Context().Value(contextKey{}).(string); !ok {
			return errors.New("")
		}
		return nil
	}

	err := StreamServerInterceptor()(nil, stream, &info, handler)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}
