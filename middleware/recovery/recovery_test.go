package recovery

import (
	"context"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/log"
	"github.com/zmzhang8/grpc_example/test"
)

func TestUnaryServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("")
	}
	wantErr := status.Errorf(codes.Internal, "panicked")

	_, err := UnaryServerInterceptor(logger)(ctx, nil, &info, handler)

	if err.Error() != wantErr.Error() {
		t.Errorf("err %v; want %v", err, wantErr)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	stream := test.ServerStreamMock{Ctx: context.TODO()}
	info := grpc.StreamServerInfo{}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		panic("")
	}
	wantErr := status.Errorf(codes.Internal, "panicked")

	err := StreamServerInterceptor(logger)(nil, stream, &info, handler)

	if err.Error() != wantErr.Error() {
		t.Errorf("err %v; want %v", err, wantErr)
	}
}
