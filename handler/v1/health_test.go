package v1

import (
	"context"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/log"
	"github.com/zmzhang8/grpc_example/middleware/logging"
	pb "github.com/zmzhang8/grpc_example/proto/v1"
	"github.com/zmzhang8/grpc_example/test"
)

type healthWatchServerMock struct {
	grpc.ServerStream
}

func (s *healthWatchServerMock) Send(resp *pb.HealthCheckResponse) error {
	return s.ServerStream.SendMsg(resp)
}

func TestHealthServer_Check_success(t *testing.T) {
	s := NewHealthServer()
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), logging.ContextKey(), logger)
	req := pb.HealthCheckRequest{}

	_, err := s.Check(ctx, &req)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
}

func TestHealthServer_Watch_success(t *testing.T) {
	s := NewHealthServer()
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx, cancelFunc := context.WithTimeout(
		context.WithValue(context.TODO(), logging.ContextKey(), logger),
		time.Second,
	)
	defer cancelFunc()
	stream := healthWatchServerMock{test.ServerStreamMock{Ctx: ctx}}
	req := pb.HealthCheckRequest{}
	wantErr := status.Errorf(codes.Canceled, "")

	err := s.Watch(&req, &stream)

	if err == nil {
		t.Errorf("err <nil>; want %v", wantErr)
	}
}
