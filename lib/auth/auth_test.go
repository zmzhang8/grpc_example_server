package auth

import (
	"context"
	"testing"

	grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMustGetAuthMetadata_success(t *testing.T) {
	ctx := context.WithValue(context.TODO(), contextKey{}, "dummy")

	gotMetadata := MustGetAuthMetadata(ctx)

	if gotMetadata != "dummy" {
		t.Errorf("metadata %v; want dummy", gotMetadata)
	}
}

func TestMustGetAuthMetadata_failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panicked false; want true")
		}
	}()

	MustGetAuthMetadata(context.TODO())
}

func TestRejectAll(t *testing.T) {
	wantErr := status.Error(codes.Unauthenticated, "")

	gotCtx, gotErr := RejectAll(context.TODO())

	if gotCtx != nil {
		t.Errorf("context %v; want <nil>", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestAllowAll(t *testing.T) {
	ctx := context.TODO()
	wantCtx := context.WithValue(ctx, contextKey{}, "")

	gotCtx, gotErr := AllowAll(ctx)

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

func TestSessionAuth_failureNotBearer(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "basic worldhello",
	))
	_, wantErr := grpc_middleware_auth.AuthFromMD(ctx, "bearer")

	gotCtx, gotErr := SessionAuth(ctx)

	if gotCtx != nil {
		t.Errorf("context %v; want nil", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}
