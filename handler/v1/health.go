package v1

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/auth"
	pb "github.com/zmzhang8/grpc_example/proto/v1"
)

type healthServer struct {
	pb.UnimplementedHealthServer
}

func (s *healthServer) Check(
	ctx context.Context,
	in *pb.HealthCheckRequest,
) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(
	in *pb.HealthCheckRequest,
	stream pb.Health_WatchServer,
) error {
	stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "")
		case <-ticker.C:
			stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})
		}
	}
}

func (s *healthServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return auth.NoAuth(ctx)
}

func NewHealthServer() *healthServer {
	return &healthServer{}
}
