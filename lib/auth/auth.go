package auth

import (
	"context"

	grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey struct{}

func MustGetAuthMetadata(ctx context.Context) string {
	metadata, ok := ctx.Value(contextKey{}).(string)
	if !ok {
		panic("cannot get auth metadata in context")
	}
	return metadata
}

func RejectAll(ctx context.Context) (context.Context, error) {
	return nil, status.Error(codes.Unauthenticated, "")
}

func AllowAll(ctx context.Context) (context.Context, error) {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/auth#pkg-types
	// The `Context` returned must be a child `Context` of the one passed in
	newCtx := context.WithValue(ctx, contextKey{}, "")
	return newCtx, nil
}

// Expected header
// key: authorization
// value: bearer {token}
func SessionAuth(ctx context.Context) (context.Context, error) {
	_, err := grpc_middleware_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.Unauthenticated, "")
}
