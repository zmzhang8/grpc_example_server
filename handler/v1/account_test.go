package v1

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/zmzhang8/grpc_example/proto/v1"
)

func TestAccountServer_Login_success(t *testing.T) {
	s := NewAccountServer()
	ctx := context.TODO()
	req := pb.LoginRequest{
		Username: "hello",
		Password: "world",
	}
	wantToken := "worldhello"

	resp, err := s.Login(ctx, &req)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}
	if resp.Token != wantToken {
		t.Errorf("token %v; want %v", resp.Token, wantToken)
	}
}

func TestAccountServer_Login_failure(t *testing.T) {
	s := NewAccountServer()
	ctx := context.TODO()
	req := pb.LoginRequest{
		Username: "hello",
		Password: "earth",
	}
	wantErr := status.Error(codes.Unauthenticated, "Authentication failed")

	_, err := s.Login(ctx, &req)

	if err == nil {
		t.Errorf("err <nil>; want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("err %v; want %v", err, wantErr)
	}
}
