package v1

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zmzhang8/grpc_example/lib/auth"
	pb "github.com/zmzhang8/grpc_example/proto/v1"
)

type accountServer struct {
	pb.UnimplementedAccountServer
}

func (s *accountServer) Login(
	ctx context.Context,
	in *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	if in.Username == "hello" && in.Password == "world" {
		return &pb.LoginResponse{
			Token:      "worldhello",
			Expiration: timestamppb.New(time.Now().UTC().Add(time.Hour)), // valid for one hour
		}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Authentication failed")
}

func (s *accountServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return auth.NoAuth(ctx)
}

func NewAccountServer() *accountServer {
	return &accountServer{}
}
