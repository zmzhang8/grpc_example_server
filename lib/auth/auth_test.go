package auth

import (
	"context"
	"testing"

	grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestDefaultAuth(t *testing.T) {
	wantErr := status.Error(codes.Unauthenticated, "Unauthenticated")

	gotCtx, gotErr := DefaultAuth(context.TODO())

	if gotCtx != nil {
		t.Errorf("context %v; want <nil>", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestNoAuth(t *testing.T) {
	ctx := context.TODO()
	wantCtx := context.WithValue(ctx, contextKey{}, "")

	gotCtx, gotErr := NoAuth(ctx)

	if gotCtx == nil {
		t.Errorf("context <nil>; want %v", wantCtx)
	}
	gotCtxValue, ok := gotCtx.Value(contextKey{}).(string)
	if !ok || gotCtxValue != "" {
		t.Errorf("context %v; want %v", gotCtx, wantCtx)
	}
	if gotErr != nil {
		t.Errorf("error %v; want <nil>", gotErr)
	}
}

func TestUserAuth_failureNotBasic(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "basic worldhello",
	))
	_, wantErr := grpc_middleware_auth.AuthFromMD(ctx, "bearer")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx != nil {
		t.Errorf("context %v; want nil", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestUserAuth_failureWrongToken(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "bearer xxx",
	))
	wantErr := status.Error(codes.Unauthenticated, "Unauthenticated")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx != nil {
		t.Errorf("context %v; want nil", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestUserAuth_success(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "bearer worldhello",
	))
	wantCtx := context.WithValue(ctx, contextKey{}, "")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx == nil {
		t.Errorf("context <nil>; want %v", wantCtx)
	}
	gotCtxValue, ok := gotCtx.Value(contextKey{}).(string)
	if !ok || gotCtxValue != "" {
		t.Errorf("context %v; want %v", gotCtx, wantCtx)
	}
	if gotErr != nil {
		t.Errorf("error %v; want <nil>", gotErr)
	}
}
