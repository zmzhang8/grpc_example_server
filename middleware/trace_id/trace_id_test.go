package trace_id

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"

	"github.com/zmzhang8/grpc_example/test"
)

func TestRequestIdFromContext_success(t *testing.T) {
	wantRequestId := "xxx"
	ctx := context.WithValue(context.TODO(), ContextKey, wantRequestId)

	gotRequestId := TraceIdFromContext(ctx)

	if gotRequestId != wantRequestId {
		t.Errorf("trace id %v; want %v", gotRequestId, wantRequestId)
	}
}

func TestRequestIdFromContext_failure(t *testing.T) {
	wantRequestId := ""
	ctx := context.TODO()

	gotRequestId := TraceIdFromContext(ctx)

	if gotRequestId != wantRequestId {
		t.Errorf("trace id %v; want %v", gotRequestId, wantRequestId)
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "yyy", errors.New("zzz")
	}

	resp, err := UnaryServerInterceptor()(ctx, nil, &info, handler)

	if resp != "yyy" {
		t.Errorf("resp %v; want yyy", resp)
	}
	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	stream := test.ServerStreamMock{Ctx: ctx}
	info := grpc.StreamServerInfo{}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New("zzz")
	}

	err := StreamServerInterceptor()(nil, stream, &info, handler)

	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}
