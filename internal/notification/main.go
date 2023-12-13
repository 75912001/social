package notification

import (
	"context"
	pkgserver "social/pkg/server"
)

type Server struct {
	*pkgserver.Normal
}

func (p *Server) Stop(ctx context.Context) (err error) {
	return nil
}

func (p *Server) Start(ctx context.Context) (err error) {
	return nil
}
