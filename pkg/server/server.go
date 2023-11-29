package server

import "context"

type IServer interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}
