package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zmzhang8/grpc_example/lib/auth"
	"github.com/zmzhang8/grpc_example/middleware/logging"
	pb "github.com/zmzhang8/grpc_example/proto/v1"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(
	ctx context.Context,
	in *pb.HelloRequest,
) (*pb.HelloReply, error) {
	logger := logging.MustGetLogger(ctx)
	logger.Infow("Received ", "name", in.GetName())
	if len(in.Name) == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Name cannot be empty",
		)
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *greeterServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return auth.UserAuth(ctx)
}

func NewGreeterServer() *greeterServer {
	return &greeterServer{}
}
