package test

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type ServerStreamMock struct {
	Ctx context.Context
}

func (s ServerStreamMock) SetHeader(metadata.MD) error { return nil }

func (s ServerStreamMock) SendHeader(metadata.MD) error { return nil }

func (s ServerStreamMock) SetTrailer(metadata.MD) {}

func (s ServerStreamMock) Context() context.Context { return s.Ctx }

func (s ServerStreamMock) SendMsg(m interface{}) error { return nil }

func (s ServerStreamMock) RecvMsg(m interface{}) error { return nil }
